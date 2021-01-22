package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterToolRepository interface {
	List(clusterName string) ([]model.ClusterTool, error)
	Save(tool *model.ClusterTool) error
	Get(clusterName string, name string) (model.ClusterTool, error)
}

func NewClusterToolRepository() ClusterToolRepository {
	return &clusterToolRepository{}
}

type clusterToolRepository struct{}

func (c clusterToolRepository) List(clusterName string) ([]model.ClusterTool, error) {
	var cluster model.Cluster
	var tools []model.ClusterTool
	if err := db.DB.
		Where(&model.Cluster{Name: clusterName}).
		Preload("Spec").
		First(&cluster).Error; err != nil {
		return tools, err
	}
	if err := db.DB.
		Where(&model.ClusterTool{ClusterID: cluster.ID}).
		Where("architecture in (?)", []string{cluster.Spec.Architectures, "all"}).
		Find(&tools).Error; err != nil {
		return tools, err
	}
	return tools, nil
}

func (c clusterToolRepository) Save(tool *model.ClusterTool) error {
	var item model.ClusterTool
	notFound := db.DB.Where(&model.ClusterTool{ClusterID: tool.ClusterID, Name: tool.Name}).First(&item).RecordNotFound()
	if notFound {
		if err := db.DB.Create(tool).Error; err != nil {
			return err
		}
	} else {
		tool.ID = item.ID
		if err := db.DB.Save(tool).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c clusterToolRepository) Get(clusterName string, name string) (model.ClusterTool, error) {
	var tool model.ClusterTool
	var cluster model.Cluster
	if err := db.DB.
		Where(&model.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return tool, err
	}
	if err := db.DB.Where(&model.ClusterTool{ClusterID: tool.ClusterID, Name: name}).First(&tool).Error; err != nil {
		return tool, err
	}
	return tool, nil
}
