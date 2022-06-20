package adm

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	ConditionTypeDone = "EnsureDone"
)

type Handler func(*AnsibleHelper) error

func (h Handler) name() string {
	name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	i := strings.Index(name, "Ensure")
	if i == -1 {
		return ""
	}
	return strings.TrimSuffix(name[i:], "-fm")
}

func (c *AnsibleHelper) setCondition(newDetail model.TaskLogDetail) {
	var details []model.TaskLogDetail
	exist := false
	for _, detail := range c.LogDetail {
		if detail.Task == newDetail.Task {
			exist = true
			if newDetail.Status != detail.Status {
				detail.Status = newDetail.Status
			}
			if newDetail.Message != detail.Message {
				detail.Message = newDetail.Message
			}
			if !newDetail.LastProbeTime.IsZero() && newDetail.LastProbeTime != detail.LastProbeTime {
				detail.LastProbeTime = newDetail.LastProbeTime
			}
		}
		details = append(details, detail)
	}
	if !exist {
		if newDetail.LastProbeTime.IsZero() {
			newDetail.LastProbeTime = time.Now()
		}
		details = append(details, newDetail)
	}
	c.LogDetail = details
}

type AnsibleHelper struct {
	Status    string
	Message   string
	LogDetail []model.TaskLogDetail

	ClusterVersion        string
	ClusterUpgradeVersion string
	ClusterRuntime        string

	Writer io.Writer
	Kobe   kobe.Interface
}

func NewAnsibleHelper(cluster model.Cluster, writer ...io.Writer) *AnsibleHelper {
	c := &AnsibleHelper{
		Status:                constant.TaskLogStatusRunning,
		ClusterVersion:        cluster.Version,
		ClusterUpgradeVersion: cluster.UpgradeVersion,
		LogDetail:             cluster.TaskLog.Details,
	}
	if writer != nil {
		c.Writer = writer[0]
	}
	c.Kobe = kobe.NewAnsible(&kobe.Config{
		Inventory: cluster.ParseInventory(),
	})
	for i := range facts.DefaultFacts {
		c.Kobe.SetVar(i, facts.DefaultFacts[i])
	}
	clusterVars := cluster.GetKobeVars()
	for k, v := range clusterVars {
		c.Kobe.SetVar(k, v)
	}
	c.Kobe.SetVar(facts.ClusterNameFactName, cluster.Name)
	ntpServerRepo := repository.NewNtpServerRepository()
	ntps, _ := ntpServerRepo.GetAddressStr()
	c.Kobe.SetVar(facts.NtpServerFactName, ntps)
	maniFest, _ := GetManiFestBy(cluster.Version)
	if maniFest.Name != "" {
		vars := maniFest.GetVars()
		for k, v := range vars {
			c.Kobe.SetVar(k, v)
		}
	}
	return c
}

func NewAnsibleHelperWithNewWorker(cluster model.Cluster, workers []string, writer ...io.Writer) *AnsibleHelper {
	c := &AnsibleHelper{
		Status:                constant.TaskLogStatusRunning,
		ClusterVersion:        cluster.Version,
		ClusterUpgradeVersion: cluster.UpgradeVersion,
		LogDetail:             cluster.TaskLog.Details,
	}
	if writer != nil {
		c.Writer = writer[0]
	}
	inventory := cluster.ParseInventory()
	for i := range inventory.Groups {
		if inventory.Groups[i].Name == "new-worker" {
			inventory.Groups[i].Hosts = append(inventory.Groups[i].Hosts, workers...)
		}
	}
	c.Kobe = kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})
	for i := range facts.DefaultFacts {
		c.Kobe.SetVar(i, facts.DefaultFacts[i])
	}
	clusterVars := cluster.GetKobeVars()
	for k, v := range clusterVars {
		c.Kobe.SetVar(k, v)
	}
	c.Kobe.SetVar(facts.ClusterNameFactName, cluster.Name)
	ntpServerRepo := repository.NewNtpServerRepository()
	ntps, _ := ntpServerRepo.GetAddressStr()
	c.Kobe.SetVar(facts.NtpServerFactName, ntps)
	maniFest, _ := GetManiFestBy(cluster.Version)
	if maniFest.Name != "" {
		vars := maniFest.GetVars()
		for k, v := range vars {
			c.Kobe.SetVar(k, v)
		}
	}
	return c
}

type ClusterAdm struct {
	createHandlers    []Handler
	upgradeHandlers   []Handler
	addWorkerHandlers []Handler
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
		ca.EnsureInitMetricsServer,
		ca.EnsureInitIngressController,
		ca.EnsurePostInit,
	}
	ca.upgradeHandlers = []Handler{
		ca.EnsureUpgradeTaskStart,
		ca.EnsureBackupETCD,
		ca.EnsureUpgradeRuntime,
		ca.EnsureUpgradeETCD,
		ca.EnsureUpgradeKubernetes,
		ca.EnsureUpdateCertificates,
	}
	ca.addWorkerHandlers = []Handler{
		ca.EnsureAddWorkerTaskStart,
		ca.EnsureAddWorkerBaseSystemConfig,
		ca.EnsureAddWorkerContainerRuntime,
		ca.EnsureAddWorkerKubernetesComponent,
		ca.EnsureAddWorkerLoadBalancer,
		ca.EnsureAddWorkerCertificates,
		ca.EnsureAddWorkerWorker,
		ca.EnsureAddWorkerNetwork,
		ca.EnsureAddWorkerPost,
		ca.EnsureAddWorkerStorage,
	}
	return ca
}

func (ca *ClusterAdm) OnInitialize(ansible *AnsibleHelper) error {
	err := ca.Create(ansible)
	return err
}

func (ca *ClusterAdm) OnUpgrade(ansible *AnsibleHelper) error {
	err := ca.Upgrade(ansible)
	return err
}

func (ca *ClusterAdm) OnAddWorker(ansible *AnsibleHelper) error {
	err := ca.AddWorker(ansible)
	return err
}

func GetManiFestBy(name string) (dto.ClusterManifest, error) {
	var clusterManifest dto.ClusterManifest
	var mo model.ClusterManifest
	err := db.DB.Where(model.ClusterManifest{Name: name}).First(&mo).Error
	if err != nil {
		return clusterManifest, err
	}
	clusterManifest.Name = mo.Name
	clusterManifest.Version = mo.Version
	clusterManifest.IsActive = mo.IsActive
	var core []dto.NameVersion
	if err := json.Unmarshal([]byte(mo.CoreVars), &core); err != nil {
		return clusterManifest, err
	}
	clusterManifest.CoreVars = core
	var network []dto.NameVersion
	if err := json.Unmarshal([]byte(mo.NetworkVars), &network); err != nil {
		return clusterManifest, err
	}
	clusterManifest.NetworkVars = network
	var other []dto.NameVersion
	if err := json.Unmarshal([]byte(mo.OtherVars), &other); err != nil {
		return clusterManifest, err
	}
	clusterManifest.OtherVars = other
	return clusterManifest, err
}

func writeLog(msg string, writer io.Writer) {
	_, err := fmt.Fprintln(writer, msg)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}
