package adm

import (
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/backup"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/initial"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/prepare"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/upgrade"
)

func (ca *ClusterAdm) Upgrade(c *Cluster) error {
	condition := ca.getUpgradeCurrentCondition(c)
	if condition != nil {
		now := time.Now()
		f := ca.getUpgradeHandler(condition.Name)
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

		nextConditionType := ca.getNextUpgradeConditionName(condition.Name)
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

func (ca *ClusterAdm) getUpgradeCurrentCondition(c *Cluster) *model.ClusterStatusCondition {
	if len(c.Status.ClusterStatusConditions) == 0 {
		return &model.ClusterStatusCondition{
			Name:          ca.upgradeHandlers[0].name(),
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

func (ca *ClusterAdm) getUpgradeHandler(conditionName string) Handler {
	for _, f := range ca.upgradeHandlers {
		if conditionName == f.name() {
			return f
		}
	}
	return nil
}
func (ca *ClusterAdm) getNextUpgradeConditionName(conditionName string) string {
	var (
		i int
		f Handler
	)
	for i, f = range ca.upgradeHandlers {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if strings.Contains(name, conditionName) {
			break
		}
	}
	if i == len(ca.upgradeHandlers)-1 {
		return ConditionTypeDone
	}
	next := ca.upgradeHandlers[i+1]
	return next.name()
}

func (ca *ClusterAdm) EnsureUpgradeTaskStart(c *Cluster) error {
	time.Sleep(5 * time.Second)
	writeLog("----upgrade task start----", c.writer)
	return nil
}

func (ca *ClusterAdm) EnsureBackupETCD(c *Cluster) error {
	time.Sleep(5 * time.Second)

	phase := backup.BackupClusterPhase{}
	return phase.Run(c.Kobe, c.writer)
}
func (ca *ClusterAdm) EnsureUpgradeRuntime(c *Cluster) error {
	time.Sleep(5 * time.Second)
	phase := prepare.ContainerRuntimePhase{
		Upgrade: true,
	}
	return phase.Run(c.Kobe, c.writer)

}
func (ca *ClusterAdm) EnsureUpgradeETCD(c *Cluster) error {
	time.Sleep(5 * time.Second)
	phase := initial.EtcdPhase{
		Upgrade: true,
	}
	return phase.Run(c.Kobe, c.writer)
}
func (ca *ClusterAdm) EnsureUpgradeKubernetes(c *Cluster) error {
	time.Sleep(5 * time.Second)
	phase := upgrade.UpgradeClusterPhase{
		Version: c.Spec.UpgradeVersion,
	}
	return phase.Run(c.Kobe, c.writer)

}
func (ca *ClusterAdm) EnsureUpdateCertificates(c *Cluster) error {
	time.Sleep(5 * time.Second)
	phase := prepare.CertificatesPhase{}
	return phase.Run(c.Kobe, c.writer)

}
