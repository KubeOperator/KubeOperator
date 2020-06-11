package tools

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	toolModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster/tool"
)

func Get(clusterName, name string) (toolModel.Tool, error) {
	var cluster clusterModel.Cluster
	var tool toolModel.Tool
	if err := db.DB.Where(clusterModel.Cluster{Name: clusterName}).First(&cluster).Error; err != nil {
		return tool, err
	}
	if err := db.DB.Where(toolModel.Tool{Name: name, ClusterID: cluster.ID}).
		First(&tool).Error; err != nil {
		return tool, err
	}
	if err := db.DB.First(&tool).Related(&(tool.Status)).Error; err != nil {
		return tool, err
	}
	return tool, nil
}

func Save(clusterName string, item toolModel.Tool) error {
	var cluster clusterModel.Cluster
	if err := db.DB.Where(clusterModel.Cluster{Name: clusterName}).First(&cluster).Error; err != nil {
		return err
	}
	item.ClusterID = cluster.ID
	if err := db.DB.Create(&item).Error; err != nil {
		return err
	}
	return nil
}
