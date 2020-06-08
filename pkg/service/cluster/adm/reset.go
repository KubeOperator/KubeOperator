package adm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/reset"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"time"
)

func (ca *ClusterAdm) Reset(c *Cluster) error {
	condition, err := ca.getResetCurrentCondition(c)
	if err != nil {
		return err
	}
	now := time.Now()
	f := ca.getCreateHandler(condition.Name)
	if f == nil {
		return fmt.Errorf("can't get handler by %s", condition.Name)
	}
	resp, err := f(c)
	if err != nil {
		c.setCondition(clusterModel.Condition{
			Name:          condition.Name,
			Status:        constant.ConditionFalse,
			LastProbeTime: now,
			Message:       err.Error(),
		})
		c.Status.Phase = constant.ClusterFailed
		c.Status.Message = err.Error()
		return nil
	}
	if resp.HostFailedInfo != nil && len(resp.HostFailedInfo) > 0 {
		by, _ := json.Marshal(resp.HostFailedInfo)
		c.setCondition(clusterModel.Condition{
			Name:          condition.Name,
			Status:        constant.ConditionFalse,
			LastProbeTime: now,
			Message:       string(by),
		})
		c.Status.Phase = constant.ClusterFailed
		c.Status.Message = string(by)
		return nil
	}
	c.setCondition(clusterModel.Condition{
		Name:          condition.Name,
		Status:        constant.ConditionTrue,
		LastProbeTime: now,
	})

	nextConditionType := ca.getNextConditionName(condition.Name)
	if nextConditionType == ConditionTypeDone {
		c.Status.Phase = constant.ClusterTerminated
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
func (ca *ClusterAdm) getResetCurrentCondition(c *Cluster) (*clusterModel.Condition, error) {
	if c.Status.Phase == constant.ClusterTerminated {
		return nil, errors.New("cluster phase is running now")
	}
	if len(ca.createHandlers) == 0 {
		return nil, errors.New("no create handlers")
	}
	if len(c.Status.Conditions) == 0 {
		return &clusterModel.Condition{
			Name:          ca.createHandlers[0].name(),
			Status:        constant.ConditionUnknown,
			LastProbeTime: time.Now(),
			Message:       "waiting process",
		}, nil
	}
	for _, condition := range c.Status.Conditions {
		if condition.Status == constant.ConditionFalse || condition.Status == constant.ConditionUnknown {
			return &condition, nil
		}
	}
	return nil, errors.New("no condition need process")
}

func (ca *ClusterAdm) EnsureRestTaskStart(c *Cluster) (kobe.Result, error) {
	return kobe.Result{}, nil
}

func (ca *ClusterAdm) EnsureRestCluster(c *Cluster) (kobe.Result, error) {
	phase := reset.ResetClusterPhase{}
	return phase.Run(c.Kobe)
}
