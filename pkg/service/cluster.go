package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
)

type ClusterService interface {
	Get(name string) (dto.Cluster, error)
	GetStatus(name string) (dto.ClusterStatus, error)
	GetSecrets(name string) (dto.ClusterSecret, error)
	GetEndpoint(name string) (string, error)
	Delete(name string) error
	Create(creation dto.ClusterCreate) error
	List() ([]dto.Cluster, error)
	Page(num, size int) (dto.ClusterPage, error)
}

func NewClusterService() ClusterService {
	return &clusterService{
		clusterRepo:                repository.NewClusterRepository(),
		clusterSpecRepo:            repository.NewClusterSpecRepository(),
		clusterNodeRepo:            repository.NewClusterNodeRepository(),
		clusterStatusRepo:          repository.NewClusterStatusRepository(),
		clusterSecretRepo:          repository.NewClusterSecretRepository(),
		clusterStatusConditionRepo: repository.NewClusterStatusConditionRepository(),
	}
}

type clusterService struct {
	clusterRepo                repository.ClusterRepository
	clusterSpecRepo            repository.ClusterSpecRepository
	clusterNodeRepo            repository.ClusterNodeRepository
	clusterStatusRepo          repository.ClusterStatusRepository
	clusterSecretRepo          repository.ClusterSecretRepository
	clusterStatusConditionRepo repository.ClusterStatusConditionRepository
}

func (c clusterService) Get(name string) (dto.Cluster, error) {
	var clusterDTO dto.Cluster
	mo, err := c.clusterRepo.Get(name)
	if err != nil {
		return clusterDTO, err
	}
	clusterDTO.Cluster = mo
	clusterDTO.NodeSize = len(mo.Nodes)
	clusterDTO.Status = mo.Status.Phase
	return clusterDTO, nil
}

func (c clusterService) List() ([]dto.Cluster, error) {
	var clusterDTOS []dto.Cluster
	mos, err := c.clusterRepo.List()
	if err != nil {
		return clusterDTOS, nil
	}
	for _, mo := range mos {
		clusterDTOS = append(clusterDTOS, dto.Cluster{
			Cluster:  mo,
			NodeSize: len(mo.Nodes),
			Status:   mo.Status.Phase,
		})
	}
	return clusterDTOS, err
}

func (c clusterService) Page(num, size int) (dto.ClusterPage, error) {
	var page dto.ClusterPage
	total, mos, err := c.clusterRepo.Page(num, size)
	if err != nil {
		return page, nil
	}
	for _, mo := range mos {
		page.Items = append(page.Items, dto.Cluster{
			Cluster:  mo,
			NodeSize: len(mo.Nodes),
			Status:   mo.Status.Phase,
		})
	}
	page.Total = total
	return page, err
}

func (c clusterService) GetSecrets(name string) (dto.ClusterSecret, error) {
	var secret dto.ClusterSecret
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return secret, err
	}
	cs, err := c.clusterSecretRepo.Get(cluster.SecretID)
	if err != nil {
		return secret, err
	}
	secret.ClusterSecret = cs

	return secret, nil
}

func (c clusterService) GetStatus(name string) (dto.ClusterStatus, error) {
	var status dto.ClusterStatus
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return status, err
	}
	cs, err := c.clusterStatusRepo.Get(cluster.StatusID)
	if err != nil {
		return status, err
	}
	status.ClusterStatus = cs
	return status, nil

}

func (c clusterService) Create(creation dto.ClusterCreate) error {
	tx := db.DB.Begin()
	spec := model.ClusterSpec{
		RuntimeType:          creation.RuntimeType,
		DockerStorageDir:     creation.DockerStorageDIr,
		ContainerdStorageDir: creation.ContainerdStorageDIr,
		NetworkType:          creation.NetworkType,
		ClusterCIDR:          creation.ClusterCIDR,
		ServiceCIDR:          creation.ServiceCIDR,
		Version:              creation.Version,
		AppDomain:            creation.AppDomain,
	}
	if err := c.clusterSpecRepo.Save(&spec); err != nil {
		tx.Rollback()
		return err
	}
	status := model.ClusterStatus{Phase: constant.ClusterWaiting}
	if err := c.clusterStatusRepo.Save(&status); err != nil {
		tx.Rollback()
		return err
	}
	secret := model.ClusterSecret{
		KubeadmToken: clusterUtil.GenerateKubeadmToken(),
	}
	if err := c.clusterSecretRepo.Save(&secret); err != nil {
		tx.Rollback()
		return err
	}
	cluster := model.Cluster{
		BaseModel: common.BaseModel{},
		Name:      creation.Name,
		SpecID:    spec.ID,
		StatusID:  status.ID,
		SecretID:  secret.ID,
	}
	if err := c.clusterRepo.Save(&cluster); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c clusterService) GetEndpoint(name string) (string, error) {
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return "", err
	}
	if cluster.Spec.LbKubeApiserverIp != "" {
		return cluster.Spec.LbKubeApiserverIp, nil
	}
	master, err := c.clusterNodeRepo.FistMaster(cluster.ID)
	if err != nil {
		return "", err
	}
	return master.Host.Ip, nil
}

func (c clusterService) Delete(name string) error {
	return c.clusterRepo.Delete(name)
}
