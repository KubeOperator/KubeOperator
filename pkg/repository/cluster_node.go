package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterNodeRepository interface {
	List(clusterName string) ([]model.ClusterNode, error)
	FistMaster(ClusterId string) (model.ClusterNode, error)
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
		First(&master).Error; err != nil {
		return master, err
	}
	return master, nil
}
