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

func (ca *ClusterAdm) AddWorker(c *Cluster, status *model.ClusterStatus) error {
	condition := ca.getAddWorkerCurrentCondition(status)
	if condition != nil {
		now := time.Now()
		f := ca.getAddWorkerHandler(condition.Name)
		err := f(c)
		if err != nil {
			ca.setAddCondition(status, model.ClusterStatusCondition{
				Name:          condition.Name,
				Status:        constant.ConditionFalse,
				LastProbeTime: now,
				Message:       err.Error(),
			})
			status.Phase = constant.ClusterFailed
			status.Message = err.Error()
			return nil
		}
		ca.setAddCondition(status, model.ClusterStatusCondition{
			Name:          condition.Name,
			Status:        constant.ConditionTrue,
			LastProbeTime: now,
		})

		nextConditionType := ca.getNextAddWorkerConditionName(condition.Name)
		if nextConditionType == ConditionTypeDone {
			status.Phase = constant.ClusterRunning
		} else {
			ca.setAddCondition(status, model.ClusterStatusCondition{
				Name:          nextConditionType,
				Status:        constant.ConditionUnknown,
				LastProbeTime: time.Now(),
				Message:       "",
			})
		}
	}
	return nil
}

func (ca *ClusterAdm) getAddWorkerCurrentCondition(status *model.ClusterStatus) *model.ClusterStatusCondition {
	if len(status.ClusterStatusConditions) == 0 {
		return &model.ClusterStatusCondition{
			Name:          ca.addWorkerHandlers[0].name(),
			Status:        constant.ConditionUnknown,
			LastProbeTime: time.Now(),
			Message:       "",
		}
	}
	for _, condition := range status.ClusterStatusConditions {
		if condition.Status == constant.ConditionFalse || condition.Status == constant.ConditionUnknown {
			return &condition
		}
	}
	return nil
}

func (ca *ClusterAdm) getAddWorkerHandler(conditionName string) Handler {
	for _, f := range ca.addWorkerHandlers {
		if conditionName == f.name() {
			return f
		}
	}
	return nil
}
func (ca *ClusterAdm) setAddCondition(status *model.ClusterStatus, newCondition model.ClusterStatusCondition) {
	var conditions []model.ClusterStatusCondition
	exist := false
	for _, condition := range status.ClusterStatusConditions {
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
	status.ClusterStatusConditions = conditions

}
func (ca *ClusterAdm) getNextAddWorkerConditionName(conditionName string) string {
	var (
		i int
		f Handler
	)
	for i, f = range ca.addWorkerHandlers {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if strings.Contains(name, conditionName) {
			break
		}
	}
	if i == len(ca.addWorkerHandlers)-1 {
		return ConditionTypeDone
	}
	next := ca.addWorkerHandlers[i+1]
	return next.name()
}

func (ca *ClusterAdm) EnsureAddWorkerTaskStart(c *Cluster) error {
	time.Sleep(5 * time.Second)
	writeLog("----add worker task start----", c.Writer)
	return nil
}

func (ca *ClusterAdm) EnsureAddWorkerBaseSystemConfig(c *Cluster) error {
	phase := prepare.AddWorkerBaseSystemConfigPhase{}
	err := phase.Run(c.Kobe, c.Writer)
	return err
}

func (ca *ClusterAdm) EnsureAddWorkerContainerRuntime(c *Cluster) error {
	phase := prepare.AddWorkerContainerRuntimePhase{}
	return phase.Run(c.Kobe, c.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerKubernetesComponent(c *Cluster) error {
	phase := prepare.AddWorkerKubernetesComponentPhase{}
	return phase.Run(c.Kobe, c.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerLoadBalancer(c *Cluster) error {
	phase := prepare.AddWorkerLoadBalancerPhase{}
	return phase.Run(c.Kobe, c.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerCertificates(c *Cluster) error {
	phase := prepare.AddWorkerCertificatesPhase{}
	return phase.Run(c.Kobe, c.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerWorker(c *Cluster) error {
	phase := initial.AddWorkerMasterPhase{}
	return phase.Run(c.Kobe, c.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerNetwork(c *Cluster) error {
	phase := initial.AddWorkerNetworkPhase{}
	return phase.Run(c.Kobe, c.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerPost(c *Cluster) error {
	phase := initial.AddWorkerPostPhase{}
	return phase.Run(c.Kobe, c.Writer)
}

func (ca *ClusterAdm) EnsureAddWorkerStorage(c *Cluster) error {
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
	return phase.Run(c.Kobe, c.Writer)
}
