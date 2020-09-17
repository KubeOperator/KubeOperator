package service

import (
	"context"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	uuid "github.com/satori/go.uuid"
	"time"
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
		clusterSpecRepo:            repository.NewClusterSpecRepository(),
		clusterIaasService:         NewClusterIaasService(),
	}
}

type clusterInitService struct {
	clusterRepo                repository.ClusterRepository
	clusterNodeRepo            repository.ClusterNodeRepository
	clusterStatusRepo          repository.ClusterStatusRepository
	clusterSecretRepo          repository.ClusterSecretRepository
	clusterStatusConditionRepo repository.ClusterStatusConditionRepository
	clusterSpecRepo            repository.ClusterSpecRepository
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
	if len(cluster.Nodes) < 1 {
		cluster.Status.Phase = constant.ClusterCreating
		_ = c.clusterStatusRepo.Save(&cluster.Status)
		err := c.clusterIaasService.Init(cluster.Name)
		if err != nil {
			cluster.Status.Phase = constant.ClusterFailed
			cluster.Status.Message = err.Error()
			_ = c.clusterStatusRepo.Save(&cluster.Status)
			return
		}
	}
	cluster.Nodes, _ = c.clusterNodeRepo.List(cluster.Name)
	ctx, cancel := context.WithCancel(context.Background())
	statusChan := make(chan adm.Cluster, 0)
	cluster.Status.Phase = constant.ClusterInitializing
	_ = c.clusterStatusRepo.Save(&cluster.Status)
	logId := uuid.NewV4().String()
	writer, err := ansible.CreateAnsibleLogWriter(cluster.Name, logId)
	if err != nil {
		log.Error(err.Error())
	}
	admCluster := adm.NewCluster(cluster, writer)
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
				cluster.Spec.KubeRouter = cluster.Nodes[0].Host.Ip
				_ = c.clusterSpecRepo.Save(&cluster.Spec)
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
		resp, err := ad.OnInitialize(cluster)
		if err != nil {
			cluster.Status.Message = err.Error()
		}
		cluster.Status = resp.Status
		select {
		case <-ctx.Done():
			return
		case statusChan <- cluster:
		}
		time.Sleep(5 * time.Second)
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
