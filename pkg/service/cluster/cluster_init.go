package cluster

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
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
	err = c.SetSecret(clusterModel.Secret{
		KubeadmToken:    "abcdefg",
		KubernetesToken: "",
	})
	if err != nil {
		return err
	}
	if status.Phase == constant.ClusterRunning || status.Phase == constant.ClusterInitializing {
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
	c.Status.Phase = constant.ClusterInitializing
	if err := SaveClusterStatus(&(c.Status)); err != nil {
		return err
	}
	statusChan := make(chan *clusterModel.Status, 0)
	stopChan := make(chan int, 0)
	go DoInitCluster(c, statusChan, stopChan)
	go SyncStatus(statusChan, stopChan)
	err = GetAndSaveClusterApiToken(c)
	if err != nil {
		log.Println(err.Error())
	}
	return nil
}

func DoInitCluster(c clusterModel.Cluster, statusChan chan *clusterModel.Status, stopChan chan int) {
	ad, _ := adm.NewClusterAdm()
	for {
		resp, err := ad.OnInitialize(c)
		if err != nil {
			c.Status.Message = err.Error()
		}
		c.Status = resp.Status
		if c.Status.Phase == constant.ClusterRunning {
			err := GetAndSaveClusterApiToken(c)
			if err != nil {
				c.Status.Message = err.Error()
			}
		}
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

func GetAndSaveClusterApiToken(c clusterModel.Cluster) error {
	secret, err := GetClusterSecret(c.Name)
	if err != nil {
		return err
	}
	master := c.FistMaster()
	sshConfig := master.ToSSHConfig()
	client, err := ssh.New(&sshConfig)
	if err != nil {
		return err
	}
	token, err := clusterUtil.GetClusterToken(client)
	if err != nil {
		return err
	}
	secret.KubernetesToken = token
	err = c.SetSecret(secret)
	if err != nil {
		return err
	}
	return nil
}
