package adm

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"io"
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
	writer io.Writer
	Kobe   kobe.Interface
}

func NewCluster(cluster model.Cluster, writer ...io.Writer) *Cluster {
	c := &Cluster{
		Cluster: cluster,
	}
	if writer != nil {
		c.writer = writer[0]
	}
	c.Kobe = kobe.NewAnsible(&kobe.Config{
		Inventory: c.ParseInventory(),
	})
	for name, _ := range facts.DefaultFacts {
		c.Kobe.SetVar(name, facts.DefaultFacts[name])
	}
	clusterVars := cluster.GetKobeVars()
	for k, v := range clusterVars {
		c.Kobe.SetVar(k, v)
	}
	c.Kobe.SetVar(facts.ClusterNameFactName, cluster.Name)
	repo := repository.NewSystemSettingRepository()
	val, _ := repo.Get("ip")
	c.Kobe.SetVar(facts.LocalHostnameFactName, val.Value)
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
		ca.EnsureInitNetwork,
		ca.EnsureInitHelm,
		ca.EnsureInitNpd,
		ca.EnsureInitMetricsServer,
		ca.EnsureInitIngressController,
		ca.EnsurePostInit,
	}
	return ca
}

func (ca *ClusterAdm) OnInitialize(c Cluster) (Cluster, error) {
	err := ca.Create(&c)
	return c, err
}
