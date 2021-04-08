package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	oceanStor           = "10-plugin-cluster-storage-oceanstor.yml"
	externalCephStorage = "10-plugin-cluster-storage-external-ceph.yml"
	NfsStorage          = "10-plugin-cluster-storage-nfs.yml"
	rookCephStorage     = "10-plugin-cluster-storage-rook-ceph.yml"
	vsphereStorage      = "10-plugin-cluster-storage-vsphere.yml"
	cinderStorage       = "10-plugin-cluster-storage-cinder.yml"
	glusterfsStorage    = "10-plugin-cluster-storage-glusterfs.yml"
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
}

func NewClusterStorageProvisionerService() ClusterStorageProvisionerService {
	return &clusterStorageProvisionerService{
		provisionerRepo: repository.NewClusterStorageProvisionerRepository(),
		clusterService:  NewClusterService(),
	}
}

func (c clusterStorageProvisionerService) ListStorageProvisioner(clusterName string) ([]dto.ClusterStorageProvisioner, error) {
	clusterStorageProvisionerDTOS := []dto.ClusterStorageProvisioner{}
	// 获取 k8s client
	client, err := c.getBaseParam(clusterName)
	if err != nil {
		return clusterStorageProvisionerDTOS, err
	}
	ps, err := c.provisionerRepo.List(clusterName)
	if err != nil {
		return clusterStorageProvisionerDTOS, err
	}
	for _, p := range ps {
		if p.Status != constant.StatusWaiting && p.Status != constant.StatusInitializing {
			var syncModel dto.ClusterStorageProvisionerSync
			syncModel.Name = p.Name
			syncModel.Status = p.Status
			syncModel.Type = p.Type
			if err := c.sync(client, syncModel); err != nil {
				p.Status = constant.ClusterNotReady
				p.Message = err.Error()
				db.DB.Save(&p)
			}
		}

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
		Name:   creation.Name,
		Type:   creation.Type,
		Vars:   string(vars),
		Status: constant.ClusterInitializing,
	}

	cluster, err := c.clusterService.Get(clusterName)
	if err != nil {
		return dp, err
	}
	err = c.provisionerRepo.Save(clusterName, &p)
	if err != nil {
		return dp, err
	}
	//playbook
	go c.do(cluster.Cluster, p)
	dp.ClusterStorageProvisioner = p
	_ = json.Unmarshal([]byte(p.Vars), &dp.Vars)
	return dp, nil
}

func (c clusterStorageProvisionerService) do(cluster model.Cluster, provisioner model.ClusterStorageProvisioner) {
	admCluster := adm.NewCluster(cluster)
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, provisioner.ID)
	if err != nil {
		log.Error(err)
	}

	// 获取创建参数
	if err := c.getVars(admCluster, cluster, provisioner); err != nil {
		c.errCreateStorageProvisioner(cluster.Name, provisioner, err)
		return
	}
	// 获取 k8s client
	client, err := c.getBaseParam(cluster.Name)
	if err != nil {
		c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create kubernetes Clientset error %s", err.Error()))
	}

	switch provisioner.Type {
	case "nfs":
		admCluster.Kobe.SetVar("storage_nfs_provisioner_name", provisioner.Name)
		if err := phases.RunPlaybookAndGetResult(admCluster.Kobe, NfsStorage, "", writer); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning("kube-system", provisioner.Name, client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "rook-ceph":
		if err := phases.RunPlaybookAndGetResult(admCluster.Kobe, rookCephStorage, "", writer); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning("rook-ceph", "rook-ceph-operator", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "vsphere":
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, vsphereStorage, "", writer); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForStatefulSetsRunning("kube-system", "vsphere-csi-controller", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "external-ceph":
		admCluster.Kobe.SetVar("storage_rbd_provisioner_name", provisioner.Name)
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, externalCephStorage, "", writer); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning("kube-system", "external-ceph", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "oceanstor":
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, oceanStor, "", writer); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForDeployRunning("kube-system", "huawei-csi-controller", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "cinder":
		admCluster.Kobe.SetVar("cinder_csi_version", "v1.20.0")
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, cinderStorage, "", writer); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
		if err := phases.WaitForStatefulSetsRunning("kube-system", "csi-cinder-controllerplugin", client); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("waitting provisioner running error %s", err.Error()))
			return
		}
	case "glusterfs":
		admCluster.Kobe.SetVar("type", provisioner.Type)
		if err = phases.RunPlaybookAndGetResult(admCluster.Kobe, glusterfsStorage, "", writer); err != nil {
			c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("create provisioner error %s", err.Error()))
			return
		}
		provisioner.Status = constant.StatusWaiting
		if err := c.provisionerRepo.Save(cluster.Name, &provisioner); err != nil {
			log.Errorf("save provisioner status err: %s", err.Error())
			return
		}
	}
	provisioner.Status = constant.ClusterRunning
	_ = c.provisionerRepo.Save(cluster.Name, &provisioner)
}

