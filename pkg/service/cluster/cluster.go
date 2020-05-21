package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"time"
)

var log = logger.Default

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

func Get(name string) (*clusterModel.Cluster, error) {
	var result clusterModel.Cluster
	err := db.DB.Model(clusterModel.Cluster{}).
		Where(&result).
		First(&result).
		Error
	return &result, err
}

func Save(item *clusterModel.Cluster) error {
	if db.DB.NewRecord(item) {
		item.Status = clusterModel.Status{
			Version:    item.Spec.Version,
			Message:    "",
			Phase:      constant.ClusterWaiting,
			Conditions: nil,
		}
		err := db.DB.Create(&item).Error
		if err != nil {
			return err
		}
		go initCluster(*item)
		return nil
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

func initCluster(c clusterModel.Cluster) {
	ad, err := adm.NewClusterAdm()
	if err != nil {
		log.Fatal(err)
	}
	c.Status.Phase = constant.ClusterWaiting
	err = Save(&c)
	if err != nil {
		log.Debugf("can not save cluster status, msg: %s", err.Error())
	}
	for {
		resp, err := ad.OnInitialize(c)
		if err != nil {
		}
		condition := resp.Status.Conditions[len(resp.Status.Conditions)-1]
		switch condition.Status {
		case constant.ConditionFalse:
			log.Debug("cluster %s init fail, message:%s", c.Name, c.Status.Message)
			return
		case constant.ConditionUnknown:
			log.Debugf("cluster %s init...", c.Name)
		case constant.ConditionTrue:
			log.Debug("cluster %s init success")
			return
		}
		err = Save(&resp)
		if err != nil {
			log.Debugf("can not save cluster status, msg: %s", err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}
