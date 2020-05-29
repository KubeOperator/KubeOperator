package adm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	phase "github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/base"
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

func (ca *ClusterAdm) EnsureSystemConfig(c *Cluster) error {
	ph := phase.SystemConfigPhase{}
	resp, err := ph.Run(c.Kobe)
	if err != nil {
		return err
	}
	if resp.HostFailedInfo != nil && len(resp.HostFailedInfo) > 0 {
		by, _ := json.Marshal(resp.HostFailedInfo)
		return errors.New(string(by))
	}
	return nil
}
