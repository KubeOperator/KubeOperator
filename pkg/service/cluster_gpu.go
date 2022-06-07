package service

import (
	"io"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/jinzhu/gorm"
	"k8s.io/client-go/kubernetes"
)

const (
	gpuOperator = "16-gpu-operator.yml"
)

type ClusterGpuService interface {
	Get(clusterName string) (*dto.ClusterGpu, error)
	Add(clusterID string) error
	HandleGPU(clusterName, operation string) (*dto.ClusterGpu, error)
}

func NewClusterGpuService() ClusterGpuService {
	return &clusterGpuService{
		clusterRepo:     repository.NewClusterRepository(),
		clusterSpecRepo: repository.NewClusterSpecRepository(),
		clusterGpuRepo:  repository.NewClusterGpuRepository(),
		clusterService:  NewClusterService(),
	}
}

type clusterGpuService struct {
	clusterRepo     repository.ClusterRepository
	clusterSpecRepo repository.ClusterSpecRepository
	clusterService  ClusterService
	clusterGpuRepo  repository.ClusterGpuRepository
}

func (c clusterGpuService) Get(clusterName string) (*dto.ClusterGpu, error) {
	gpuInfo, err := c.clusterGpuRepo.GetByClusterName(clusterName)
	return &dto.ClusterGpu{ID: gpuInfo.ID, Status: gpuInfo.Status, Message: gpuInfo.Message}, err
}

func (c clusterGpuService) Add(clusterID string) error {
	return c.clusterGpuRepo.Add(clusterID)
}

func (c clusterGpuService) HandleGPU(clusterName string, operation string) (*dto.ClusterGpu, error) {
	cluster, err := c.clusterRepo.Get(clusterName)
	if err != nil {
		return nil, err
	}

	status := constant.StatusCreating
	if operation == "disable" {
		status = constant.StatusTerminating
	}

	cluster.SpecConf.SupportGpu = status
	if err := c.clusterSpecRepo.SaveConf(&cluster.SpecConf); err != nil {
		return nil, err
	}
	var gpuInfo model.ClusterGpu
	if operation == "disable" {
		gpuInfo, err = c.clusterGpuRepo.GetByClusterID(cluster.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	} else {
		if err := c.clusterGpuRepo.Delete(cluster.ID); err != nil {
			return nil, err
		}
		gpuInfo.ClusterID = cluster.ID
		gpuInfo.Status = status
		if err := c.clusterGpuRepo.Save(&gpuInfo); err != nil {
			return nil, err
		}
	}

	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, gpuInfo.ID)
	if err != nil {
		return nil, err
	}

	go c.handleGpu(gpuInfo, operation, cluster, writer)

	return &dto.ClusterGpu{ID: gpuInfo.ID, Status: gpuInfo.Status}, nil
}

func (c clusterGpuService) handleGpu(gpuInfo model.ClusterGpu, operation string, cluster model.Cluster, writer io.Writer) {
	client := adm.NewCluster(cluster)
	client.Kobe.SetVar(facts.SupportGpuFactName, operation)
	if err := phases.RunPlaybookAndGetResult(client.Kobe, gpuOperator, "", writer); err != nil {
		c.errHandleGpu(cluster, gpuInfo, constant.StatusFailed, err)
		return
	}

	if operation == "disable" {
		_ = c.clusterGpuRepo.Delete(cluster.ID)

		cluster.SpecConf.SupportGpu = constant.StatusDisabled
		_ = c.clusterSpecRepo.SaveConf(&cluster.SpecConf)
		return
	}

	cluster.SpecConf.SupportGpu = constant.StatusWaiting
	_ = c.clusterSpecRepo.SaveConf(&cluster.SpecConf)

	k8sClient, err := c.getBaseParam(cluster.Name)
	if err != nil {
		c.errHandleGpu(cluster, gpuInfo, constant.StatusFailed, err)
		return
	}

	if err := phases.WaitForDeployRunning("kube-operator", "gpu-operator", k8sClient); err != nil {
		c.errHandleGpu(cluster, gpuInfo, constant.StatusNotReady, err)
		return
	}

	gpuInfo.Status = constant.StatusRunning
	_ = c.clusterGpuRepo.Save(&gpuInfo)

	cluster.SpecConf.SupportGpu = constant.StatusEnabled
	_ = c.clusterSpecRepo.SaveConf(&cluster.SpecConf)
}

func (c clusterGpuService) errHandleGpu(cluster model.Cluster, gpuInfo model.ClusterGpu, Status string, err error) {
	logger.Log.Errorf(err.Error())
	gpuInfo.Status = Status
	gpuInfo.Message = err.Error()
	_ = c.clusterGpuRepo.Save(&gpuInfo)

	cluster.SpecConf.SupportGpu = Status
	_ = c.clusterSpecRepo.SaveConf(&cluster.SpecConf)
}

func (c clusterGpuService) getBaseParam(clusterName string) (*kubernetes.Clientset, error) {
	var client *kubernetes.Clientset
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return client, err
	}

	endpoints, err := c.clusterService.GetApiServerEndpoints(clusterName)
	if err != nil {
		return client, err
	}

	client, err = kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Token: secret.KubernetesToken,
		Hosts: endpoints,
	})
	if err != nil {
		return client, err
	}
	return client, nil
}
