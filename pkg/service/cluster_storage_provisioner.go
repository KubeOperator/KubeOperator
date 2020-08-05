package service

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/plugin/storage"
)

type ClusterStorageProvisionerService interface {
	ListStorageProvisioner(clusterName string) ([]dto.ClusterStorageProvisioner, error)
	CreateStorageProvisioner(clusterName string, creation dto.ClusterStorageProvisionerCreation) (dto.ClusterStorageProvisioner, error)
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
func (c clusterStorageProvisionerService) DeleteStorageProvisioner(clusterName string, provisioner string) error {
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
		Status: constant.ClusterWaiting,
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
	p := getPhase(provisioner)
	err := p.Run(admCluster.Kobe)
	if err != nil {
		provisioner.Status = constant.ClusterFailed
		provisioner.Message = err.Error()
	} else {
		provisioner.Status = constant.ClusterRunning
	}
	_ = c.provisionerRepo.Save(cluster.Name, &provisioner)

}

func getPhase(provisioner model.ClusterStorageProvisioner) phases.Interface {
	vars := map[string]string{}
	_ = json.Unmarshal([]byte(provisioner.Vars), &vars)
	var p phases.Interface
	switch provisioner.Type {
	case "nfs":
		p = &storage.NfsStoragePhase{
			NfsServer:        vars["storage_nfs_server"],
			NfsServerPath:    vars["storage_nfs_server_path"],
			NfsServerVersion: vars["storage_nfs_server_version"],
			ProvisionerName:  provisioner.Name,
		}
	case "rook-ceph":
		p = &storage.RookCephStoragePhase{
			StorageRookPath: vars["storage_rook_path"],
		}
	case "vsphere":
		p = &storage.VsphereStoragePhase{
			VcUsername: vars["username"],
			VcPassword: vars["password"],
			VcHost:     vars["host"],
			VcPort:     vars["port"],
			Datacenter: vars["datacenter"],
			Datastore:  vars["datastore"],
			Folder:     vars["folder"],
		}
	case "external-ceph":
		p = &storage.ExternalCephStoragePhase{
			CephMonitor:               vars["ceph_monitor"],
			CephOsdPool:               vars["ceph_osd_pool"],
			CephAdminId:               vars["ceph_admin_id"],
			CephAdminSecret:           vars["ceph_admin_secret"],
			CephUserId:                vars["ceph_user_id"],
			CephUserSecret:            vars["ceph_user_secret"],
			CephFsType:                vars["ceph_fsType"],
			CephImageFormat:           vars["ceph_imageFormat"],
			StorageRbdProvisionerName: vars["storage_rbd_provisioner_name"],
			ProvisionerName:           provisioner.Name,
		}
	}
	return p
}
