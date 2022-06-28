package adm

import (
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/initial"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/plugin/ingress"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/prepare"
)

func (ca *ClusterAdm) Create(aHelper *AnsibleHelper) error {
	task := ca.getCreateCurrentTask(aHelper)
	if task != nil {
		f := ca.getCreateHandler(task.Task)
		if err := f(aHelper); err != nil {
			aHelper.setCondition(model.TaskLogDetail{
				Task:          task.Task,
				Status:        constant.TaskLogStatusFailed,
				LastProbeTime: time.Now().Unix(),
				StartTime:     task.StartTime,
				EndTime:       time.Now().Unix(),
				Message:       err.Error(),
			})
			aHelper.Status = constant.TaskLogStatusFailed
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

		nextConditionType := ca.getNextCreateConditionName(task.Task)
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

func (ca *ClusterAdm) getCreateCurrentTask(aHelper *AnsibleHelper) *model.TaskLogDetail {
	if len(aHelper.LogDetail) == 0 {
		taskItem := &model.TaskLogDetail{
			Task:          ca.createHandlers[0].name(),
			Status:        constant.TaskLogStatusRunning,
			LastProbeTime: time.Now().Unix(),
			StartTime:     time.Now().Unix(),
			EndTime:       time.Now().Unix(),
			Message:       "",
		}
		aHelper.LogDetail = append(aHelper.LogDetail, *taskItem)
		return taskItem
	}
	for _, detail := range aHelper.LogDetail {
		if detail.Status == constant.TaskLogStatusFailed || detail.Status == constant.TaskLogStatusRunning {
			return &detail
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

func (ca *ClusterAdm) EnsureInitTaskStart(aHelper *AnsibleHelper) error {
	time.Sleep(5 * time.Second)
	writeLog("----init task start----", aHelper.Writer)
	return nil
}

func (ca *ClusterAdm) EnsurePrepareBaseSystemConfig(aHelper *AnsibleHelper) error {
	phase := prepare.BaseSystemConfigPhase{}
	err := phase.Run(aHelper.Kobe, aHelper.Writer)
	return err
}

func (ca *ClusterAdm) EnsurePrepareContainerRuntime(aHelper *AnsibleHelper) error {
	phase := prepare.ContainerRuntimePhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsurePrepareKubernetesComponent(aHelper *AnsibleHelper) error {
	phase := prepare.KubernetesComponentPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsurePrepareLoadBalancer(aHelper *AnsibleHelper) error {
	phase := prepare.LoadBalancerPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsurePrepareCertificates(aHelper *AnsibleHelper) error {
	phase := prepare.CertificatesPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureInitEtcd(aHelper *AnsibleHelper) error {
	phase := initial.EtcdPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureInitMaster(aHelper *AnsibleHelper) error {
	phase := initial.MasterPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureInitWorker(aHelper *AnsibleHelper) error {
	phase := initial.WorkerPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}
func (ca *ClusterAdm) EnsureInitNetwork(aHelper *AnsibleHelper) error {
	phase := initial.NetworkPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureInitHelm(aHelper *AnsibleHelper) error {
	phase := initial.HelmPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureInitMetricsServer(aHelper *AnsibleHelper) error {
	phase := initial.MetricsServerPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureInitIngressController(aHelper *AnsibleHelper) error {
	phase := ingress.ControllerPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsurePostInit(aHelper *AnsibleHelper) error {
	phase := initial.PostPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}
