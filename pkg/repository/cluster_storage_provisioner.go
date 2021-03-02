package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterStorageProvisionerRepository interface {
	List(clusterName string) ([]model.ClusterStorageProvisioner, error)
	Save(clusterName string, provisioner *model.ClusterStorageProvisioner) error
	Delete(clusterName string, provisionerName string) error
	BatchDelete(clusterName string, items []dto.ClusterStorageProvisioner) error
}

type clusterStorageProvisionerRepository struct {
}

func NewClusterStorageProvisionerRepository() ClusterStorageProvisionerRepository {
	return &clusterStorageProvisionerRepository{}
}

func (c clusterStorageProvisionerRepository) List(clusterName string) ([]model.ClusterStorageProvisioner, error) {
	var cluster model.Cluster
	var ps []model.ClusterStorageProvisioner
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return ps, err
	}
	if err := db.DB.Where("cluster_id = ?", cluster.ID).Find(&ps).Error; err != nil {
		return ps, err
	}
	return ps, nil
}

func (c clusterStorageProvisionerRepository) Save(clusterName string, provisioner *model.ClusterStorageProvisioner) error {
	var cluster model.Cluster
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
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

func (c clusterStorageProvisionerRepository) Delete(clusterName string, provisionerName string) error {
	var cluster model.Cluster
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return err
	}
	var provisioner model.ClusterStorageProvisioner
	if err := db.DB.Where("name = ?", provisionerName).First(&provisioner).Error; err != nil {
		return err
	}
	err := db.DB.Delete(&provisioner).Error
	if err != nil {
		return err
	}
	return nil
}

func (c clusterStorageProvisionerRepository) BatchDelete(clusterName string, items []dto.ClusterStorageProvisioner) error {
	tx := db.DB.Begin()
	for _, item := range items {
		err := c.Delete(clusterName, item.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
