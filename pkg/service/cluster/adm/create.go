package adm

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/docker"
	"os"
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
		c.setCondition(cluster.Condition{
			Name:          condition.Name,
			Status:        constant.ConditionFalse,
			LastProbeTime: now,
			Message:       err.Error(),
		})
		c.Status.Message = err.Error()
		return nil
	}

	c.setCondition(cluster.Condition{
		Name:          condition.Name,
		Status:        constant.ConditionTrue,
		LastProbeTime: now,
	})

	nextConditionType := ca.getNextConditionName(condition.Name)
	if nextConditionType == ConditionTypeDone {
		c.Status.Phase = constant.ClusterRunning
	} else {
		c.setCondition(cluster.Condition{
			Name:          nextConditionType,
			Status:        constant.ConditionUnknown,
			LastProbeTime: now,
			Message:       "waiting process",
		})

	}
	return nil
}

func (ca *ClusterAdm) EnsureDockerInstall(c *Cluster) error {
	taskId, err := docker.Install(c.Kobe)
	if err != nil {
		return err
	}
	err = c.Kobe.Watch(os.Stdout, taskId)
	if err != nil {
		return err
	}
	result, err := c.Kobe.GetResult(taskId)
	if err != nil {
		return err
	}
	if result.Success {
		return nil
	}
	return nil
}
