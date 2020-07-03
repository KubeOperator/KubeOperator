package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterToolRepository interface {
	List(clusterName string) ([]model.ClusterTool, error)
	Save(clusterName string, tool *model.ClusterTool) error
}

func NewClusterToolRepository() ClusterToolRepository {
	return &clusterToolRepository{}
}

type clusterToolRepository struct{}

func (c clusterToolRepository) List(clusterName string) ([]model.ClusterTool, error) {
	var cluster model.Cluster
	var tools []model.ClusterTool
	if err := db.DB.
		Where(model.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return tools, err
	}
	if err := db.DB.
		Where(model.ClusterTool{ClusterID: cluster.ID}).
		Find(&tools).Error; err != nil {
		return tools, err
	}
	return tools, nil
}

func (c clusterToolRepository) Save(clusterName string, tool *model.ClusterTool) error {
	var cluster model.Cluster
	if err := db.DB.
		Where(model.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return err
	}
	tool.ClusterID = cluster.ID
	if db.DB.NewRecord(tool) {
		if err := db.DB.Create(tool).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(tool).Error; err != nil {
			return nil
		}
	}
	return nil
}
