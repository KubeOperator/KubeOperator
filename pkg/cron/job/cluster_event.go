package job

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterEvent struct {
	clusterService      service.ClusterService
	clusterEventService service.ClusterEventService
	messageService      service.MessageService
}

func NewClusterEvent() *ClusterEvent {
	return &ClusterEvent{
		clusterService:      service.NewClusterService(),
		clusterEventService: service.NewClusterEventService(),
		messageService:      service.NewMessageService(),
	}
}

func (c *ClusterEvent) Run() {
	fmt.Println("start event")
	logger.Log.Infof("start cluster event")
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // 信号量
	clusters, _ := c.clusterService.List()
	for _, cluster := range clusters {
		if cluster.Status != constant.StatusRunning {
			return
		}
		secret, err := c.clusterService.GetSecrets(cluster.Name)
		if err != nil {
			continue
		}
		endpoints, err := c.clusterService.GetApiServerEndpoints(cluster.Name)
		if cluster.Status == constant.ClusterRunning {
			client, err := kubernetes.NewKubernetesClient(&kubernetes.Config{
				Token: secret.KubernetesToken,
				Hosts: endpoints,
			})
			if err != nil {
				continue
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				namespaceList, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
				if err != nil {
					logger.Log.Errorf("list cluster %s namespace error : %s", cluster.Name, err.Error())
					return
				}
				for _, namespace := range namespaceList.Items {
					eventList, err := client.EventsV1beta1().Events(namespace.Name).List(context.Background(), metav1.ListOptions{})
					if err != nil {
						logger.Log.Errorf("list namespace %s event error : %s", namespace.Name, err.Error())
						return
					}
					for _, event := range eventList.Items {
						exist, err := c.clusterEventService.ExistEventUid(string(event.UID), cluster.ID)
						if err != nil {
							return
						}
						if !exist {
							clusterEvent := new(model.ClusterEvent)
							clusterEvent.UID = string(event.UID)
							clusterEvent.Name = event.Name
							clusterEvent.Type = event.Type
							clusterEvent.Namespace = event.Namespace
							clusterEvent.Reason = event.Reason
							clusterEvent.Kind = event.Regarding.Kind
							clusterEvent.Component = event.DeprecatedSource.Component
							clusterEvent.Host = event.DeprecatedSource.Host
							clusterEvent.ClusterID = cluster.ID
							clusterEvent.Message = event.Note

							if clusterEvent.Type == "Warning" {
								content, _ := json.Marshal(clusterEvent)
								err := c.messageService.SendMessage(constant.Cluster, false, string(content), cluster.Name, constant.ClusterEventWarning)
								if err != nil {
									logger.Log.Errorf("send cluster  %s event error : %s", cluster.Name, err.Error())
								}
							}
							err := c.clusterEventService.Save(*clusterEvent)
							if err != nil {
								return
							}
						}
					}
				}
			}()
		}
	}
	wg.Wait()
}