func (c clusterStorageProvisionerService) errCreateStorageProvisioner(clusterName string, provisioner model.ClusterStorageProvisioner, err error) {
	log.Errorf(err.Error())
	provisioner.Status = constant.ClusterFailed
	provisioner.Message = err.Error()
	_ = c.provisionerRepo.Save(clusterName, &provisioner)
}

func (c clusterStorageProvisionerService) SyncStorageProvisioner(clusterName string, provisioners []dto.ClusterStorageProvisionerSync) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 2)
	for _, provisioner := range provisioners {
		if provisioner.Status == constant.ClusterInitializing || provisioner.Status == constant.ClusterTerminating {
			continue
		}
		if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("name = ?", provisioner.Name).Update("status", constant.ClusterSynchronizing).Error; err != nil {
			log.Errorf("update host status to synchronizing error: %s", err.Error())
		}

		wg.Add(1)
		go func(provisioner dto.ClusterStorageProvisionerSync) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			log.Infof("gather provisioner [%s] info", provisioner.Name)
			client, err := c.getBaseParam(clusterName)
			if err != nil {
				log.Errorf("get kubernetes Clientset err, error: %s", err.Error())
			}
			if err := c.sync(client, provisioner); err != nil {
				log.Errorf("gather provisioner info error: %s", err.Error())
				if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("name = ?", provisioner.Name).
					Updates(map[string]interface{}{
						"status":  constant.ClusterFailed,
						"message": err.Error(),
					}).Error; err != nil {
					log.Errorf("update host status to failed error: %s", err.Error())
				}
			} else {
				if err := db.DB.Model(&model.ClusterStorageProvisioner{}).Where("name = ?", provisioner.Name).
					Updates(map[string]interface{}{"status": constant.ClusterRunning}).Error; err != nil {
					log.Errorf("update host status to running error: %s", err.Error())
				}
			}
		}(provisioner)
	}
	return nil
}

func (c clusterStorageProvisionerService) sync(client *kubernetes.Clientset, provisioner dto.ClusterStorageProvisionerSync) error {
	switch provisioner.Type {
	case "external-ceph":
		ex, err := client.AppsV1().Deployments("kube-system").Get(context.TODO(), "external-ceph", metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
		}
		if ex.Status.ReadyReplicas < 1 {
			return fmt.Errorf("not ready")
		}
	case "nfs":
		nfs, err := client.AppsV1().Deployments("kube-system").Get(context.TODO(), provisioner.Name, metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
		}
		if nfs.Status.ReadyReplicas < 1 {
			return fmt.Errorf("not ready")
		}
	case "vsphere":
		vs, err := client.AppsV1().StatefulSets("kube-system").Get(context.TODO(), "vsphere-csi-controller", metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
		}
		if vs.Status.ReadyReplicas < 1 {
			return fmt.Errorf("not ready")
		}
	case "rook-ceph":
		rook, err := client.AppsV1().Deployments("rook-ceph").Get(context.TODO(), "rook-ceph-operator", metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
		}
		if rook.Status.ReadyReplicas < 1 {
			return fmt.Errorf("not ready")
		}
	case "oceanstor":
		oc, err := client.AppsV1().Deployments("kube-system").Get(context.TODO(), "huawei-csi-controller", metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
		}
		if oc.Status.ReadyReplicas < 1 {
			return fmt.Errorf("not ready")
		}
	case "cinder":
		oc, err := client.AppsV1().StatefulSets("kube-system").Get(context.TODO(), "csi-cinder-controllerplugin", metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
		}
		if oc.Status.ReadyReplicas < 1 {
			return fmt.Errorf("not ready")
		}
	}
	return nil
}

