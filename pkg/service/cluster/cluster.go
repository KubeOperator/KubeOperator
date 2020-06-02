package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
)

func Page(num, size int) (clusters []clusterModel.Cluster, total int, err error) {
	err = db.DB.Model(clusterModel.Cluster{}).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Preload("Status").
		Preload("Spec").
		Find(&clusters).
		Error
	return
}

func List() (clusters []clusterModel.Cluster, err error) {
	err = db.DB.Model(clusterModel.Cluster{}).
		Preload("Spec").
		Preload("Status").
		Find(&clusters).Error
	return
}

func Get(name string) (clusterModel.Cluster, error) {
	var result clusterModel.Cluster
	result.Name = name
	if err := db.DB.Where(result).First(&result).Error; err != nil {
		return result, err
	}
	if err := db.DB.First(&result).
		Related(&result.Spec).
		Related(&result.Status).Error; err != nil {
		return result, err
	}
	return result, nil
}

func Save(item *clusterModel.Cluster) error {
	if err := db.DB.Create(&item).Error; err != nil {
		return err
	}
	//go InitCluster(*item)
	return nil
}

func Delete(name string) error {
	var cluster clusterModel.Cluster
	if err := db.DB.Where(clusterModel.Cluster{Name: name}).
		First(&cluster).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&cluster).Error; err != nil {
		return err
	}
	return nil
}

func Batch(operation string, items []clusterModel.Cluster) ([]clusterModel.Cluster, error) {
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, c := range items {
			if err := Delete(c.Name); err != nil {
				tx.Rollback()
				return nil, err
			}

		}
		tx.Commit()
	default:
		return nil, constant.NotSupportedBatchOperation
	}
	return items, nil
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

func GetClusterStatus(clusterName string) (clusterModel.Status, error) {
	var cluster clusterModel.Cluster
	var status clusterModel.Status
	if err := db.DB.
		Where(&clusterModel.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return status, err
	}
	if err := db.DB.
		Where(status).
		First(&status).Error; err != nil {
		return status, err
	}
	if err := db.DB.
		First(&status).
		Order("last_probe_time asc").
		Related(&status.Conditions).Error; err != nil {
		return status, err
	}
	return status, nil
}
