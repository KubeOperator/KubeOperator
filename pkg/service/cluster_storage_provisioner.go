package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	CreateStorageProvisioner(clusterName string, creation dto.ClusterStorageProvisionerCreation) (dto.ClusterStorageProvisioner, error)
	SyncStorageProvisioner(clusterName string, syncs []dto.ClusterStorageProvisionerSync) error
	DeleteStorageProvisioner(clusterName string, provisioner string) error
	BatchStorageProvisioner(clusterName string, batch dto.ClusterStorageProvisionerBatch) error
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
	clusterStorageProvisionerDTOS := []dto.ClusterStorageProvisioner{}
	ps, err := c.provisionerRepo.List(clusterName)
	if err != nil {
		return clusterStorageProvisionerDTOS, err
	}

	var syncList []dto.ClusterStorageProvisionerSync
	for _, p := range ps {
		syncList = append(syncList, dto.ClusterStorageProvisionerSync{
			Name:      p.Name,
			Namespace: p.Namespace,
			Type:      p.Type,
			Status:    p.Status,
		})
	}
	if err := c.SyncStorageProvisioner(clusterName, syncList); err != nil {
		return clusterStorageProvisionerDTOS, err
	}

	ps, err = c.provisionerRepo.List(clusterName)
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

func (c clusterStorageProvisionerService) DeleteStorageProvisioner(clusterName string, provisioner string) error {
	err := c.deleteProvisioner(clusterName, provisioner)
	if err != nil {
		return err
	}
	return c.provisionerRepo.Delete(clusterName, provisioner)
}

func (c clusterStorageProvisionerService) BatchStorageProvisioner(clusterName string, batch dto.ClusterStorageProvisionerBatch) error {
	switch batch.Operation {
	case constant.BatchOperationDelete:
		return c.provisionerRepo.BatchDelete(clusterName, batch.Items)
	default:
		return errors.New("not supported")
	}
}

func (c clusterStorageProvisionerService) CreateStorageProvisioner(clusterName string, creation dto.ClusterStorageProvisionerCreation) (dto.ClusterStorageProvisioner, error) {
	vars, _ := json.Marshal(creation.Vars)
	var dp dto.ClusterStorageProvisioner
	p := model.ClusterStorageProvisioner{
		Name:      creation.Name,
		Namespace: creation.Namespace,
		Type:      creation.Type,
		Vars:      string(vars),
		Status:    constant.StatusCreating,
	}

	if creation.IsInCluster {
		p.Status = constant.StatusRunning
		err := c.provisionerRepo.Save(clusterName, &p)
		return dp, err
	}

	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return dp, err
	}
	num := 0
	_ = db.DB.Model(&model.ClusterStorageProvisioner{}).Where("name = ? AND type = ? AND cluster_id = ?", p.Name, p.Type, cluster.ID).Count(&num).Error
	if num != 0 {
		return dp, errors.New("PROVISIONER_EXSIT")
	}

	err = c.provisionerRepo.Save(clusterName, &p)
	if err != nil {
		return dp, err
	}
	var registery model.SystemRegistry
	if cluster.Architectures == constant.ArchAMD64 {
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfAMD64).First(&registery).Error; err != nil {
			return dp, errors.New("load image pull port failed")
		}
	} else {
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfARM64).First(&registery).Error; err != nil {
			return dp, errors.New("load image pull port failed")
		}
	}
	playbook := c.loadPlayBookName(p.Type)
	task := model.TaskLogDetail{
		ID:            p.ID,
		Task:          playbook,
		ClusterID:     cluster.ID,
		LastProbeTime: time.Now().Unix(),
		Status:        constant.TaskLogStatusRunning,
	}
	if err := c.taskLogService.StartDetail(&task); err != nil {
		return dp, fmt.Errorf("save tasklog failed, err: %v", err)
	}

	//playbook
	go c.do(cluster, p, task, registery.RegistryPort)
	dp.ClusterStorageProvisioner = p
	_ = json.Unmarshal([]byte(p.Vars), &dp.Vars)
	return dp, nil
}

