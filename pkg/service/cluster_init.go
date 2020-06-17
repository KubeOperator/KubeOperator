package service

import (
	"context"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"time"
)

type ClusterInitService interface {
	Init(name string) error
}

func NewClusterInitService() ClusterInitService {
	return &clusterInitService{
		clusterRepo:                repository.NewClusterRepository(),
		clusterNodeRepo:            repository.NewClusterNodeRepository(),
		clusterStatusRepo:          repository.NewClusterStatusRepository(),
		clusterSecretRepo:          repository.NewClusterSecretRepository(),
		clusterStatusConditionRepo: repository.NewClusterStatusConditionRepository(),
	}
}

type clusterInitService struct {
	clusterRepo                repository.ClusterRepository
	clusterNodeRepo            repository.ClusterNodeRepository
	clusterStatusRepo          repository.ClusterStatusRepository
	clusterSecretRepo          repository.ClusterSecretRepository
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
	if len(status.ClusterStatusConditions) > 0 {
		for i, _ := range status.ClusterStatusConditions {
			if status.ClusterStatusConditions[i].Status == constant.ConditionFalse {
				status.ClusterStatusConditions[i].Status = constant.ConditionUnknown
				status.ClusterStatusConditions[i].Message = ""
				err := c.clusterStatusConditionRepo.Save(&status.ClusterStatusConditions[i])
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
	ctx, cancel := context.WithCancel(context.Background())
	admCluster := adm.NewCluster(cluster)
	statusChan := make(chan adm.Cluster, 0)
	go c.do(ctx, *admCluster, statusChan)
	go c.pollingStatus(cancel, statusChan)
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

func (c clusterInitService) pollingStatus(cancel context.CancelFunc, statusChan chan adm.Cluster) {
	for {
		cluster := <-statusChan
		_ = c.clusterStatusRepo.Save(&cluster.Status)
		switch cluster.Status.Phase {
		case constant.ClusterFailed:
			cancel()
			return
		case constant.ClusterRunning:
			cancel()
			err := c.gatherKubernetesToken(cluster.Cluster)
			if err != nil {
				cluster.Status.Phase = constant.ClusterNotConnected
				cluster.Status.Message = err.Error()
			}
			return
		}
	}
}

func (c clusterInitService) gatherKubernetesToken(cluster model.Cluster) error {
	secret, err := c.clusterSecretRepo.Get(cluster.SecretID)
	if err != nil {
		return err
	}
	master, err := c.clusterNodeRepo.FistMaster(cluster.ID)
	if err != nil {
		return err
	}
	sshConfig := master.ToSSHConfig()
	client, err := ssh.New(&sshConfig)
	if err != nil {
		return err
	}
	token, err := clusterUtil.GetClusterToken(client)
	if err != nil {
		return err
	}
	secret.KubernetesToken = token
	return c.clusterSecretRepo.Save(&secret)
}
