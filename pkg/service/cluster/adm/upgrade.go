package adm

import (
	"fmt"
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
	"github.com/KubeOperator/KubeOperator/pkg/util/version"
)

func (ca *ClusterAdm) Upgrade(aHelper *AnsibleHelper) error {
	task := ca.getUpgradeCurrentTask(aHelper)
	if task != nil {
		f := ca.getUpgradeHandler(task.Task)
		err := f(aHelper)
		if err != nil {
			aHelper.setCondition(model.TaskLogDetail{
				Task:          task.Task,
				Status:        constant.TaskLogStatusFailed,
				LastProbeTime: time.Now().Unix(),
				StartTime:     task.StartTime,
				EndTime:       time.Now().Unix(),
				Message:       err.Error(),
			})
			aHelper.Status = constant.StatusFailed
			aHelper.Message = err.Error()
			return nil
		}
		aHelper.setCondition(model.TaskLogDetail{
			Task:          task.Task,
			Status:        constant.TaskLogStatusSuccess,
			LastProbeTime: time.Now().Unix(),
			StartTime:     task.StartTime,
			EndTime:       time.Now().Unix(),
		})

		nextConditionType := ca.getNextUpgradeConditionName(task.Task)
		if nextConditionType == ConditionTypeDone {
			aHelper.Status = constant.TaskLogStatusSuccess
		} else {
			aHelper.setCondition(model.TaskLogDetail{
				Task:          nextConditionType,
				Status:        constant.TaskLogStatusRunning,
				LastProbeTime: time.Now().Unix(),
				StartTime:     time.Now().Unix(),
			})
		}
	}
	return nil
}

func (ca *ClusterAdm) getUpgradeCurrentTask(aHelper *AnsibleHelper) *model.TaskLogDetail {
	if len(aHelper.LogDetail) == 0 {
		return &model.TaskLogDetail{
			Task:          ca.upgradeHandlers[0].name(),
			Status:        constant.TaskLogStatusRunning,
			LastProbeTime: time.Now().Unix(),
			StartTime:     time.Now().Unix(),
			EndTime:       time.Now().Unix(),
			Message:       "",
		}
	}
	for _, task := range aHelper.LogDetail {
		if task.Status == constant.TaskLogStatusFailed || task.Status == constant.TaskLogStatusRunning {
			return &task
		}
	}
	return nil
}

func (ca *ClusterAdm) getUpgradeHandler(taskName string) Handler {
	for _, f := range ca.upgradeHandlers {
		if taskName == f.name() {
			return f
		}
	}
	return nil
}
func (ca *ClusterAdm) getNextUpgradeConditionName(taskName string) string {
	var (
		i int
		f Handler
	)
	for i, f = range ca.upgradeHandlers {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if strings.Contains(name, taskName) {
			break
		}
	}
	if i == len(ca.upgradeHandlers)-1 {
		return ConditionTypeDone
	}
	next := ca.upgradeHandlers[i+1]
	return next.name()
}

func (ca *ClusterAdm) EnsureUpgradeTaskStart(aHelper *AnsibleHelper) error {
	time.Sleep(5 * time.Second)
	writeLog("----upgrade task start----", aHelper.Writer)
	return nil
}

func (ca *ClusterAdm) EnsureBackupETCD(aHelper *AnsibleHelper) error {
	time.Sleep(5 * time.Second)
	phase := backup.BackupClusterPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}
func (ca *ClusterAdm) EnsureUpgradeRuntime(aHelper *AnsibleHelper) error {
	time.Sleep(5 * time.Second)
	phase := prepare.ContainerRuntimePhase{
		Upgrade: true,
	}
	oldManiFest, _ := GetManiFestBy(aHelper.ClusterVersion)
	newManiFest, _ := GetManiFestBy(aHelper.ClusterUpgradeVersion)
	oldVars := oldManiFest.GetVars()
	newVars := newManiFest.GetVars()
	var runtimeVersionKey = "runtime_version"
	switch aHelper.ClusterRuntime {
	case "docker":
		runtimeVersionKey = strings.Replace(runtimeVersionKey, "runtime", "docker", -1)
	case "containerd":
		runtimeVersionKey = strings.Replace(runtimeVersionKey, "runtime", "containerd", -1)
	}
	oldVersion := oldVars[runtimeVersionKey]
	newVersion := newVars[runtimeVersionKey]
	_, _ = fmt.Fprintf(aHelper.Writer, "%s -> %s", oldVersion, newVersion)
	newer := version.IsNewerThan(newVersion, oldVersion)
	if !newer {
		_, _ = fmt.Fprintln(aHelper.Writer, "runtime version is newest.skip upgrade")
		return nil
	}
	aHelper.Kobe.SetVar(runtimeVersionKey, newVersion)
	return phase.Run(aHelper.Kobe, aHelper.Writer)

}
func (ca *ClusterAdm) EnsureUpgradeETCD(aHelper *AnsibleHelper) error {
	time.Sleep(5 * time.Second)
	phase := initial.EtcdPhase{
		Upgrade: true,
	}
	oldManiFest, _ := GetManiFestBy(aHelper.ClusterVersion)
	newManiFest, _ := GetManiFestBy(aHelper.ClusterUpgradeVersion)
	oldVars := oldManiFest.GetVars()
	newVars := newManiFest.GetVars()
	var etcdVersionKey = "etcd_version"
	oldVersion := oldVars[etcdVersionKey]
	newVersion := newVars[etcdVersionKey]
	_, _ = fmt.Fprintf(aHelper.Writer, "%s -> %s", oldVersion, newVersion)
	newer := version.IsNewerThan(newVersion, oldVersion)
	if !newer {
		_, _ = fmt.Fprintln(aHelper.Writer, "etcd version is newest.skip upgrade")
		return nil
	}
	aHelper.Kobe.SetVar(etcdVersionKey, newVersion)
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}
func (ca *ClusterAdm) EnsureUpgradeKubernetes(aHelper *AnsibleHelper) error {
	time.Sleep(5 * time.Second)
	index := strings.Index(aHelper.ClusterUpgradeVersion, "-")
	phase := upgrade.UpgradeClusterPhase{
		Version: aHelper.ClusterUpgradeVersion[:index],
	}
	return phase.Run(aHelper.Kobe, aHelper.Writer)

}
func (ca *ClusterAdm) EnsureUpdateCertificates(aHelper *AnsibleHelper) error {
	time.Sleep(5 * time.Second)
	phase := prepare.CertificatesPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)

}
