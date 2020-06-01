package adm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/initial"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/prepare"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"reflect"
	"runtime"
	"strings"
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
	resp, err := f(c)
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
	if resp.HostFailedInfo != nil && len(resp.HostFailedInfo) > 0 {
		by, _ := json.Marshal(resp.HostFailedInfo)
		c.setCondition(clusterModel.Condition{
			Name:          condition.Name,
			Status:        constant.ConditionFalse,
			LastProbeTime: now,
			Message:       string(by),
		})
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
func (ca *ClusterAdm) getCreateCurrentCondition(c *Cluster) (*clusterModel.Condition, error) {
	if c.Status.Phase == constant.ClusterRunning {
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

func (ca *ClusterAdm) getCreateHandler(conditionName string) Handler {
	for _, f := range ca.createHandlers {
		if conditionName == f.name() {
			return f
		}
	}
	return nil
}
func (ca *ClusterAdm) getNextConditionName(conditionName string) string {
	var (
		i int
		f Handler
	)
	for i, f = range ca.createHandlers {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if strings.Contains(name, conditionName) {
			break
		}
	}
	if i == len(ca.createHandlers)-1 {
		return ConditionTypeDone
	}
	next := ca.createHandlers[i+1]
	return next.name()
}

func (ca *ClusterAdm) EnsurePrepareBaseSystemConfig(c *Cluster) (kobe.Result, error) {
	phase := prepare.BaseSystemConfigPhase{}
	return phase.Run(c.Kobe)
}

func (ca *ClusterAdm) EnsurePrepareContainerRuntime(c *Cluster) (kobe.Result, error) {
	phase := prepare.ContainerRuntimePhase{
		ContainerRuntime: c.Spec.RuntimeType,
	}
	return phase.Run(c.Kobe)
}

func (ca *ClusterAdm) EnsurePrepareKubernetesComponent(c *Cluster) (kobe.Result, error) {
	phase := prepare.KubernetesComponentPhase{}
	return phase.Run(c.Kobe)
}

func (ca *ClusterAdm) EnsurePrepareLoadBalancer(c *Cluster) (kobe.Result, error) {
	phase := prepare.LoadBalancerPhase{}
	return phase.Run(c.Kobe)
}

func (ca *ClusterAdm) EnsurePrepareCertificates(c *Cluster) (kobe.Result, error) {
	phase := prepare.CertificatesPhase{}
	return phase.Run(c.Kobe)
}

func (ca *ClusterAdm) EnsureInitEtcd(c *Cluster) (kobe.Result, error) {
	phase := initial.EtcdPhase{}
	return phase.Run(c.Kobe)
}
func (ca *ClusterAdm) EnsureInitKubeConfig(c *Cluster) (kobe.Result, error) {
	phase := initial.KubeConfigPhase{}
	return phase.Run(c.Kobe)
}
func (ca *ClusterAdm) EnsureInitMaster(c *Cluster) (kobe.Result, error) {
	phase := initial.MasterPhase{}
	return phase.Run(c.Kobe)
}
