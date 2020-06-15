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

}

func NewClusterRepository() ClusterRepository {
	return &clusterRepository{}
}

type clusterRepository struct{}

func (c clusterRepository) Get(name string) (model.Cluster, error) {
	var cluster model.Cluster
	if err := db.DB.Where(model.Cluster{Name: name}).Find(&cluster).Error; err != nil {
		return cluster, err
	}
	return cluster, nil
}

func (c clusterRepository) List() ([]model.Cluster, error) {
	var clusters []model.Cluster
	if err := db.DB.Find(&clusters).Error; err != nil {
		return clusters, err
	}
	return clusters, nil
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

func (c clusterRepository) Delete(id string) error {
	cluster := model.Cluster{ID: id}
	if err := db.DB.First(&cluster).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&cluster).Error; err != nil {
		return err
	}
	return nil
}
