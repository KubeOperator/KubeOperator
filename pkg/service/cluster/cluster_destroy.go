package cluster

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"time"
)

func DestroyCluster(c clusterModel.Cluster) error {
	status, err := GetClusterStatus(c.Name)
	if err != nil {
		return err
	}
	if status.Phase == constant.ClusterInitializing || status.Phase == constant.ClusterTerminating {
		return errors.New(fmt.Sprintf("invalid status: %s", status.Phase))
	}
	c.Status = status
	nodes, err := GetClusterNodes(c.Name)
	if err != nil {
		return err
	}
	if len(nodes) < 1 {
		return errors.New(fmt.Sprintf("node size is : %d", len(nodes)))
	}
	c.Nodes = nodes
	c.Status.Phase = constant.ClusterTerminating
	if err := SaveClusterStatus(&(c.Status)); err != nil {
		return err
	}
	statusChan := make(chan *clusterModel.Cluster, 0)
	stopChan := make(chan int, 0)
	go DoDestroyCluster(c, statusChan, stopChan)
	go SyncDestroyStatus(statusChan, stopChan)
	return nil
}

func DoDestroyCluster(c clusterModel.Cluster, statusChan chan *clusterModel.Cluster, stopChan chan int) {
	ad, _ := adm.NewClusterAdm()
	for {
		resp, err := ad.OnReset(c)
		if err != nil {
			c.Status.Message = err.Error()
		}
		c.Status = resp.Status
		select {
		case <-stopChan:
			return
		default:
			statusChan <- &c
		}
		time.Sleep(5 * time.Second)
	}
}

func SyncDestroyStatus(statusChan chan *clusterModel.Cluster, stopChan chan int) {
	for {
		c := <-statusChan
		if err := db.DB.Save(&(c.Status)).Error; err != nil {
			stopChan <- 1
			return
		}
		switch c.Status.Phase {
		case constant.ClusterFailed:
			stopChan <- 1
			return
		case constant.ClusterTerminated:
			stopChan <- 1
			db.DB.Delete(&c)
			return
		}
	}
}
