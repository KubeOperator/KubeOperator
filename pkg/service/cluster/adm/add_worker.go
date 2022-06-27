package adm

import (
	"encoding/json"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/initial"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/plugin/storage"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/prepare"
)

func (ca *ClusterAdm) AddWorker(aHelper *AnsibleHelper) error {
	detail := ca.getAddWorkerCurrentTask(aHelper)
	if detail != nil {
		now := time.Now()
		f := ca.getAddWorkerHandler(detail.Task)
		err := f(aHelper)
		if err != nil {
			aHelper.setCondition(model.TaskLogDetail{
				Task:          detail.Task,
				Status:        constant.TaskDetailStatusFalse,
				LastProbeTime: now.Unix(),
				EndTime:       now.Unix(),
				Message:       err.Error(),
			})
			aHelper.Status = constant.TaskLogStatusFailed
			aHelper.Message = err.Error()
			return nil
		}
		aHelper.setCondition(model.TaskLogDetail{
			Task:          detail.Task,
			Status:        constant.TaskDetailStatusTrue,
			LastProbeTime: now.Unix(),
			EndTime:       now.Unix(),
		})

		nextConditionType := ca.getNextAddWorkerConditionName(detail.Task)
		if nextConditionType == ConditionTypeDone {
			aHelper.Status = constant.TaskLogStatusSuccess
		} else {
			aHelper.setCondition(model.TaskLogDetail{
				Task:          nextConditionType,
				Status:        constant.TaskDetailStatusUnknown,
				LastProbeTime: time.Now().Unix(),
				StartTime:     time.Now().Unix(),
				Message:       "",
			})
		}
	}
	return nil
}

func (ca *ClusterAdm) getAddWorkerCurrentTask(aHelper *AnsibleHelper) *model.TaskLogDetail {
	if len(aHelper.LogDetail) == 0 {
		return &model.TaskLogDetail{
			Task:          ca.addWorkerHandlers[0].name(),
			Status:        constant.TaskDetailStatusUnknown,
			LastProbeTime: time.Now().Unix(),
			StartTime:     time.Now().Unix(),
			Message:       "",
		}
	}
	for _, detail := range aHelper.LogDetail {
		if detail.Status == constant.TaskDetailStatusFalse || detail.Status == constant.TaskDetailStatusUnknown {
			return &detail
		}
	}
	return nil
}

func (ca *ClusterAdm) getAddWorkerHandler(detailName string) Handler {
	for _, f := range ca.addWorkerHandlers {
		if detailName == f.name() {
			return f
		}
	}
	return nil
}

func (ca *ClusterAdm) getNextAddWorkerConditionName(detailName string) string {
	var (
		i int
		f Handler
	)
	for i, f = range ca.addWorkerHandlers {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if strings.Contains(name, detailName) {
			break
		}
	}
	if i == len(ca.addWorkerHandlers)-1 {
		return ConditionTypeDone
	}
	next := ca.addWorkerHandlers[i+1]
	return next.name()
}

func (ca *ClusterAdm) EnsureAddWorkerTaskStart(aHelper *AnsibleHelper) error {
	time.Sleep(5 * time.Second)
	writeLog("----add worker task start----", aHelper.Writer)
	return nil
}

func (ca *ClusterAdm) EnsureAddWorkerBaseSystemConfig(aHelper *AnsibleHelper) error {
	phase := prepare.AddWorkerBaseSystemConfigPhase{}
	err := phase.Run(aHelper.Kobe, aHelper.Writer)
	return err
}

func (ca *ClusterAdm) EnsureAddWorkerContainerRuntime(aHelper *AnsibleHelper) error {
	phase := prepare.AddWorkerContainerRuntimePhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerKubernetesComponent(aHelper *AnsibleHelper) error {
	phase := prepare.AddWorkerKubernetesComponentPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerLoadBalancer(aHelper *AnsibleHelper) error {
	phase := prepare.AddWorkerLoadBalancerPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerCertificates(aHelper *AnsibleHelper) error {
	phase := prepare.AddWorkerCertificatesPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerWorker(aHelper *AnsibleHelper) error {
	phase := initial.AddWorkerMasterPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerNetwork(aHelper *AnsibleHelper) error {
	phase := initial.AddWorkerNetworkPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerPost(aHelper *AnsibleHelper) error {
	phase := initial.AddWorkerPostPhase{}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerStorage(aHelper *AnsibleHelper) error {
	var provisoners []model.ClusterStorageProvisioner
	phase := storage.AddWorkerStoragePhase{
		AddWorker:                          true,
		EnableNfsProvisioner:               "disable",
		NfsVersion:                         "v4",
		EnableGfsProvisioner:               "disable",
		EnableExternalCephBlockProvisioner: "disable",
		EnableExternalCephFsProvisioner:    "disable",
	}
	_ = db.DB.Where("status = ?", constant.StatusRunning).Find(&provisoners).Error
	for _, p := range provisoners {
		switch p.Type {
		case "nfs":
			phase.EnableNfsProvisioner = "enable"
			if phase.NfsVersion == "v3" {
				continue
			}
			var vars map[string]string
			if err := json.Unmarshal([]byte(p.Vars), &vars); err != nil {
				continue
			}
			if _, ok := vars["storage_nfs_server_version"]; ok {
				phase.NfsVersion = vars["storage_nfs_server_version"]
			}
		case "gfs":
			phase.EnableGfsProvisioner = "enable"
			continue
		case "external-ceph-block":
			phase.EnableExternalCephBlockProvisioner = "enable"
			continue
		case "external-cephfs":
			phase.EnableExternalCephFsProvisioner = "enable"
			continue
		}
	}
	return phase.Run(aHelper.Kobe, aHelper.Writer)
}
