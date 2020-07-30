package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterLogRepository interface {
	Save(clusterName string, log *model.ClusterLog) error
	List(clusterName string) ([]model.ClusterLog, error)
}

func NewClusterLogRepository() ClusterLogRepository {
	return &clusterLogRepository{}
}

type clusterLogRepository struct {
}

func (c *clusterLogRepository) Save(clusterName string, log *model.ClusterLog) error {
	var cluster model.Cluster
	if err := db.DB.Where(model.Cluster{Name: clusterName}).First(&cluster).Error; err != nil {
		return err
	}
	log.ClusterID = cluster.ID
	if db.DB.NewRecord(log) {
		return db.DB.Create(log).Error
	} else {
		return db.DB.Save(log).Error
	}
}

func (c *clusterLogRepository) List(clusterName string) ([]model.ClusterLog, error) {
	var cluster model.Cluster
	if err := db.DB.Where(model.Cluster{Name: clusterName}).First(&cluster).Error; err != nil {
		return nil, err
	}
	var items []model.ClusterLog
	if err := db.DB.Where(model.ClusterLog{ClusterID: cluster.ID}).
		Find(&items).
		Order("create_at desc").
		Error; err != nil {
		return nil, err
	}
	return items, nil
}
