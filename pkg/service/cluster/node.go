package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
)

func ListClusterNodes(clusterName string) ([]clusterModel.Node, error) {
	var cluster clusterModel.Cluster
	var nodes []clusterModel.Node
	if err := db.DB.
		Where(clusterModel.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return nodes, err
	}
	if err := db.DB.
		Where(clusterModel.Node{ClusterID: cluster.ID}).
		Find(&nodes).Error; err != nil {
		return nodes, err
	}
	return nodes, nil
}
func PageClusterNodes(clusterName string, num, size int) (nodes []clusterModel.Node, total int, err error) {
	var cluster clusterModel.Cluster
	if err = db.DB.
		Where(clusterModel.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return
	}
	if err = db.DB.
		Where(clusterModel.Node{ClusterID: cluster.ID}).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Find(&nodes).Error; err != nil {
		return
	}
	return
}

func GetClusterNodes(name string) ([]clusterModel.Node, error) {
	var cluster clusterModel.Cluster
	if err := db.DB.
		Where(clusterModel.Cluster{Name: name}).
		Preload("Nodes").
		First(&cluster).Error; err != nil {
		return nil, err
	}
	for i, _ := range cluster.Nodes {
		if err := db.DB.
			Preload("Credential").
			Where(hostModel.Host{
				NodeID: cluster.Nodes[i].ID,
			}).First(&cluster.Nodes[i].Host).Error; err != nil {
			return nil, err
		}
	}
	return cluster.Nodes, nil
}
