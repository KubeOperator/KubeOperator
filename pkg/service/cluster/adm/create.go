package adm

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"log"
	"time"
)

func (ca *ClusterAdm) Create(c *Cluster) error {
	condition, err := ca.getCreateCurrentCondition(c)
	if err != nil {
		return err
	}
	now := time.Now()
	f := ca.getCreateHandler(condition.Name)
	if f == nil {
		return fmt.Errorf("can't get handler by %s", condition.Name)
	}
	err = f(c)
	if err != nil {
		c.setCondition(clusterModel.Condition{
			Name:          condition.Name,
			Status:        constant.ConditionFalse,
			LastProbeTime: now,
			Message:       err.Error(),
		})
		c.Status.Message = err.Error()
		return nil
	}

	c.setCondition(clusterModel.Condition{
		Name:          condition.Name,
		Status:        constant.ConditionTrue,
		LastProbeTime: now,
	})

	nextConditionType := ca.getNextConditionName(condition.Name)
	if nextConditionType == ConditionTypeDone {
		c.Status.Phase = constant.ClusterRunning
	} else {
		c.setCondition(clusterModel.Condition{
			Name:          nextConditionType,
			Status:        constant.ConditionUnknown,
			LastProbeTime: time.Now(),
			Message:       "waiting process",
		})

	}
	return nil
}

func (ca *ClusterAdm) EnsureDockerInstall(c *Cluster) error {
	log.Println("install docker...")
	time.Sleep(5 * time.Second)
	return nil
}

func (ca *ClusterAdm) EnsureKubeletInstall(c *Cluster) error {
	log.Println("install kubelet...")
	time.Sleep(5 * time.Second)
	return errors.New("aaabbbccc")
}

func (ca *ClusterAdm) EnsureClusterInit(c *Cluster) error {
	log.Println("install cluster...")
	time.Sleep(5 * time.Second)
	return nil
}
