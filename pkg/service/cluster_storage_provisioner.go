package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/jinzhu/gorm"
	"k8s.io/client-go/kubernetes"
)

const (
	oceanStor              = "10-plugin-cluster-storage-oceanstor.yml"
	externalCephRbdStorage = "10-plugin-cluster-storage-external-ceph-block.yml"
	externalCephFsStorage  = "10-plugin-cluster-storage-external-cephfs.yml"
	NfsStorage             = "10-plugin-cluster-storage-nfs.yml"
	rookCephStorage        = "10-plugin-cluster-storage-rook-ceph.yml"
	vsphereStorage         = "10-plugin-cluster-storage-vsphere.yml"
	cinderStorage          = "10-plugin-cluster-storage-cinder.yml"
	glusterfsStorage       = "10-plugin-cluster-storage-glusterfs.yml"
)

type ClusterStorageProvisionerService interface {
	ListStorageProvisioner(clusterName string) ([]dto.ClusterStorageProvisioner, error)
	CreateStorageProvisioner(clusterName string, creation dto.ClusterStorageProvisionerCreation) error
	SyncStorageProvisioner(clusterName string, syncs []dto.ClusterStorageProvisionerSync) error
	DeleteStorageProvisioner(clusterName string, provisioner string) error
}

type clusterStorageProvisionerService struct {
	provisionerRepo repository.ClusterStorageProvisionerRepository
	clusterService  ClusterService
	clusterRepo     repository.ClusterRepository
	taskLogService  TaskLogService
}

func NewClusterStorageProvisionerService() ClusterStorageProvisionerService {
	return &clusterStorageProvisionerService{
		provisionerRepo: repository.NewClusterStorageProvisionerRepository(),
		clusterService:  NewClusterService(),
		clusterRepo:     repository.NewClusterRepository(),
		taskLogService:  NewTaskLogService(),
	}
}

func (c clusterStorageProvisionerService) ListStorageProvisioner(clusterName string) ([]dto.ClusterStorageProvisioner, error) {
	var clusterStorageProvisionerDTOS []dto.ClusterStorageProvisioner
	ps, err := c.provisionerRepo.List(clusterName)
	if err != nil {
		return clusterStorageProvisionerDTOS, err
	}
	for _, p := range ps {
		var vars map[string]interface{}
		_ = json.Unmarshal([]byte(p.Vars), &vars)
		clusterStorageProvisionerDTOS = append(clusterStorageProvisionerDTOS, dto.ClusterStorageProvisioner{
			ClusterStorageProvisioner: p,
			Vars:                      vars,
		})
	}

	return clusterStorageProvisionerDTOS, nil
}

func (c clusterStorageProvisionerService) SyncStorageProvisioner(clusterName string, syncs []dto.ClusterStorageProvisionerSync) error {
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return err
	}
	go c.dosync(client, syncs)

	return nil
}

func (c clusterStorageProvisionerService) CreateStorageProvisioner(clusterName string, creation dto.ClusterStorageProvisionerCreation) error {
	var dp model.ClusterStorageProvisioner
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	if err := db.DB.Model(&model.ClusterStorageProvisioner{}).
		Where("name = ? AND type = ? AND cluster_id = ?", creation.Name, creation.Type, cluster.ID).
		First(&dp).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if dp.ID != "" {
		return errors.New("PROVISIONER_EXSIT")
	}
	dp = creation.ProvisionerCreate2Mo()
	dp.ClusterID = cluster.ID
	if err := db.DB.Create(&dp).Error; err != nil {
		return err
	}
	if creation.IsInCluster {
		return nil
	}

	var registery model.SystemRegistry
	if cluster.Architectures == constant.ArchAMD64 {
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfAMD64).First(&registery).Error; err != nil {
			return errors.New("load image pull port failed")
		}
	} else {
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfARM64).First(&registery).Error; err != nil {
			return errors.New("load image pull port failed")
		}
	}
	creation.Vars["registry_port"] = fmt.Sprint(registery.RegistryPort)
	playbook := c.loadPlayBookName(dp, creation.Vars)
	task := model.TaskLogDetail{
		ID:            dp.ID,
		Task:          fmt.Sprintf("%s (%s)", playbook, constant.StatusEnabled),
		ClusterID:     cluster.ID,
		LastProbeTime: time.Now().Unix(),
		Status:        constant.TaskLogStatusRunning,
	}
	if err := c.taskLogService.StartDetail(&task); err != nil {
		return fmt.Errorf("save tasklog failed, err: %v", err)
	}

	if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("id = ?", dp.ID).
		Updates(map[string]interface{}{"status": constant.StatusInitializing, "message": ""}).Error; err != nil {
		return err
	}

	//playbook
	go c.docreate(&cluster, task, dp, creation.Vars)
	return nil
}

