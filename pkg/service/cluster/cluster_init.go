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
	go InitCluster(c)
	return nil
}

func InitCluster(c clusterModel.Cluster) {
	ad, err := adm.NewClusterAdm()
	if err != nil {
		log.Println(err)
	}
	db.DB.
		First(&c.Status).
		Order("last_probe_time asc").
		Related(&c.Status.Conditions)
	nodes, err := GetClusterNodes(c.Name)
	if err != nil {
		log.Println(err)
	}
	c.Nodes = nodes
	if c.Status.Phase == constant.ClusterInitializing {
		return
	}
	c.Status.Phase = constant.ClusterInitializing
	err = db.DB.Save(&c.Status).Error
	if err != nil {
		log.Printf("can not save cluster status, msg: %s", err.Error())
	}
	for {
		resp, err := ad.OnInitialize(c)
		if err != nil {
			log.Fatal(err)
		}
		finished := false
		current := resp.Status.Conditions[len(resp.Status.Conditions)-1]
		switch current.Status {
		case constant.ConditionFalse:
			log.Printf("cluster %s initial fail, message:%s", c.Name, c.Status.Message)
			resp.Status.Phase = constant.ClusterFailed
			finished = true
		case constant.ConditionUnknown:
			log.Printf("cluster %s initial...", c.Name)
		case constant.ConditionTrue:
			log.Printf("cluster %s initial success", c.Name)
			finished = true
		}
		c.Status = resp.Status
		err = db.DB.Save(&c.Status).Error
		if err != nil {
			log.Println(err.Error())
		}
		for i, _ := range c.Status.Conditions {
			c.Status.Conditions[i].StatusID = c.Status.ID
			if db.DB.NewRecord(c.Status.Conditions[i]) {
				err := db.DB.Create(&c.Status.Conditions[i]).Error
				if err != nil {
					log.Println(err.Error())
				}
			} else {
				err := db.DB.Save(&c.Status.Conditions[i]).Error
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
		if finished {
			return
		}
		time.Sleep(5 * time.Second)
	}
}
