package cluster

import (
	"ko3-gin/pkg/constant"
	"ko3-gin/pkg/db"
	clusterModel "ko3-gin/pkg/model/cluster"
	"ko3-gin/pkg/model/common"
)

func Page(num, size int) (clusters []clusterModel.Cluster, total int, err error) {
	err = db.DB.Model(clusterModel.Cluster{}).
		Find(&clusters).
		Offset((num - 1) * size).
		Limit(size).
		Count(&total).Error
	return
}

func List() (clusters []clusterModel.Cluster, err error) {
	err = db.DB.Model(clusterModel.Cluster{}).Find(&clusters).Error
	return
}

func Get(name string) (cluster clusterModel.Cluster, err error) {
	err = db.DB.Model(clusterModel.Cluster{}).
		Where(&clusterModel.Cluster{
			BaseModel: common.BaseModel{
				Name: name,
			},
		}).First(&cluster).Error
	return
}

func Save(item *clusterModel.Cluster) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}

func Delete(name string) error {
	var c clusterModel.Cluster
	c.Name = name
	return db.DB.Delete(&c).Error
}

func Batch(operation string, items []clusterModel.Cluster) ([]clusterModel.Cluster, error) {
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, item := range items {
			err := db.DB.Model(clusterModel.Cluster{}).Delete(&item).Error
			if err != nil {
				tx.Rollback()
			}
		}
		tx.Commit()
	default:
		return nil, constant.NotSupportedBatchOperation
	}
	return items, nil
}

func Nodes(clusterName string) (nodes []clusterModel.Node, err error) {
	err = db.DB.Model(clusterModel.Node{}).
		Where(&clusterModel.Node{ClusterID: clusterName}).
		Find(&nodes).Error
	return
}

func DeleteNode(clusterName, nodeName string) error {
	cluster, err := Get(clusterName)
	if err != nil {
		return err
	}
	var node clusterModel.Node
	node.Name = nodeName
	node.ClusterID = cluster.ID
	return db.DB.Delete(&node).Error
}

func AddNode(clusterName string, node *clusterModel.Node) error {
	cluster, err := Get(clusterName)
	if err != nil {
		return err
	}
	node.ClusterID = cluster.ID
	return db.DB.Create(&cluster).Error
}
