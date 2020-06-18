package adm

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
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

type Handler func(*Cluster) error

func (h Handler) name() string {
	name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	i := strings.Index(name, "Ensure")
	if i == -1 {
		return ""
	}
	return strings.TrimSuffix(name[i:], "-fm")
}

func (c *Cluster) setCondition(newCondition model.ClusterStatusCondition) {
	var conditions []model.ClusterStatusCondition
	exist := false
	for _, condition := range c.Status.ClusterStatusConditions {
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
	c.Status.ClusterStatusConditions = conditions

}

type Cluster struct {
	model.Cluster
	Kobe kobe.Interface
}

func NewCluster(cluster model.Cluster) *Cluster {
	c := &Cluster{
		Cluster: cluster,
	}
	c.Kobe = kobe.NewAnsible(&kobe.Config{
		Inventory: c.ParseInventory(),
	})
	for name, _ := range facts.DefaultFacts {
		c.Kobe.SetVar(name, facts.DefaultFacts[name])
	}
	return c
}

type ClusterAdm struct {
	createHandlers []Handler
}

func NewClusterAdm() *ClusterAdm {
	ca := new(ClusterAdm)
	ca.createHandlers = []Handler{
		ca.EnsureInitTaskStart,
		ca.EnsurePrepareBaseSystemConfig,
		ca.EnsurePrepareContainerRuntime,
		ca.EnsurePrepareKubernetesComponent,
		ca.EnsurePrepareLoadBalancer,
		ca.EnsurePrepareCertificates,
		ca.EnsureInitEtcd,
		ca.EnsureInitMaster,
		ca.EnsureInitWorker,
		ca.EnsurePostInit,
		ca.EnsureInitNetwork,
	}
	return ca
}

func (ca *ClusterAdm) OnInitialize(c Cluster) (Cluster, error) {
	err := ca.Create(&c)
	return c, err
}
