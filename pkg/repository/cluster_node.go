package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterNodeRepository interface {
	List(clusterName string) ([]model.ClusterNode, error)
	FistMaster(ClusterId string) (model.ClusterNode, error)
	Delete(id string) error
}

func NewClusterNodeRepository() ClusterNodeRepository {
	return &clusterNodeRepository{}
}

type clusterNodeRepository struct{}

func (c clusterNodeRepository) List(clusterName string) ([]model.ClusterNode, error) {
	var cluster model.Cluster
	var nodes []model.ClusterNode
	if err := db.DB.
		Where(model.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return nodes, err
	}
	if err := db.DB.
		Where(model.ClusterNode{ClusterID: cluster.ID}).
		Find(&nodes).Error; err != nil {
		return nodes, err
	}
	return nodes, nil
}

func (c clusterNodeRepository) FistMaster(ClusterId string) (model.ClusterNode, error) {
	var master model.ClusterNode
	if err := db.DB.
		Where(model.ClusterNode{ClusterID: ClusterId, Role: constant.NodeRoleNameMaster}).
		Preload("Host").
		First(&master).
		Error; err != nil {
		return master, err
	}
	return master, nil
}

func (c clusterNodeRepository) Delete(id string) error {
	node := model.ClusterNode{ID: id}
	tx := db.DB.Begin()
	if err := db.DB.
		First(&node).
		Related(&node.Host).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&node).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
