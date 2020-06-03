package cluster

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"log"
	"time"
)

func RetryInitCluster(c clusterModel.Cluster) error {
	if c.Status.Phase != constant.ClusterFailed {
		return errors.New("cluster status is not failed")
	}
	db.DB.
		First(&c.Status).
		Order("last_probe_time asc").
		Related(&c.Status.Conditions)
	if len(c.Status.Conditions) > 0 {
		for i, _ := range c.Status.Conditions {
			if c.Status.Conditions[i].Status == constant.ConditionFalse {
				c.Status.Conditions[i].Status = constant.ConditionUnknown
				c.Status.Conditions[i].Message = ""
			}
			err := db.DB.Save(&c.Status.Conditions[i]).Error
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	return InitCluster(c)
}

func InitCluster(c clusterModel.Cluster) error {
	status, err := GetClusterStatus(c.Name)
	if err != nil {
		return err
	}
	c.Status = status
	nodes, err := GetClusterNodes(c.Name)
	if err != nil {
		return err
	}
	c.Nodes = nodes
	c.Status.Phase = constant.ClusterInitializing
	if err := SaveClusterStatus(&(c.Status)); err != nil {
		return err
	}
	statusChan := make(chan *clusterModel.Status, 0)
	stopChan := make(chan int, 0)
	go DoInitCluster(c, statusChan, stopChan)
	go SyncStatus(statusChan, stopChan)
	return nil
}

func DoInitCluster(c clusterModel.Cluster, statusChan chan *clusterModel.Status, stopChan chan int) {
	ad, _ := adm.NewClusterAdm()
	for {
		resp, err := ad.OnInitialize(c)
		if err != nil {
			log.Fatal(err)
		}
		c.Status = resp.Status
		select {
		case <-stopChan:
			return
		default:
			statusChan <- &(c.Status)
		}
		time.Sleep(5 * time.Second)
	}
}

func SyncStatus(statusChan chan *clusterModel.Status, stopChan chan int) {
	for {
		status := <-statusChan
		if err := db.DB.Save(status).Error; err != nil {
			stopChan <- 1
			return
		}
		switch status.Phase {
		case constant.ClusterFailed, constant.ClusterRunning:
			stopChan <- 1
			return
		}
	}
}