func (c clusterStorageProvisionerService) do(cluster model.Cluster, provisioner model.ClusterStorageProvisioner, task model.TaskLogDetail, repoPort int) {
	admCluster := adm.NewAnsibleHelper(cluster)
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, provisioner.ID)
	if err != nil {
		logger.Log.Error(err)
	}
	if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("id = ?", provisioner.ID).Update("status", constant.StatusInitializing).Error; err != nil {
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
		c.errCreateStorageProvisioner(cluster.Name, provisioner, err)
		return
	}

	// 获取创建参数
	if err := c.getVars(admCluster, cluster, provisioner); err != nil {
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
		c.errCreateStorageProvisioner(cluster.Name, provisioner, err)
		return
	}
	admCluster.Kobe.SetVar("registry_port", fmt.Sprint(repoPort))
	// 获取 k8s client
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
		c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create kubernetes Clientset error %s", err.Error()))
	}

	switch provisioner.Type {
	case "nfs":
		admCluster.Kobe.SetVar("storage_nfs_provisioner_name", provisioner.Name)
		if err := phases.RunPlaybookAndGetResult(admCluster.Kobe, NfsStorage, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			logger.Log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning(provisioner.Namespace, provisioner.Name, client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "rook-ceph":
		if err := phases.RunPlaybookAndGetResult(admCluster.Kobe, rookCephStorage, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			logger.Log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning(provisioner.Namespace, "rook-ceph-operator", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "vsphere":
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, vsphereStorage, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			logger.Log.Errorf("save provisioner status err: %s", err.Error())
			return
		}

		if err := phases.WaitForStatefulSetsRunning(provisioner.Namespace, "vsphere-csi-controller", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "external-ceph-block":
		admCluster.Kobe.SetVar("storage_rbd_provisioner_name", provisioner.Name)
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, externalCephRbdStorage, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			logger.Log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning(provisioner.Namespace, "external-ceph-block", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "external-cephfs":
		admCluster.Kobe.SetVar("storage_fs_provisioner_name", provisioner.Name)
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, externalCephFsStorage, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			logger.Log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning(provisioner.Namespace, "external-cephfs", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "oceanstor":
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, oceanStor, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			logger.Log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning(provisioner.Namespace, "huawei-csi-controller", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "cinder":
		admCluster.Kobe.SetVar("cinder_csi_version", "v1.20.0")
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, cinderStorage, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			logger.Log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForStatefulSetsRunning(provisioner.Namespace, "csi-cinder-controllerplugin", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "glusterfs":
		admCluster.Kobe.SetVar("type", provisioner.Type)
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, glusterfsStorage, "", writer); err != nil {
			_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusFailed, err.Error())
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		_ = c.taskLogService.EndDetail(&task, constant.TaskLogStatusSuccess, "")
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			logger.Log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
	}
	provisioner.Status = constant.StatusRunning
	_ = c.provisionerRepo.Save(cluster.Name, &provisioner)
}

func (c clusterStorageProvisionerService) errCreateStorageProvisioner(clusterName string, provisioner model.ClusterStorageProvisioner, err error) {
	logger.Log.Errorf(err.Error())
	provisioner.Status = constant.StatusFailed
	provisioner.Message = err.Error()
	_ = c.provisionerRepo.Save(clusterName, &provisioner)
}

func (c clusterStorageProvisionerService) loadPlayBookName(provisionerType string) string {
	switch provisionerType {
	case "nfs":
		return NfsStorage
	case "rook-ceph":
		return rookCephStorage
	case "vsphere":
		return vsphereStorage
	case "external-ceph-block":
		return externalCephRbdStorage
	case "external-cephfs":
		return rookCephStorage
	case "oceanstor":
		return oceanStor
	case "cinder":
		return cinderStorage
	case "glusterfs":
		return glusterfsStorage
	}
	return ""
}

func (c clusterStorageProvisionerService) SyncStorageProvisioner(clusterName string, provisioners []dto.ClusterStorageProvisionerSync) error {
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		logger.Log.Errorf("get kubernetes Clientset err, error: %s", err.Error())
	}
	for _, provisioner := range provisioners {
		if provisioner.Status == constant.StatusInitializing || provisioner.Status == constant.StatusTerminating {
			continue
		}
		if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("name = ?", provisioner.Name).Update("status", constant.StatusSynchronizing).Error; err != nil {
			logger.Log.Errorf("update host status to synchronizing error: %s", err.Error())
		}
		logger.Log.Infof("gather provisioner [%s] info", provisioner.Name)

		switch provisioner.Type {
		case "external-ceph-block":
			c.updateProvisionerStatus(client, "deployment", provisioner.Namespace, provisioner.Name)
		case "external-cephfs":
			c.updateProvisionerStatus(client, "deployment", provisioner.Namespace, provisioner.Name)
		case "nfs":
			c.updateProvisionerStatus(client, "deployment", provisioner.Namespace, provisioner.Name)
		case "vsphere":
			c.updateProvisionerStatus(client, "statefulSets", provisioner.Namespace, "vsphere-csi-controller")
		case "rook-ceph":
			c.updateProvisionerStatus(client, "deployment", provisioner.Namespace, "rook-ceph-operator")
		case "oceanstor":
			c.updateProvisionerStatus(client, "deployment", provisioner.Namespace, "huawei-csi-controller")
		case "cinder":
			c.updateProvisionerStatus(client, "statefulSets", provisioner.Namespace, "csi-cinder-controllerplugin")
		}
	}
	return nil
}

func (c clusterStorageProvisionerService) deleteProvisioner(clusterName string, provisionerName string) error {
	var provisioner model.ClusterStorageProvisioner
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	db.DB.Where("name = ? AND cluster_id = ?", provisionerName, cluster.ID).First(&provisioner)
	if provisioner.ID == "" {
		return errors.New("not found")
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return err
	}
	switch provisioner.Type {
	case "nfs":
		contextTo := context.TODO()
		_ = client.CoreV1().ServiceAccounts(provisioner.Namespace).Delete(contextTo, "nfs-client-provisioner", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "run-nfs-client-provisioner", metav1.DeleteOptions{})
		_ = client.AppsV1().Deployments(provisioner.Namespace).Delete(contextTo, provisioner.Name, metav1.DeleteOptions{})
	case "external-ceph":
		contextTo := context.TODO()
		_ = client.CoreV1().ServiceAccounts(provisioner.Namespace).Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoles().Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().Roles(provisioner.Namespace).Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().RoleBindings(provisioner.Namespace).Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		_ = client.PolicyV1beta1().PodSecurityPolicies().Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		_ = client.AppsV1().Deployments(provisioner.Namespace).Delete(contextTo, provisioner.Name, metav1.DeleteOptions{})
	case "oceanstor":
		contextTo := context.TODO()
		_ = client.CoreV1().ConfigMaps(provisioner.Namespace).Delete(contextTo, "huawei-csi-configmap", metav1.DeleteOptions{})
		_ = client.AppsV1().Deployments(provisioner.Namespace).Delete(contextTo, "huawei-csi-controller", metav1.DeleteOptions{})
		_ = client.AppsV1().DaemonSets(provisioner.Namespace).Delete(contextTo, "huawei-csi-node", metav1.DeleteOptions{})
		_ = client.CoreV1().ServiceAccounts(provisioner.Namespace).Delete(contextTo, "huawei-csi-controller", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoles().Delete(contextTo, "huawei-csi-provisioner-runner", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "huawei-csi-provisioner-role", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoles().Delete(contextTo, "huawei-csi-attacher-runner", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "huawei-csi-attacher-role", metav1.DeleteOptions{})
		_ = client.CoreV1().ServiceAccounts(provisioner.Namespace).Delete(contextTo, "huawei-csi-node", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoles().Delete(contextTo, "huawei-csi-driver-registrar-runner", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "huawei-csi-driver-registrar-role", metav1.DeleteOptions{})
	case "vsphere":
		contextTo := context.TODO()
		_ = client.CoreV1().ServiceAccounts(provisioner.Namespace).Delete(contextTo, "vsphere-csi-controller", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "vsphere-csi-controller-binding", metav1.DeleteOptions{})
		_ = client.RbacV1beta1().ClusterRoles().Delete(contextTo, "vsphere-csi-controller-role", metav1.DeleteOptions{})
		_ = client.AppsV1().StatefulSets(provisioner.Namespace).Delete(contextTo, "vsphere-csi-controller", metav1.DeleteOptions{})
		_ = client.StorageV1().CSIDrivers().Delete(contextTo, "csi.vsphere.vmware.com", metav1.DeleteOptions{})
		_ = client.AppsV1().DaemonSets(provisioner.Namespace).Delete(contextTo, "vsphere-csi-node", metav1.DeleteOptions{})
	case "cinder":
		contextTo := context.TODO()
		_ = client.CoreV1().ServiceAccounts(provisioner.Namespace).Delete(contextTo, "csi-cinder-controller-service", metav1.DeleteOptions{})
		_ = client.CoreV1().ServiceAccounts(provisioner.Namespace).Delete(contextTo, "csi-cinder-controller-sa", metav1.DeleteOptions{})
		_ = client.CoreV1().ServiceAccounts(provisioner.Namespace).Delete(contextTo, "csi-cinder-node-sa", metav1.DeleteOptions{})

		_ = client.AppsV1().StatefulSets(provisioner.Namespace).Delete(contextTo, "csi-cinder-controllerplugin", metav1.DeleteOptions{})
		_ = client.AppsV1().DaemonSets(provisioner.Namespace).Delete(contextTo, "csi-cinder-nodeplugin", metav1.DeleteOptions{})
	}
	return nil
}

func (c clusterStorageProvisionerService) updateProvisionerStatus(client *kubernetes.Clientset, source, namespace, name string) {
	status, errMsg := "", ""
	if source == "deployment" {
		ex, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			status = constant.StatusFailed
			errMsg = err.Error()
		} else {
			if ex.Status.ReadyReplicas == ex.Status.Replicas {
				status = constant.StatusRunning
			} else {
				status = constant.StatusWaiting
			}
		}
	}
	if source == "statefulSets" {
		ex, err := client.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			status = constant.StatusFailed
			errMsg = err.Error()
		} else {
			if ex.Status.ReadyReplicas == ex.Status.Replicas {
				status = constant.StatusRunning
			} else {
				status = constant.StatusWaiting
			}
		}
	}
	if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("name = ?", name).
		Updates(map[string]interface{}{"status": status, "Message": errMsg}).Error; err != nil {
		logger.Log.Errorf("update host status to failed error: %s", err.Error())
	}
}

func (c clusterStorageProvisionerService) getVars(admCluster *adm.AnsibleHelper, cluster model.Cluster, provisioner model.ClusterStorageProvisioner) error {
	if provisioner.Type == "glusterfs" {
		return nil
	}
	var (
		manifest       model.ClusterManifest
		storageVars    []model.VersionHelp
		storageDic     model.StorageProvisionerDic
		storageDicVars map[string]interface{}
		commonVars     map[string]interface{}
	)

	// 获取版本
	if err := db.DB.Where("name = ?", cluster.Version).First(&manifest).Error; err != nil {
		return fmt.Errorf("can't find manifest version: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(manifest.StorageVars), &storageVars); err != nil {
		return fmt.Errorf("unmarshal manifest.storageVars error %s", err.Error())
	}
	// 获取存储字典
	isExist := false
	for _, storage := range storageVars {
		if storage.Name == provisioner.Type {
			isExist = true
			if err := db.DB.Where("name = ? AND version = ?", storage.Name, storage.Version).First(&storageDic).Error; err != nil {
				return fmt.Errorf("can't find storage provisioner dic : %s", err.Error())
			}
			break
		}
	}
	if !isExist {
		return fmt.Errorf("can't find storage provisioner dic: %s", provisioner.Type)
	}
	if err := json.Unmarshal([]byte(storageDic.Vars), &storageDicVars); err != nil {
		return fmt.Errorf("unmarshal storageDic.Vars error %s", err.Error())
	}
	// 获取前端参数
	if err := json.Unmarshal([]byte(provisioner.Vars), &commonVars); err != nil {
		return fmt.Errorf("unmarshal provisioner.Vars error %s", err.Error())
	}

	for k, v := range storageDicVars {
		if v != nil {
			admCluster.Kobe.SetVar(k, fmt.Sprintf("%v", v))
		}
	}
	for k, v := range commonVars {
		if v != nil {
			admCluster.Kobe.SetVar(k, fmt.Sprintf("%v", v))
		}
	}
	return nil
}
