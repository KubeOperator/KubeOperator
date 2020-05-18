package cluster

import (
	uuid "github.com/satori/go.uuid"
	"ko3-gin/pkg/db"
	"ko3-gin/pkg/model/cluster"
)

func Page(num, size int) (clusters []cluster.Cluster, total int, err error) {
	err = db.DB.Model(cluster.Cluster{}).
		Find(&clusters).
		Offset((num - 1) * size).
		Limit(size).
		Count(total).Error
	return
}

func List() (clusters []cluster.Cluster, err error) {
	err = db.DB.Model(cluster.Cluster{}).Find(clusters).Error
	return
}

func Save(item cluster.Cluster) (cluster cluster.Cluster, err error) {
	if db.DB.NewRecord(item) {
		item.ID = uuid.NewV4().String()
		err = db.DB.Create(&item).Error
	} else {
		err = db.DB.Save(&item).Error
	}
	return item, err
}

func Delete(name string) (cluster cluster.Cluster, err error) {
	cluster.Name = name
	err = db.DB.First(&cluster).Error
	if err != nil {
		return
	}
	err = db.DB.Delete(cluster).Error
	return cluster, err
}

