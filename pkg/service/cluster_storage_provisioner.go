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
	"github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	oceanStor           = "10-plugin-cluster-storage-oceanstor.yml"
	externalCephStorage = "10-plugin-cluster-storage-external-ceph.yml"
	NfsStorage          = "10-plugin-cluster-storage-nfs.yml"
	rookCephStorage     = "10-plugin-cluster-storage-rook-ceph.yml"
	vsphereStorage      = "10-plugin-cluster-storage-vsphere.yml"
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
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return clusterStorageProvisionerDTOS, err
	}
	endpoints, err := c.clusterService.GetApiServerEndpoints(clusterName)
	client, err := kubernetes.NewKubernetesClient(&kubernetes.Config{
		Token: secret.KubernetesToken,
		Hosts: endpoints,
	})
	if err != nil {
		return clusterStorageProvisionerDTOS, err
	}
	deploymentsList, err := client.AppsV1().Deployments("kube-system").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return clusterStorageProvisionerDTOS, err
	}
	ps, err := c.provisionerRepo.List(clusterName)
	if err != nil {
		return clusterStorageProvisionerDTOS, err
	}
	for _, p := range ps {
		for _, item := range deploymentsList.Items {
			if p.Name == item.Name {
				if item.Status.ReadyReplicas < item.Status.Replicas {
					p.Status = "NotReady"
					var message string
					for _, condition := range item.Status.Conditions {
						message = condition.Message + message
					}
					p.Message = message
					db.DB.Save(&p)
				}
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
	var (
		manifest       model.ClusterManifest
		storageVars    []model.VersionHelp
		storageDic     model.StorageProvisionerDic
		storageDicVars map[string]interface{}
		commonVars     map[string]interface{}
	)

	admCluster := adm.NewCluster(cluster)
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, provisioner.ID)
	if err != nil {
		log.Error(err)
	}

	// 获取版本
	if err := db.DB.Where("name = ?", cluster.Spec.Version).First(&manifest).Error; err != nil {
		c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("can find manifest version: %s", err.Error()))
		return
	}
	if err := json.Unmarshal([]byte(manifest.StorageVars), &storageVars); err != nil {
		c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("unmarshal manifest.storageVars error %s", err.Error()))
		return
	}
	// 获取存储字典
	isExist := false
	for _, storage := range storageVars {
		if storage.Name == provisioner.Type {
			isExist = true
			if err := db.DB.Where("name = ? AND version = ?", storage.Name, storage.Version).First(&storageDic).Error; err != nil {
				c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("can find storage provisioner dic : %s", err.Error()))
				return
			}
			break
		}
	}
	if !isExist {
		c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("can find storage provisioner dic: %s", provisioner.Type))
		return
	}
	if err := json.Unmarshal([]byte(storageDic.Vars), &storageDicVars); err != nil {
		c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("unmarshal storageDic.Vars error %s", err.Error()))
		return
	}
	// 获取前端参数
	if err := json.Unmarshal([]byte(provisioner.Vars), &commonVars); err != nil {
		c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("unmarshal provisioner.Vars error %s", err.Error()))
		return
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
	switch provisioner.Type {
	case "nfs":
		admCluster.Kobe.SetVar("storage_nfs_provisioner_name", provisioner.Name)
		err = phases.RunPlaybookAndGetResult(admCluster.Kobe, NfsStorage, "", writer)
	case "rook-ceph":
		err = phases.RunPlaybookAndGetResult(admCluster.Kobe, rookCephStorage, "", writer)
	case "vsphere":
		err = phases.RunPlaybookAndGetResult(admCluster.Kobe, vsphereStorage, "", writer)
	case "external-ceph":
		admCluster.Kobe.SetVar("storage_rbd_provisioner_name", provisioner.Name)
		err = phases.RunPlaybookAndGetResult(admCluster.Kobe, externalCephStorage, "", writer)
	case "oceanstor":
		err = phases.RunPlaybookAndGetResult(admCluster.Kobe, oceanStor, "", writer)
	}
	if err != nil {
		c.errCreateStorageProvisioner(cluster.Name, provisioner, fmt.Errorf("unmarshal provisioner.Vars error %s", err.Error()))
		return
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
			if err := c.sync(clusterName, provisioner); err != nil {
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

func (c clusterStorageProvisionerService) sync(clusterName string, provisioner dto.ClusterStorageProvisionerSync) error {
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return err
	}
	endpoints, err := c.clusterService.GetApiServerEndpoints(clusterName)
	client, err := kubernetes.NewKubernetesClient(&kubernetes.Config{
		Token: secret.KubernetesToken,
		Hosts: endpoints,
	})
	if err != nil {
		return err
	}
	switch provisioner.Type {
	case "external-ceph":
		_, err := client.AppsV1().Deployments("kube-system").Get(context.TODO(), "external-ceph", metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
		}
	case "nfs":
		_, err := client.AppsV1().Deployments("kube-system").Get(context.TODO(), provisioner.Name, metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
		}
	case "vsphere":
		_, err := client.AppsV1().StatefulSets("kube-system").Get(context.TODO(), "vsphere-csi-controller", metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
		}
	case "rook-ceph":
		_, err := client.AppsV1().Deployments("kube-system").Get(context.TODO(), "rook-ceph-operator", metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
		}
	case "oceanstor":
		_, err := client.AppsV1().Deployments("kube-system").Get(context.TODO(), "huawei-csi-controller", metav1.GetOptions{})
		if err != nil && checkError(err) {
			return err
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
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return err
	}
	endpoints, err := c.clusterService.GetApiServerEndpoints(clusterName)
	client, err := kubernetes.NewKubernetesClient(&kubernetes.Config{
		Token: secret.KubernetesToken,
		Hosts: endpoints,
	})
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
