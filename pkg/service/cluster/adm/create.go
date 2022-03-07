package adm

import (
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/initial"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/plugin/ingress"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/prepare"
)

var log = logger.Default

func (ca *ClusterAdm) Create(c *Cluster) error {
	condition := ca.getCreateCurrentCondition(c)
	if condition != nil {
		now := time.Now()
		f := ca.getCreateHandler(condition.Name)
		err := f(c)
		if err != nil {
			c.setCondition(model.ClusterStatusCondition{
				Name:          condition.Name,
				Status:        constant.ConditionFalse,
				LastProbeTime: now,
				Message:       err.Error(),
			})
			c.Status.Phase = constant.ClusterFailed
			c.Status.Message = err.Error()
			return nil
		}
		c.setCondition(model.ClusterStatusCondition{
			Name:          condition.Name,
			Status:        constant.ConditionTrue,
			LastProbeTime: now,
		})

		nextConditionType := ca.getNextCreateConditionName(condition.Name)
		if nextConditionType == ConditionTypeDone {
			c.Status.Phase = constant.ClusterRunning
		} else {
			c.setCondition(model.ClusterStatusCondition{
				Name:          nextConditionType,
				Status:        constant.ConditionUnknown,
				LastProbeTime: time.Now(),
				Message:       "",
			})
		}
	}
	return nil
}

func (ca *ClusterAdm) getCreateCurrentCondition(c *Cluster) *model.ClusterStatusCondition {
	if len(c.Status.ClusterStatusConditions) == 0 {
		return &model.ClusterStatusCondition{
			Name:          ca.createHandlers[0].name(),
			Status:        constant.ConditionUnknown,
			LastProbeTime: time.Now(),
			Message:       "",
		}
	}
	for _, condition := range c.Status.ClusterStatusConditions {
		if condition.Status == constant.ConditionFalse || condition.Status == constant.ConditionUnknown {
			return &condition
		}
	}
	return nil
}

func (ca *ClusterAdm) getCreateHandler(conditionName string) Handler {
	for _, f := range ca.createHandlers {
		if conditionName == f.name() {
			return f
		}
	}
	return nil
}
func (ca *ClusterAdm) getNextCreateConditionName(conditionName string) string {
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

func (ca *ClusterAdm) EnsureInitTaskStart(c *Cluster) error {
	time.Sleep(5 * time.Second)
	writeLog("----init task start----", c.FileName)
	return nil
}

func (ca *ClusterAdm) EnsurePrepareBaseSystemConfig(c *Cluster) error {
	phase := prepare.BaseSystemConfigPhase{}
	err := phase.Run(c.Kobe, c.FileName)
	return err
}

func (ca *ClusterAdm) EnsurePrepareContainerRuntime(c *Cluster) error {
	phase := prepare.ContainerRuntimePhase{}
	return phase.Run(c.Kobe, c.FileName)
}

func (ca *ClusterAdm) EnsurePrepareKubernetesComponent(c *Cluster) error {
	phase := prepare.KubernetesComponentPhase{}
	return phase.Run(c.Kobe, c.FileName)
}

func (ca *ClusterAdm) EnsurePrepareLoadBalancer(c *Cluster) error {
	phase := prepare.LoadBalancerPhase{}
	return phase.Run(c.Kobe, c.FileName)
}

func (ca *ClusterAdm) EnsurePrepareCertificates(c *Cluster) error {
	phase := prepare.CertificatesPhase{}
	return phase.Run(c.Kobe, c.FileName)
}

func (ca *ClusterAdm) EnsureInitEtcd(c *Cluster) error {
	phase := initial.EtcdPhase{}
	return phase.Run(c.Kobe, c.FileName)
}

func (ca *ClusterAdm) EnsureInitMaster(c *Cluster) error {
	phase := initial.MasterPhase{}
	return phase.Run(c.Kobe, c.FileName)
}

func (ca *ClusterAdm) EnsureInitWorker(c *Cluster) error {
	phase := initial.WorkerPhase{}
	return phase.Run(c.Kobe, c.FileName)
}
func (ca *ClusterAdm) EnsureInitNetwork(c *Cluster) error {
	phase := initial.NetworkPhase{}
	return phase.Run(c.Kobe, c.FileName)
}

func (ca *ClusterAdm) EnsureInitHelm(c *Cluster) error {
	phase := initial.HelmPhase{}
	return phase.Run(c.Kobe, c.FileName)
}

func (ca *ClusterAdm) EnsureInitMetricsServer(c *Cluster) error {
	phase := initial.MetricsServerPhase{}
	return phase.Run(c.Kobe, c.FileName)
}

func (ca *ClusterAdm) EnsureInitIngressController(c *Cluster) error {
	phase := ingress.ControllerPhase{}
	return phase.Run(c.Kobe, c.FileName)
}

func (ca *ClusterAdm) EnsurePostInit(c *Cluster) error {
	phase := initial.PostPhase{}
	return phase.Run(c.Kobe, c.FileName)
}