func (c clusterStorageProvisionerService) deleteProvisioner(clusterName string, provisionerName string) error {
	var provisioner model.ClusterStorageProvisioner
	db.DB.Where("name = ?", provisionerName).First(&provisioner)
	if provisioner.ID == "" {
		return errors.New("not found")
	}
	client, err := c.getBaseParam(clusterName)
	if err != nil {
		return err
	}
	switch provisioner.Type {
	case "nfs":
		contextTo := context.TODO()
		err := client.CoreV1().ServiceAccounts("kube-system").Delete(contextTo, "nfs-client-provisioner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "run-nfs-client-provisioner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.AppsV1().Deployments("kube-system").Delete(contextTo, provisioner.Name, metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
	case "external-ceph":
		contextTo := context.TODO()
		err := client.CoreV1().ServiceAccounts("kube-system").Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoles().Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().Roles("kube-system").Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().RoleBindings("kube-system").Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.PolicyV1beta1().PodSecurityPolicies().Delete(contextTo, "rbd-provisioner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.AppsV1().Deployments("kube-system").Delete(contextTo, provisioner.Name, metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
	case "oceanstor":
		contextTo := context.TODO()
		err = client.CoreV1().ConfigMaps("kube-system").Delete(contextTo, "huawei-csi-configmap", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.AppsV1().Deployments("kube-system").Delete(contextTo, "huawei-csi-controller", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.AppsV1().DaemonSets("kube-system").Delete(contextTo, "huawei-csi-node", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err := client.CoreV1().ServiceAccounts("kube-system").Delete(contextTo, "huawei-csi-controller", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoles().Delete(contextTo, "huawei-csi-provisioner-runner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "huawei-csi-provisioner-role", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoles().Delete(contextTo, "huawei-csi-attacher-runner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "huawei-csi-attacher-role", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.CoreV1().ServiceAccounts("kube-system").Delete(contextTo, "huawei-csi-node", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoles().Delete(contextTo, "huawei-csi-driver-registrar-runner", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "huawei-csi-driver-registrar-role", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
	case "vsphere":
		contextTo := context.TODO()
		err := client.CoreV1().ServiceAccounts("kube-system").Delete(contextTo, "vsphere-csi-controller", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoleBindings().Delete(contextTo, "vsphere-csi-controller-binding", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.RbacV1beta1().ClusterRoles().Delete(contextTo, "vsphere-csi-controller-role", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.AppsV1().StatefulSets("kube-system").Delete(contextTo, "vsphere-csi-controller", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.StorageV1().CSIDrivers().Delete(contextTo, "csi.vsphere.vmware.com", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.AppsV1().DaemonSets("kube-system").Delete(contextTo, "vsphere-csi-node", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
	case "cinder":
		contextTo := context.TODO()
		err := client.CoreV1().ServiceAccounts("kube-system").Delete(contextTo, "csi-cinder-controller-service", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.CoreV1().ServiceAccounts("kube-system").Delete(contextTo, "csi-cinder-controller-sa", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.CoreV1().ServiceAccounts("kube-system").Delete(contextTo, "csi-cinder-node-sa", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}

		err = client.AppsV1().StatefulSets("kube-system").Delete(contextTo, "csi-cinder-controllerplugin", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.StorageV1().CSIDrivers().Delete(contextTo, "cinder.csi.openstack.org", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
		err = client.AppsV1().DaemonSets("kube-system").Delete(contextTo, "csi-cinder-nodeplugin", metav1.DeleteOptions{})
		if err != nil && checkError(err) {
			return err
		}
	}
	return nil
}

func checkError(err error) bool {
	if e, ok := err.(*errors2.StatusError); ok {
		if e.ErrStatus.Code == 404 {
			return false
		} else {
			return true
		}
	}
	return true
}

func (c clusterStorageProvisionerService) getBaseParam(clusterName string) (*kubernetes.Clientset, error) {
	var client *kubernetes.Clientset
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return client, err
	}
	endpoints, err := c.clusterService.GetApiServerEndpoints(clusterName)
	client, err = kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Token: secret.KubernetesToken,
		Hosts: endpoints,
	})
	if err != nil {
		return client, err
	}
	return client, nil
}

func (c clusterStorageProvisionerService) getVars(admCluster *adm.Cluster, cluster model.Cluster, provisioner model.ClusterStorageProvisioner) error {
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
	if err := db.DB.Where("name = ?", cluster.Spec.Version).First(&manifest).Error; err != nil {
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