func (c clusterStorageProvisionerService) docreate(cluster *model.Cluster, task model.TaskLogDetail, dp model.ClusterStorageProvisioner, vars map[string]interface{}) {
	admCluster, writer, err := c.loadAdmCluster(*cluster, dp, vars, constant.StatusEnabled)
	if err != nil {
		_ = c.taskLogService.EndDetail(&task, dp.Name, "provisioner", constant.TaskLogStatusFailed, err.Error())
		c.errHandlerProvisioner(dp, constant.StatusFailed, err)
		return
	}

	client, err := clusterUtil.NewClusterClient(cluster)
	if err != nil {
		_ = c.taskLogService.EndDetail(&task, dp.Name, "provisioner", constant.TaskLogStatusFailed, err.Error())
		c.errHandlerProvisioner(dp, constant.StatusFailed, err)
		return
	}

	playbook := strings.ReplaceAll(task.Task, " (enable)", "")
	if err := phases.RunPlaybookAndGetResult(admCluster.Kobe, playbook, "", writer); err != nil {
		_ = c.taskLogService.EndDetail(&task, dp.Name, "provisioner", constant.TaskLogStatusFailed, err.Error())
		c.errHandlerProvisioner(dp, constant.StatusFailed, err)
		return
	}
	_ = c.taskLogService.EndDetail(&task, dp.Name, "provisioner", constant.TaskLogStatusSuccess, "")
	dp.Status = constant.StatusWaiting
	if err := db.DB.Save(&dp).Error; err != nil {
		logger.Log.Errorf("save storage provisioner status err: %s", err.Error())
		return
	}
	syncItem := dto.ClusterStorageProvisionerSync{
		Name:      dp.Name,
		Type:      dp.Type,
		Namespace: dp.Namespace,
		Status:    dp.Status,
	}
	c.dosync(client, []dto.ClusterStorageProvisionerSync{syncItem})
}

func (c clusterStorageProvisionerService) DeleteStorageProvisioner(clusterName string, provisionerName string) error {
	var provisioner model.ClusterStorageProvisioner
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	db.DB.Where("name = ? AND cluster_id = ?", provisionerName, cluster.ID).First(&provisioner)
	if provisioner.ID == "" {
		return errors.New("not found")
	}

	Vars := make(map[string]interface{})
	if err := json.Unmarshal([]byte(provisioner.Vars), &Vars); err != nil {
		return err
	}
	playbook := c.loadPlayBookName(provisioner, Vars)
	task := model.TaskLogDetail{
		ID:            fmt.Sprintf("%s (%s)", provisioner.ID, constant.StatusDisabled),
		Task:          fmt.Sprintf("%s (%s)", playbook, constant.StatusDisabled),
		ClusterID:     cluster.ID,
		LastProbeTime: time.Now().Unix(),
		Status:        constant.TaskLogStatusRunning,
	}
	if err := c.taskLogService.StartDetail(&task); err != nil {
		return fmt.Errorf("save tasklog failed, err: %v", err)
	}

	if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("id = ?", provisioner.ID).
		Updates(map[string]interface{}{"status": constant.StatusTerminating, "message": ""}).Error; err != nil {
		return err
	}

	go c.dodelete(&cluster, task, provisioner, Vars)

	return nil
}

