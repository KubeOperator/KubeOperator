package adm

import (
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	ConditionTypeDone = "EnsureDone"
)

type Handler func(*Cluster) (kobe.Result, error)

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
	c.Kobe = kobe.NewAnsible(&kobe.Config{
		Inventory: c.ParseInventory(),
	})
	// set default vars
	for name, _ := range facts.DefaultFacts {
		c.Kobe.SetVar(name, facts.DefaultFacts[name])
	}

	return c, nil
}

type ClusterAdm struct {
	createHandlers []Handler
}

func NewClusterAdm() (*ClusterAdm, error) {
	ca := new(ClusterAdm)
	ca.createHandlers = []Handler{
		ca.EnsurePrepareBaseSystemConfig,
		ca.EnsurePrepareContainerRuntime,
		ca.EnsurePrepareKubernetesComponent,
		ca.EnsurePrepareLoadBalancer,
		ca.EnsurePrepareCertificates,
		ca.EnsureInitEtcd,
		ca.EnsureInitKubeConfig,
		ca.EnsureInitMaster,
	}
	return ca, nil
}

func (ca *ClusterAdm) OnInitialize(cluster clusterModel.Cluster) (clusterModel.Cluster, error) {
	c, err := NewCluster(cluster)
	if err != nil {
		return cluster, err
	}
	err = ca.Create(c)
	return c.Cluster, err
}
