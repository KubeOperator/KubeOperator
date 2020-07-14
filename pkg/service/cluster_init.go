package service

import (
	"context"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
)

type ClusterInitService interface {
	Init(name string) error
	GatherKubernetesToken(cluster model.Cluster) error
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
	clusterIaasService         ClusterIaasService
}

func (c clusterInitService) Init(name string) error {
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return err
	}
	cluster.Status, err = c.clusterStatusRepo.Get(cluster.StatusID)
	if err != nil {
		return err
	}
	if len(cluster.Status.ClusterStatusConditions) > 0 {
		for i, _ := range cluster.Status.ClusterStatusConditions {
			if cluster.Status.ClusterStatusConditions[i].Status == constant.ConditionFalse {
				cluster.Status.ClusterStatusConditions[i].Status = constant.ConditionUnknown
				cluster.Status.ClusterStatusConditions[i].Message = ""
				err := c.clusterStatusConditionRepo.Save(&cluster.Status.ClusterStatusConditions[i])
				if err != nil {
					return err
				}
			}
		}
	}
	go c.do(cluster)
	return nil
}

func (c clusterInitService) do(cluster model.Cluster) {
	_ = c.clusterIaasService.Init(cluster.Name)
	cluster.Status.Phase = constant.ClusterInitializing
	_ = c.clusterStatusRepo.Save(&cluster.Status)
	ctx, cancel := context.WithCancel(context.Background())
	admCluster := adm.NewCluster(cluster)
	statusChan := make(chan adm.Cluster, 0)
	go c.doCreate(ctx, *admCluster, statusChan)
	for {
		cluster := <-statusChan
		_ = c.clusterStatusRepo.Save(&cluster.Status)
		switch cluster.Status.Phase {
		case constant.ClusterFailed:
			cancel()
			return
		case constant.ClusterRunning:
			for i, _ := range cluster.Nodes {
				cluster.Nodes[i].Status = constant.ClusterRunning
				_ = c.clusterNodeRepo.Save(&cluster.Nodes[i])
			}
			cancel()
			err := c.GatherKubernetesToken(cluster.Cluster)
			if err != nil {
				cluster.Status.Phase = constant.ClusterNotConnected
				cluster.Status.Message = err.Error()
			}
			return
		}
	}
}

func (c clusterInitService) doCreate(ctx context.Context, cluster adm.Cluster, statusChan chan adm.Cluster) {
	ad := adm.NewClusterAdm()
	for {
		select {
		case <-ctx.Done():
			return
		case statusChan <- cluster:
		}
		resp, err := ad.OnInitialize(cluster)
		if err != nil {
			cluster.Status.Message = err.Error()
		}
		cluster.Status = resp.Status
	}
}

func (c clusterInitService) GatherKubernetesToken(cluster model.Cluster) error {
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