func (c clusterStorageProvisionerService) dodelete(cluster *model.Cluster, task model.TaskLogDetail, provisioner model.ClusterStorageProvisioner, vars map[string]interface{}) {
	admCluster, writer, err := c.loadAdmCluster(*cluster, provisioner, vars, constant.StatusDisabled)
	if err != nil {
		_ = c.taskLogService.EndDetail(&task, provisioner.Name, "provisioner", constant.TaskLogStatusFailed, err.Error())
		c.errHandlerProvisioner(provisioner, constant.StatusFailed, err)
		return
	}

	playbook := strings.ReplaceAll(task.Task, " (disable)", "")
	if err := phases.RunPlaybookAndGetResult(admCluster.Kobe, playbook, "", writer); err != nil {
		_ = c.taskLogService.EndDetail(&task, provisioner.Name, "provisioner", constant.TaskLogStatusFailed, err.Error())
		c.errHandlerProvisioner(provisioner, constant.StatusFailed, err)
		return
	}
	_ = c.taskLogService.EndDetail(&task, provisioner.Name, "provisioner", constant.TaskLogStatusSuccess, "")
	_ = db.DB.Where("id = ?", provisioner.ID).Delete(&model.ClusterStorageProvisioner{})
}

func (c clusterStorageProvisionerService) errHandlerProvisioner(provisioner model.ClusterStorageProvisioner, status string, err error) {
	logger.Log.Errorf(err.Error())
	provisioner.Status = status
	provisioner.Message = err.Error()
	_ = db.DB.Save(&provisioner)
}

func (c clusterStorageProvisionerService) loadPlayBookName(provisioner model.ClusterStorageProvisioner, vars map[string]interface{}) string {
	switch provisioner.Type {
	case "nfs":
		vars["storage_nfs_provisioner_name"] = provisioner.Name
		return NfsStorage
	case "rook-ceph":
		return rookCephStorage
	case "vsphere":
		return vsphereStorage
	case "external-ceph-block":
		vars["storage_rbd_provisioner_name"] = provisioner.Name
		return externalCephRbdStorage
	case "external-cephfs":
		vars["storage_fs_provisioner_name"] = provisioner.Name
		return externalCephFsStorage
	case "oceanstor":
		return oceanStor
	case "cinder":
		vars["cinder_csi_version"] = "v1.20.0"
		return cinderStorage
	case "glusterfs":
		vars["type"] = provisioner.Type
		return glusterfsStorage
	}
	return ""
}

func (c clusterStorageProvisionerService) dosync(client *kubernetes.Clientset, provisioners []dto.ClusterStorageProvisionerSync) {
	for _, provisioner := range provisioners {
		if provisioner.Status == constant.StatusInitializing || provisioner.Status == constant.StatusTerminating {
			continue
		}
		if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("name = ?", provisioner.Name).Update("status", constant.StatusSynchronizing).Error; err != nil {
			logger.Log.Errorf("update host status to synchronizing error: %s", err.Error())
		}

		switch provisioner.Type {
		case "external-cephfs", "external-ceph-block", "nfs":
			if err := phases.WaitForDeployRunning(provisioner.Namespace, provisioner.Name, client); err != nil {
				c.changeStatus(provisioner, constant.StatusFailed, err)
				continue
			}
			c.changeStatus(provisioner, constant.StatusRunning, nil)
		case "vsphere":
			if err := phases.WaitForStatefulSetsRunning(provisioner.Namespace, "vsphere-csi-controller", client); err != nil {
				c.changeStatus(provisioner, constant.StatusFailed, err)
				continue
			}
			c.changeStatus(provisioner, constant.StatusRunning, nil)
		case "rook-ceph":
			if err := phases.WaitForDeployRunning(provisioner.Namespace, "rook-ceph-operator", client); err != nil {
				c.changeStatus(provisioner, constant.StatusFailed, err)
				continue
			}
			c.changeStatus(provisioner, constant.StatusRunning, nil)
		case "oceanstor":
			if err := phases.WaitForDeployRunning(provisioner.Namespace, "huawei-csi-controller", client); err != nil {
				c.changeStatus(provisioner, constant.StatusFailed, err)
				continue
			}
			c.changeStatus(provisioner, constant.StatusRunning, nil)
		case "cinder":
			if err := phases.WaitForStatefulSetsRunning(provisioner.Namespace, "csi-cinder-controllerplugin", client); err != nil {
				c.changeStatus(provisioner, constant.StatusFailed, err)
				continue
			}
			c.changeStatus(provisioner, constant.StatusRunning, nil)
		}
	}
}

