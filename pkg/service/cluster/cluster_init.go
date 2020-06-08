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
	statusChan := make(chan *clusterModel.Cluster, 0)
	stopChan := make(chan int, 0)
	go DoInitCluster(c, statusChan, stopChan)
	go SyncInitStatus(statusChan, stopChan)
	return nil
}

func DoInitCluster(c clusterModel.Cluster, statusChan chan *clusterModel.Cluster, stopChan chan int) {
	ad, _ := adm.NewClusterAdm()
	for {
		resp, err := ad.OnInitialize(c)
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

func SyncInitStatus(statusChan chan *clusterModel.Cluster, stopChan chan int) {
	for {
		c := <-statusChan
		db.DB.Save(&(c.Status))
		switch c.Status.Phase {
		case constant.ClusterFailed:
			stopChan <- 1
			return
		case constant.ClusterRunning:
			stopChan <- 1
			err := GetAndSaveClusterApiToken(*c)
			if err != nil {
				c.Status.Message = err.Error()
			}
			err = c.Status.ClearConditions()
			if err != nil {
				c.Status.Message = err.Error()
			}
			db.DB.Save(&(c.Status))
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
