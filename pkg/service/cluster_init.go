package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubeconfig"
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
		for i := range cluster.Status.ClusterStatusConditions {
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
	logId, fileName, err := ansible.CreateAnsibleLogWriter(cluster.Name)
	if err != nil {
		return err
	}
	cluster.LogId = logId
	_ = c.clusterRepo.Save(&cluster)
	go c.do(cluster, fileName)
	return nil
}

func (c clusterInitService) do(cluster model.Cluster, fileName string) {
	if len(cluster.Nodes) < 1 {
		return
	}
	cluster.Nodes, _ = c.clusterNodeRepo.List(cluster.Name)
	ctx, cancel := context.WithCancel(context.Background())
	statusChan := make(chan adm.Cluster)
	cluster.Status.Phase = constant.ClusterInitializing
	_ = c.clusterStatusRepo.Save(&cluster.Status)

	admCluster := adm.NewCluster(cluster, fileName)
	go c.doCreate(ctx, *admCluster, statusChan)
	for {
		cluster := <-statusChan
		_ = c.clusterStatusRepo.Save(&cluster.Status)
		switch cluster.Status.Phase {
		case constant.ClusterFailed:
			cancel()
			return
		case constant.ClusterRunning:
			for i := range cluster.Nodes {
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
	if statusChan == nil {
		statusChan = make(chan adm.Cluster)
	}
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
	master, err := c.clusterNodeRepo.FirstMaster(cluster.ID)
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

	kubeConf, err := c.getKubeconfig(cluster.Name)
	if err != nil {
		return err
	}
	secret.KubeConf = kubeConf

	return c.clusterSecretRepo.Save(&secret)
}

func (c clusterInitService) getKubeconfig(name string) (string, error) {
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return "", err
	}
	m, err := c.clusterNodeRepo.FirstMaster(cluster.ID)
	if err != nil {
		return "", err
	}
	cfg := m.ToSSHConfig()
	s, err := ssh.New(&cfg)
	if err != nil {
		return "", err
	}
	bf, err := kubeconfig.ReadKubeConfigFile(s)
	if err != nil {
		return "", err
	}
	configStr := string(bf)

	lbAddr := fmt.Sprintf("%s:%d", cluster.Spec.LbKubeApiserverIp, cluster.Spec.KubeApiServerPort)
	newStr := strings.ReplaceAll(configStr, "127.0.0.1:8443", lbAddr)

	return newStr, nil
}
