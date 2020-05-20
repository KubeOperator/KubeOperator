package adm

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	ConditionTypeDone = "EnsureDone"
)

type Handler func(*Cluster) error

func (h Handler) name() string {
	name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	i := strings.Index(name, "Ensure")
	if i == -1 {
		return ""
	}
	return strings.TrimSuffix(name[i:], "-fm")
}

type Cluster struct {
	clusterModel.Cluster
	Kobe kobe.Interface
}

func (c *Cluster) setCondition(newCondition clusterModel.Condition) {
	var conditions []clusterModel.Condition
	exist := false
	for _, condition := range c.Status.Conditions {
		if condition.Name == newCondition.Name {
			exist = true
			if newCondition.Status != condition.Status {
				condition.Status = newCondition.Status
			}
			if newCondition.Message != condition.Message {
				condition.Message = newCondition.Message
			}
			if !newCondition.LastProbeTime.IsZero() && newCondition.LastProbeTime != condition.LastProbeTime {
				condition.LastProbeTime = newCondition.LastProbeTime
			}
		}
		conditions = append(conditions, condition)
	}
	if !exist {
		if newCondition.LastProbeTime.IsZero() {
			newCondition.LastProbeTime = time.Now()
		}
		conditions = append(conditions, newCondition)
	}
	c.Status.Conditions = conditions

}

func NewCluster(cluster clusterModel.Cluster) (*Cluster, error) {
	c := &Cluster{
		Cluster: cluster,
	}
	return c, nil
}

type ClusterAdm struct {
	createHandlers []Handler
}

func NewClusterAdm() (*ClusterAdm, error) {
	ca := new(ClusterAdm)
	ca.createHandlers = []Handler{
		ca.EnsureDockerInstall,
	}
	return ca, nil
}

func (ca *ClusterAdm) OnInitialize(cluster clusterModel.Cluster) (clusterModel.Cluster, error) {
	c, err := NewCluster(cluster)
	if err != nil {
		return cluster, err
	}
	err = ca.Create(c)
	return cluster, err
}

func (ca *ClusterAdm) OnJoin(cluster clusterModel.Cluster) (clusterModel.Cluster, error) {
	c, err := NewCluster(cluster)
	if err != nil {
		return cluster, err
	}
	err = ca.Create(c)
	return cluster, err
}

func (ca *ClusterAdm) getCreateHandler(conditionType string) Handler {
	for _, f := range ca.createHandlers {
		if conditionType == f.name() {
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
