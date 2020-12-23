package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterNodeRepository interface {
	Get(clusterName string, name string) (model.ClusterNode, error)
	Page(num, size int, clusterName string) (int, []model.ClusterNode, error)
	List(clusterName string) ([]model.ClusterNode, error)
	ListByRole(clusterName string, role string) ([]model.ClusterNode, error)
	Save(node *model.ClusterNode) error
	FirstMaster(clusterId string) (model.ClusterNode, error)
	AllMaster(clusterId string) ([]model.ClusterNode, error)
	Delete(id string) error
	BatchSave(nodes []*model.ClusterNode) error
}

func NewClusterNodeRepository() ClusterNodeRepository {
	return &clusterNodeRepository{}
}

type clusterNodeRepository struct{}

func (c clusterNodeRepository) Get(clusterName string, name string) (model.ClusterNode, error) {
	var cluster model.Cluster
	var node model.ClusterNode
	if err := db.DB.
		Where(model.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return node, err
	}
	if err := db.DB.
		Where(model.ClusterNode{ClusterID: cluster.ID, Name: name}).
		First(&node).Error; err != nil {
		return node, err
	}
	return node, nil
}

func (c clusterNodeRepository) Page(num, size int, clusterName string) (int, []model.ClusterNode, error) {
	var total int
	var nodes []model.ClusterNode
	var cluster model.Cluster
	if err := db.DB.
		Where(model.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return 0, nil, err
	}
	if err := db.DB.
		Model(model.ClusterNode{}).
		Where(model.ClusterNode{ClusterID: cluster.ID}).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Preload("Host").
		Preload("Host.Credential").
		Preload("Host.Zone").
		Order("substring_index(name, '-', 2), cast(substring_index(name, '-', -1) as UNSIGNED INTEGER)").
		Find(&nodes).Error; err != nil {
		return 0, nodes, err
	}
	return total, nodes, nil
}

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
		Preload("Host").
		Preload("Host.Credential").
		Preload("Host.Zone").
		Order("substring_index(name, '-', 2), cast(substring_index(name, '-', -1) as UNSIGNED INTEGER)").
		Find(&nodes).Error; err != nil {
		return nodes, err
	}
	return nodes, nil
}

func (c clusterNodeRepository) ListByRole(clusterName string, role string) ([]model.ClusterNode, error) {
	var cluster model.Cluster
	var nodes []model.ClusterNode
	if err := db.DB.
		Where(model.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return nodes, err
	}
	if err := db.DB.
		Where(model.ClusterNode{ClusterID: cluster.ID, Role: role}).
		Find(&nodes).Error; err != nil {
		return nodes, err
	}
	return nodes, nil

}

func (c clusterNodeRepository) FirstMaster(clusterId string) (model.ClusterNode, error) {
	var master model.ClusterNode
	if err := db.DB.
		Where(model.ClusterNode{ClusterID: clusterId, Role: constant.NodeRoleNameMaster}).
		Preload("Host").
		Preload("Host.Credential").
		First(&master).
		Error; err != nil {
		return master, err
	}
	return master, nil
}

func (c clusterNodeRepository) AllMaster(clusterId string) ([]model.ClusterNode, error) {
	var masters []model.ClusterNode
	if err := db.DB.
		Where(model.ClusterNode{ClusterID: clusterId, Role: constant.NodeRoleNameMaster}).
		Preload("Host").
		Preload("Host.Credential").
		Find(&masters).
		Error; err != nil {
		return nil, err
	}
	return masters, nil
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
	node.Host.ClusterID = ""
	if err := tx.Save(node.Host).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c clusterNodeRepository) Save(node *model.ClusterNode) error {
	if db.DB.NewRecord(node) {
		if err := db.DB.Create(node).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(node).Error; err != nil {
			return nil
		}
	}
	return nil
}

func (c clusterNodeRepository) BatchSave(nodes []*model.ClusterNode) error {
	tx := db.DB.Begin()
	for i := range nodes {
		if db.DB.NewRecord(nodes[i]) {
			if err := db.DB.Create(nodes[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := db.DB.Save(nodes[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	return nil
}
