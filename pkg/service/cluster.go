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
	//GetNodes(name string)
	Delete(name string) error
	Create(creation dto.ClusterCreate) error
	List() ([]dto.Cluster, error)
}

//func NewClusterService() ClusterService {
//	return &clusterService{
//		repo: repository.NewClusterRepository(),
//	}
//}

type clusterService struct {
	clusterRepo                repository.ClusterRepository
	clusterSpecRepo            repository.ClusterSpecRepository
	clusterStatusRepo          repository.ClusterStatusRepository
	clusterSecretRepo          repository.ClusterSecretRepository
	clusterStatusConditionRepo repository.ClusterStatusConditionRepository
}

func (c clusterService) Get(name string) (dto.Cluster, error) {
	var cluster dto.Cluster
	mo, err := c.clusterRepo.Get(name)
	if err != nil {
		return cluster, err
	}
	spec, err := c.clusterSpecRepo.Get(mo.SpecID)
	if err != nil {
		return cluster, err
	}
	status, err := c.clusterStatusRepo.Get(mo.SpecID)
	if err != nil {
		return cluster, err
	}

	cluster.Cluster = mo
	cluster.Spec = spec
	cluster.Status = status
	return cluster, nil
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
	conditions, err := c.clusterStatusConditionRepo.List(cluster.StatusID)
	if err != nil {
		return status, err
	}
	status.ClusterStatus = cs
	status.Conditions = conditions
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
