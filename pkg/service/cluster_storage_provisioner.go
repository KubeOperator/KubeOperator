package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/plugin/storage"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
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
					db.DB.Model(model.ClusterStorageProvisioner{}).Save(&p)
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

	p := getPhase(provisioner)
	if err := p.Run(admCluster.Kobe, writer); err != nil {
		provisioner.Status = constant.ClusterFailed
		provisioner.Message = err.Error()
	} else {
		provisioner.Status = constant.ClusterRunning
	}
	_ = c.provisionerRepo.Save(cluster.Name, &provisioner)
}

func getPhase(provisioner model.ClusterStorageProvisioner) phases.Interface {
	vars := map[string]interface{}{}
	_ = json.Unmarshal([]byte(provisioner.Vars), &vars)
	var p phases.Interface
	switch provisioner.Type {
	case "nfs":
		p = &storage.NfsStoragePhase{
			NfsServer:        fmt.Sprintf("%v", vars["storage_nfs_server"]),
			NfsServerPath:    fmt.Sprintf("%v", vars["storage_nfs_server_path"]),
			NfsServerVersion: fmt.Sprintf("%v", vars["storage_nfs_server_version"]),
			ProvisionerName:  provisioner.Name,
		}
	case "rook-ceph":
		p = &storage.RookCephStoragePhase{
			StorageRookPath: fmt.Sprintf("%v", vars["storage_rook_path"]),
		}
	case "vsphere":
		p = &storage.VsphereStoragePhase{
			VcUsername: fmt.Sprintf("%v", vars["vc_username"]),
			VcPassword: fmt.Sprintf("%v", vars["vc_password"]),
			VcHost:     fmt.Sprintf("%v", vars["vc_host"]),
			VcPort:     fmt.Sprintf("%v", vars["vc_port"]),
			Datacenter: fmt.Sprintf("%v", vars["datacenter"]),
			Datastore:  fmt.Sprintf("%v", vars["datastore"]),
			Folder:     fmt.Sprintf("%v", vars["folder"]),
		}
	case "external-ceph":
		p = &storage.ExternalCephStoragePhase{
			ProvisionerName: provisioner.Name,
		}
	case "oceanstor":
		p = &storage.OceanStorPhase{
			OceanStorType:           fmt.Sprintf("%v", vars["oceanstor_type"]),
			OceanstorProduct:        fmt.Sprintf("%v", vars["oceanstor_product"]),
			OceanstorURLs:           fmt.Sprintf("%v", vars["oceanstor_urls"]),
			OceanstorUser:           fmt.Sprintf("%v", vars["oceanstor_user"]),
			OceanstorPassword:       fmt.Sprintf("%v", vars["oceanstor_password"]),
			OceanstorPools:          fmt.Sprintf("%v", vars["oceanstor_pools"]),
			OceanstorPortal:         fmt.Sprintf("%v", vars["oceanstor_portal"]),
			OceanstorControllerType: fmt.Sprintf("%v", vars["oceanstor_controller_type"]),
			OceanstorIsMultipath:    fmt.Sprintf("%v", vars["oceanstor_is_multipath"]),
		}
	}

	return p
}
