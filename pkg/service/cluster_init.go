package service

import (
	"context"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"time"
)

type ClusterInitService interface {
	Init(name string) error
}

type clusterInitService struct {
	clusterRepo                repository.ClusterRepository
	clusterStatusRepo          repository.ClusterStatusRepository
	clusterStatusConditionRepo repository.ClusterStatusConditionRepository
}

func (c clusterInitService) Init(name string) error {
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return err
	}
	status, err := c.clusterStatusRepo.Get(cluster.StatusID)
	if err != nil {
		return err
	}
	if len(status.Conditions) > 0 {
		for i, _ := range status.Conditions {
			if status.Conditions[i].Status == constant.ConditionFalse {
				status.Conditions[i].Status = constant.ConditionUnknown
				status.Conditions[i].Message = ""
				err := c.clusterStatusConditionRepo.Save(&status.Conditions[i])
				if err != nil {
					return err
				}
			}
		}
	}
	status.Phase = constant.ClusterInitializing
	if err := c.clusterStatusRepo.Save(&status); err != nil {
		return err
	}
	ctx := context.Background()
	admCluster := adm.NewCluster(cluster)
	statusChan := make(chan adm.Cluster, 0)
	go c.do(ctx, *admCluster, statusChan)
	go c.pollingStatus(ctx, statusChan)
	return nil

}

func (c clusterInitService) do(ctx context.Context, cluster adm.Cluster, statusChan chan adm.Cluster) {
	ad := adm.NewClusterAdm()
	for {
		resp, err := ad.OnInitialize(cluster)
		if err != nil {
			cluster.Status.Message = err.Error()
		}
		cluster.Status = resp.Status
		select {
		case <-ctx.Done():
			return
		default:
			statusChan <- cluster
		}
		time.Sleep(5 * time.Second)
	}
}

func (c clusterInitService) pollingStatus(ctx context.Context, statusChan chan adm.Cluster) {

	for {
		cluster := <-statusChan
		_ = c.clusterStatusRepo.Save(&cluster.Status)
		switch cluster.Status.Phase {
		case constant.ClusterFailed:
			ctx.Done()
			return
		case constant.ClusterRunning:
			ctx.Done()
			return
		}

	}
}
