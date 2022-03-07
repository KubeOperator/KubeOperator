package adm

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
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
	FileName string
	Kobe     kobe.Interface
}

func NewCluster(cluster model.Cluster, fileName string) *Cluster {
	c := &Cluster{
		Cluster: cluster,
	}
	c.FileName = fileName
	c.Kobe = kobe.NewAnsible(&kobe.Config{
		Inventory: c.ParseInventory(),
	})
	for i := range facts.DefaultFacts {
		c.Kobe.SetVar(i, facts.DefaultFacts[i])
	}
	clusterVars := cluster.GetKobeVars()
	for k, v := range clusterVars {
		c.Kobe.SetVar(k, v)
	}
	c.Kobe.SetVar(facts.ClusterNameFactName, cluster.Name)
	var systemSetting model.SystemSetting
	db.DB.Model(&model.SystemSetting{}).Where(model.SystemSetting{Key: "ntp_server"}).First(&systemSetting)
	if systemSetting.ID != "" {
		c.Kobe.SetVar(facts.NtpServerName, systemSetting.Value)
	}
	maniFest, err := GetManiFestBy(cluster.Spec.Version)
	if err != nil {
		log.Errorf("get manifest by version %s failed, %v", cluster.Spec.Version, err)
	}
	if maniFest.Name != "" {
		vars := maniFest.GetVars()
		for k, v := range vars {
			c.Kobe.SetVar(k, v)
		}
	}
	return c
}

type ClusterAdm struct {
	createHandlers  []Handler
	upgradeHandlers []Handler
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
	return ca
}

func (ca *ClusterAdm) OnInitialize(c Cluster) (Cluster, error) {
	err := ca.Create(&c)
	return c, err
}

func (ca *ClusterAdm) OnUpgrade(c Cluster) (Cluster, error) {
	err := ca.Upgrade(&c)
	return c, err
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

func writeLog(msg string, fileName string) {
	writer, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
	if err != nil {
		log.Error(err)
	}
	defer writer.Close()
	if _, err := fmt.Fprintln(writer, msg); err != nil {
		log.Error(err.Error())
	}
}
