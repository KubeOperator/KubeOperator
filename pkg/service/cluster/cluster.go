package cluster

import (
	"ko3-gin/pkg/db"
	clusterModel "ko3-gin/pkg/model/cluster"
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
	err := db.DB.First(&c).Error
	if err != nil {
		return err
	}
	return db.DB.Delete(&c).Error
}
