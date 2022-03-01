package job

import (
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
)

type ClusterHealthCheck struct {
	clusterService service.ClusterService
}

func NewClusterHealthCheck() *ClusterHealthCheck {
	return &ClusterHealthCheck{
		clusterService: service.NewClusterService(),
	}
}

func (c *ClusterHealthCheck) Run() {
	cs, err := c.clusterService.List()
	if err != nil {
		log.Errorf("list clusters error %s", err.Error())
		return
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // 信号量
	for i := range cs {
		if cs[i].Status != constant.StatusRunning && cs[i].Status != constant.StatusLost {
			continue
		}
		wg.Add(1)
		go func(item int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			log.Infof("test cluster  %s api  ", cs[item].Name)
			endpoints, err := c.clusterService.GetApiServerEndpoints(cs[item].Name)
			if err != nil {
				log.Errorf("get cluster %s endpoint error %s", cs[item].Name, err.Error())
				return
			}
			secret, err := c.clusterService.GetSecrets(cs[item].Name)
			if err != nil {
				log.Errorf("get cluster %s secret error %s", cs[item].Name, err.Error())
				return
			}
			_, err = kubeUtil.SelectAliveHost(endpoints)
			if err != nil {
				log.Errorf("ping cluster %s api error %s", cs[item].Name, err.Error())
				cs[item].Cluster.Status.Phase = constant.StatusLost
				if err := db.DB.Save(&cs[item].Cluster.Status).Error; err != nil {
					log.Errorf("save cluster %s status error %s", cs[item].Name, err.Error())
					return
				}
				return
			}
			client, err := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
				Hosts: endpoints,
				Token: secret.KubernetesToken,
			})
			if err != nil {
				log.Errorf("get cluster %s api client error %s", cs[item].Name, err.Error())
				return
			}
			_, err = client.ServerVersion()
			if err != nil {
				log.Errorf("ping cluster %s api error %s", cs[item].Name, err.Error())
				cs[item].Cluster.Status.Phase = constant.StatusLost
				if err := db.DB.Save(&cs[item].Cluster.Status).Error; err != nil {
					log.Errorf("save cluster %s status error %s", cs[item].Name, err.Error())
					return
				}
				return
			}
			if cs[item].Cluster.Status.Phase == constant.StatusLost {
				cs[item].Cluster.Status.Phase = constant.StatusRunning
				if err := db.DB.Save(&cs[item].Cluster.Status).Error; err != nil {
					log.Errorf("save cluster %s status error %s", cs[item].Name, err.Error())
					return
				}
			}
		}(i)
	}
	wg.Wait()
}
