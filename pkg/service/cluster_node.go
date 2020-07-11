package service

import (
	"context"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterNodeService interface {
	List(clusterName string) ([]dto.Node, error)
	Batch(clusterName string, batch dto.NodeBatch) ([]dto.Node, error)
}

var log = logger.Default

func NewClusterNodeService() ClusterNodeService {
	return &clusterNodeService{
		ClusterService: NewClusterService(),
		NodeRepo:       repository.NewClusterNodeRepository(),
		HostRepo:       repository.NewHostRepository(),
	}
}

type clusterNodeService struct {
	ClusterService ClusterService
	NodeRepo       repository.ClusterNodeRepository
	HostRepo       repository.HostRepository
}

func (c clusterNodeService) List(clusterName string) ([]dto.Node, error) {
	var nodes []dto.Node
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	cluster.Nodes, err = c.NodeRepo.List(cluster.Name)
	if err != nil {
		return nil, err
	}
	endpoint, err := c.ClusterService.GetApiServerEndpoint(clusterName)
	if err != nil {
		return nil, err
	}
	secret, err := c.ClusterService.GetSecrets(clusterName)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Host:  endpoint.Address,
		Token: secret.KubernetesToken,
		Port:  endpoint.Port,
	})
	if err != nil {
		return nil, err
	}

	kubeNodes, err := kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range cluster.Nodes {
		n := dto.Node{
			ClusterNode: node,
		}
		if node.Status == constant.ClusterRunning {
			for _, kn := range kubeNodes.Items {
				if node.Name == kn.Name {
					n.Info = kn
				}
			}
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func (c clusterNodeService) Batch(clusterName string, item dto.NodeBatch) ([]dto.Node, error) {
	switch item.Operation {
	case constant.BatchOperationCreate:
		return c.batchCreate(clusterName, item)
	case constant.BatchOperationDelete:
		return c.batchDelete(clusterName, item)
	}
	return nil, nil
}

func (c clusterNodeService) batchDelete(clusterName string, item dto.NodeBatch) ([]dto.Node, error) {
	var nodes []dto.Node
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	if cluster.Spec.Provider == constant.ClusterProviderBareMetal {
		for _, nodeName := range item.Nodes {
			n, err := c.NodeRepo.Get(clusterName, nodeName)
			nodes = append(nodes, dto.Node{ClusterNode: n})
			if err != nil {
				return nil, err
			}
			if n.Status == constant.ClusterRunning {
				go c.doDelete(n, clusterName)
			} else {
				_ = c.NodeRepo.Delete(n.ID)
			}
		}
	}
	return nodes, nil
}

func (c clusterNodeService) batchCreate(clusterName string, item dto.NodeBatch) ([]dto.Node, error) {
	var nodes []dto.Node
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	var mNodes []*model.ClusterNode

	ns, err := c.NodeRepo.List(clusterName)
	if err != nil {
		return nil, err
	}
	hash := map[string]interface{}{}
	for _, n := range ns {
		hash[n.Name] = nil
	}
	if cluster.Spec.Provider == constant.ClusterProviderBareMetal {
		for _, host := range item.Hosts {
			h, err := c.HostRepo.Get(host)
			if err != nil {
				return nil, err
			}
			h.ClusterID = cluster.ID
			err = c.HostRepo.Save(&h)
			if err != nil {
				return nil, err
			}
			var name string
			for i := 1; i < len(ns)+len(item.Hosts); i++ {
				name = fmt.Sprintf("%s-%d", constant.NodeRoleNameWorker, i)
				if _, ok := hash[name]; ok {
					continue
				}
				break
			}
			n := model.ClusterNode{
				Name:      name,
				ClusterID: cluster.ID,
				HostID:    h.ID,
				Role:      constant.NodeRoleNameWorker,
				Status:    constant.ClusterWaiting,
			}
			mNodes = append(mNodes, &n)
		}
		if err := c.NodeRepo.BatchSave(mNodes); err != nil {
			return nil, err
		}
	}
	for _, n := range mNodes {
		go c.doCreate(*n, clusterName)
		nodes = append(nodes, dto.Node{ClusterNode: *n})
	}
	return nodes, nil
}

const deleteWorkerPlaybook = "96-remove-worker.yml"

func (c clusterNodeService) doDelete(worker model.ClusterNode, clusterName string) {
	worker.Status = constant.ClusterTerminating
	_ = c.NodeRepo.Save(&worker)
	cluster, _ := c.ClusterService.Get(clusterName)
	inventory := cluster.ParseInventory()
	for i, _ := range inventory.Groups {
		if inventory.Groups[i].Name == "del-worker" {
			inventory.Groups[i].Hosts = append(inventory.Groups[i].Hosts, worker.Name)
		}
	}
	k := kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})
	for name, _ := range facts.DefaultFacts {
		k.SetVar(name, facts.DefaultFacts[name])
	}
	log.Debugf("start run delete worker: %s", worker.Name)
	_ = phases.RunPlaybookAndGetResult(k, deleteWorkerPlaybook)
	worker.Status = constant.ClusterTerminated
	_ = c.NodeRepo.Save(&worker)
	_ = c.NodeRepo.Delete(worker.ID)
}

