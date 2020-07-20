package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterRepository interface {
	Get(name string) (model.Cluster, error)
	List() ([]model.Cluster, error)
	Save(cluster *model.Cluster) error
	Delete(name string) error
	Page(num, size int) (int, []model.Cluster, error)
}

func NewClusterRepository() ClusterRepository {
	return &clusterRepository{}
}

type clusterRepository struct {
}

func (c clusterRepository) Get(name string) (model.Cluster, error) {
	var cluster model.Cluster
	if err := db.DB.
		Where(model.Cluster{Name: name}).
		Preload("Status").
		Preload("Spec").
		Preload("Nodes").
		Preload("Nodes.Host").
		Preload("Nodes.Host.Credential").
		Preload("Nodes.Host.Zone").
		Find(&cluster).Error; err != nil {
		return cluster, err
	}
	return cluster, nil
}

func (c clusterRepository) List() ([]model.Cluster, error) {
	var clusters []model.Cluster
	db.DB.Model(model.Cluster{})
	if err := db.DB.Model(model.Cluster{}).
		Preload("Status").
		Preload("Spec").
		Preload("Nodes").
		Preload("Nodes.Host").
		Preload("Nodes.Host.Credential").
		Find(&clusters).Error; err != nil {
		return clusters, err
	}
	return clusters, nil
}

func (c clusterRepository) Page(num, size int) (int, []model.Cluster, error) {
	var total int
	var clusters []model.Cluster
	if err := db.DB.Model(model.Cluster{}).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Preload("Status").
		Preload("Spec").
		Preload("Nodes").
		Find(&clusters).Error; err != nil {
		return total, clusters, err
	}
	return total, clusters, nil
}

func (c clusterRepository) Save(cluster *model.Cluster) error {
	if db.DB.NewRecord(cluster) {
		if err := db.DB.Create(cluster).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(cluster).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c clusterRepository) Delete(name string) error {
	var cluster model.Cluster
	if err := db.DB.Where(model.Cluster{Name: name}).First(&cluster).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&cluster).Error; err != nil {
		return err
	}
	return nil
}