func (c clusterStorageProvisionerService) changeStatus(provisioner dto.ClusterStorageProvisionerSync, status string, err error) {
	if status == constant.StatusRunning {
		if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("name = ?", provisioner.Name).
			Updates(map[string]interface{}{"status": status, "Message": ""}).Error; err != nil {
			logger.Log.Errorf("update provisioner status to failed error: %s", err.Error())
		}
		return
	}

	if status == constant.StatusFailed && provisioner.Status == constant.StatusWaiting {
		if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("name = ?", provisioner.Name).
			Updates(map[string]interface{}{"status": constant.StatusNotReady, "Message": err.Error()}).Error; err != nil {
			logger.Log.Errorf("update provisioner status to failed error: %s", err.Error())
		}
		return
	}

	if status == constant.StatusFailed {
		if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("name = ?", provisioner.Name).
			Updates(map[string]interface{}{"status": constant.StatusFailed, "Message": err.Error()}).Error; err != nil {
			logger.Log.Errorf("update provisioner status to failed error: %s", err.Error())
		}
	}
}

func (c clusterStorageProvisionerService) loadAdmCluster(cluster model.Cluster, provisioner model.ClusterStorageProvisioner, vars map[string]interface{}, operation string) (*adm.AnsibleHelper, io.Writer, error) {
	admCluster := adm.NewAnsibleHelper(cluster)

	if len(vars) != 0 {
		for k, v := range vars {
			if v != nil {
				admCluster.Kobe.SetVar(k, fmt.Sprintf("%v", v))
			}
		}
	}
	var (
		manifest       model.ClusterManifest
		storageVars    []model.VersionHelp
		storageDic     model.StorageProvisionerDic
		storageDicVars map[string]interface{}
	)

	// 获取版本
	if err := db.DB.Where("name = ?", cluster.Version).First(&manifest).Error; err != nil {
		logger.Log.Errorf("can't find manifest version: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(manifest.StorageVars), &storageVars); err != nil {
		logger.Log.Errorf("unmarshal manifest.storageVars error %s", err.Error())
	}
	// 获取存储字典
	isExist := false
	for _, storage := range storageVars {
		if storage.Name == provisioner.Type {
			isExist = true
			if err := db.DB.Where("name = ? AND version = ?", storage.Name, storage.Version).First(&storageDic).Error; err != nil {
				logger.Log.Errorf("can't find storage provisioner dic : %s", err.Error())
			}
			break
		}
	}
	if !isExist {
		logger.Log.Errorf("can't find storage provisioner dic: %s", provisioner.Type)
	}
	if err := json.Unmarshal([]byte(storageDic.Vars), &storageDicVars); err != nil {
		logger.Log.Errorf("unmarshal storageDic.Vars error %s", err.Error())
	}

	for k, v := range storageDicVars {
		if v != nil {
			admCluster.Kobe.SetVar(k, fmt.Sprintf("%v", v))
		}
	}

	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, fmt.Sprintf("%s (%s)", provisioner.ID, operation))
	if err != nil {
		return admCluster, writer, err
	}

	switch provisioner.Type {
	case "nfs":
		admCluster.Kobe.SetVar(facts.EnableNfsFactName, operation)
	case "gfs":
		admCluster.Kobe.SetVar(facts.EnableGfsFactName, operation)
	case "external-ceph-block":
		fmt.Println(facts.EnableCephBlockFactName, operation)
		admCluster.Kobe.SetVar(facts.EnableCephBlockFactName, operation)
	case "external-cephfs":
		admCluster.Kobe.SetVar(facts.EnableCephFsFactName, operation)
	case "cinder":
		admCluster.Kobe.SetVar(facts.EnableCinderFactName, operation)
	case "vsphere":
		admCluster.Kobe.SetVar(facts.EnableVsphereFactName, operation)
	case "oceanstor":
		admCluster.Kobe.SetVar(facts.EnableOceanstorFactName, operation)
	case "rook-ceph":
		admCluster.Kobe.SetVar(facts.EnableRookFactName, operation)
	}
	return admCluster, writer, err
}
