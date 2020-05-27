package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
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
	if db.DB.NewRecord(item) {
		tx := db.DB.Begin()
		if err := db.DB.Create(&item).Error; err != nil {
			tx.Rollback()
			return err
		}
		item.Spec.ClusterID = item.ID
		if err := db.DB.Create(&item.Spec).Error; err != nil {
			tx.Rollback()
			return err
		}
		item.Status = clusterModel.Status{
			ClusterID: item.ID,
			Message:   "",
			Phase:     constant.ClusterWaiting,
		}
		if err := db.DB.Create(&item.Status).Error; err != nil {
			tx.Rollback()
			return err
		}
		for _, node := range item.Nodes {
			node.ClusterID = item.ID
			if err := db.DB.First(&node.Host).Error; err != nil {
				tx.Rollback()
				return err
			}
			node.HostID = node.Host.ID
			if err := db.DB.Create(&node).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
		go InitCluster(*item)
		return nil
	} else {
		return db.DB.Save(&item).Error
	}
}

func Delete(name string) error {
	tx := db.DB.Begin()
	c := clusterModel.Cluster{Name: name,}
	if err := db.DB.First(&c).Delete(&c).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := db.DB.Where(clusterModel.Spec{ClusterID: c.ID,}).
		Delete(clusterModel.Spec{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := db.DB.Where(clusterModel.Status{ClusterID: c.ID,}).
		Delete(clusterModel.Status{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := db.DB.Where(clusterModel.Node{ClusterID: c.ID,}).
		Delete(clusterModel.Node{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func Batch(operation string, items []clusterModel.Cluster) ([]clusterModel.Cluster, error) {
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, c := range items {
			if err := db.DB.First(&c).Delete(&c).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			if err := db.DB.Where(clusterModel.Spec{ClusterID: c.ID,}).
				Delete(clusterModel.Spec{}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			if err := db.DB.Where(clusterModel.Status{ClusterID: c.ID,}).
				Delete(clusterModel.Status{}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			if err := db.DB.Where(clusterModel.Node{ClusterID: c.ID,}).
				Delete(clusterModel.Node{}).Error; err != nil {
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

func GetStatus(clusterName string) (clusterModel.Status, error) {
	var cluster clusterModel.Cluster
	var status clusterModel.Status
	if err := db.DB.
		Where(&clusterModel.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return status, err
	}
	status.ClusterID = cluster.ID
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
