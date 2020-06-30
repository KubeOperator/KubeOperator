package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterStorageProvisionerRepository interface {
	List(clusterName string) ([]model.ClusterStorageProvisioner, error)
	Save(clusterName string, provisioner *model.ClusterStorageProvisioner) error
}

type clusterStorageProvisionerRepository struct {
}

func NewClusterStorageProvisionerRepository() ClusterStorageProvisionerRepository {
	return &clusterStorageProvisionerRepository{}
}

func (c clusterStorageProvisionerRepository) List(clusterName string) ([]model.ClusterStorageProvisioner, error) {
	var cluster model.Cluster
	var ps []model.ClusterStorageProvisioner
	if err := db.DB.
		Where(model.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return ps, err
	}
	if err := db.DB.
		Where(model.ClusterStorageProvisioner{ClusterID: cluster.ID}).
		Find(&ps).Error; err != nil {
		return ps, err
	}
	return ps, nil
}

func (c clusterStorageProvisionerRepository) Save(clusterName string, provisioner *model.ClusterStorageProvisioner) error {
	var cluster model.Cluster
	if err := db.DB.
		Where(model.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return err
	}
	provisioner.ClusterID = cluster.ID
	if db.DB.NewRecord(provisioner) {
		if err := db.DB.Create(provisioner).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(provisioner).Error; err != nil {
			return nil
		}
	}
	return nil
}
