package job

import (
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
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
		logger.Log.Errorf("list clusters error %s", err.Error())
		return
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // 信号量
	for i := range cs {
		if cs[i].Status != constant.StatusRunning && cs[i].Status != constant.StatusLost {
			continue
		}
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			logger.Log.Infof("test cluster  %s api  ", cs[i].Name)
			client, err := kubeUtil.NewClusterClient(&cs[i].Cluster)
			if err != nil {
				logger.Log.Errorf("get cluster %s api client error %+v", cs[i].Name, err)
				return
			}
			_, err = client.ServerVersion()
			if err != nil {
				logger.Log.Errorf("ping cluster %s api error %s", cs[i].Name, err.Error())
				cs[i].Cluster.Status = constant.StatusLost
				if err := db.DB.Save(&cs[i].Cluster.Status).Error; err != nil {
					logger.Log.Errorf("save cluster %s status error %s", cs[i].Name, err.Error())
					return
				}
				return
			}
			if cs[i].Cluster.Status == constant.StatusLost {
				cs[i].Cluster.Status = constant.StatusRunning
				if err := db.DB.Save(&cs[i].Cluster.Status).Error; err != nil {
					logger.Log.Errorf("save cluster %s status error %s", cs[i].Name, err.Error())
					return
				}
			}
		}()
	}
	wg.Wait()
}
