package cluster

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"log"
	"time"
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

func Get(name string) (*clusterModel.Cluster, error) {
	var result clusterModel.Cluster
	result.Name = name
	err := db.DB.First(&result).
		Related(&result.Spec).
		Related(&result.Status).Error
	return &result, err
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
		//go initCluster(*item)
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

func InitCluster(c clusterModel.Cluster) {
	ad, err := adm.NewClusterAdm()
	if err != nil {
		log.Println(err)
	}
	c.Status.Phase = constant.ClusterInitializing
	err = db.DB.Save(&c.Status).Error
	if err != nil {
		log.Printf("can not save cluster status, msg: %s", err.Error())
	}
	c.Status.Conditions = []clusterModel.Condition{}
	for {
		resp, err := ad.OnInitialize(c)
		if err != nil {
			log.Fatal(err)
		}
		finished := false
		condition := resp.Status.Conditions[len(resp.Status.Conditions)-1]
		switch condition.Status {
		case constant.ConditionFalse:
			log.Printf("cluster %s init fail, message:%s", c.Name, c.Status.Message)
			finished = true
		case constant.ConditionUnknown:
			log.Printf("cluster %s init...", c.Name)
		case constant.ConditionTrue:
			log.Printf("cluster %s init success", c.Name)
			finished = true
		}
		c.Status = resp.Status
		for _, c := range c.Status.Conditions {
			fmt.Println(c.Status)
		}
		if finished {
			return
		}
		time.Sleep(5 * time.Second)
	}
}