const addWorkerPlaybook = "91-add-worker.yml"

func (c clusterNodeService) doCreate(worker model.ClusterNode, clusterName string) {
	worker.Status = constant.ClusterInitializing
	_ = c.NodeRepo.Save(&worker)
	cluster, _ := c.ClusterService.Get(clusterName)
	inventory := cluster.ParseInventory()
	for i, _ := range inventory.Groups {
		if inventory.Groups[i].Name == "new-worker" {
			inventory.Groups[i].Hosts = append(inventory.Groups[i].Hosts, worker.Name)
		}
	}
	k := kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})
	for name, _ := range facts.DefaultFacts {
		k.SetVar(name, facts.DefaultFacts[name])
	}
	//if cluster.Spec.NetworkType != "" {
	//	k.SetVar(facts.NetworkPluginFactName, cluster.Spec.NetworkType)
	//}
	//if cluster.Spec.RuntimeType != "" {
	//	k.SetVar(facts.ContainerRuntimeFactName, cluster.Spec.RuntimeType)
	//}
	//if cluster.Spec.DockerStorageDir != "" {
	//	k.SetVar(facts.DockerStorageDirFactName, cluster.Spec.DockerStorageDir)
	//}
	//if cluster.Spec.ContainerdStorageDir != "" {
	//	k.SetVar(facts.ContainerdStorageDirFactName, cluster.Spec.ContainerdStorageDir)
	//}
	//if cluster.Spec.LbKubeApiserverIp != "" {
	//	k.SetVar(facts.LbKubeApiserverPortFactName, cluster.Spec.LbKubeApiserverIp)
	//}
	//if cluster.Spec.KubePodSubnet != "" {
	//	k.SetVar(facts.KubePodSubnetFactName, cluster.Spec.KubePodSubnet)
	//}
	//if cluster.Spec.KubeServiceSubnet != "" {
	//	k.SetVar(facts.KubeServiceSubnetFactName, cluster.Spec.KubeServiceSubnet)
	//}
	//secret, _ := c.ClusterService.GetSecrets(clusterName)
	//k.SetVar(facts.KubeadmTokenFactName, secret.KubernetesToken)
	log.Debugf("start run add worker: %s", worker.Name)
	err := phases.RunPlaybookAndGetResult(k, addWorkerPlaybook)
	if err != nil {
		worker.Status = constant.ClusterFailed
		worker.Message = err.Error()
		_ = c.NodeRepo.Save(&worker)
		return
	}
	worker.Status = constant.ClusterRunning
	_ = c.NodeRepo.Save(&worker)
}
