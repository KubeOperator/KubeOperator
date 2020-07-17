package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider/client"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
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
		PlanRepo:       repository.NewPlanRepository(),
	}
}

type clusterNodeService struct {
	ClusterService ClusterService
	NodeRepo       repository.ClusterNodeRepository
	HostRepo       repository.HostRepository
	PlanRepo       repository.PlanRepository
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
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	cluster.Cluster.Plan, err = c.PlanRepo.GetById(cluster.PlanID)
	if err != nil {
		return nil, err
	}
	ch := make(chan int)
	var mNodes []*model.ClusterNode
	switch cluster.Spec.Provider {
	case constant.ClusterProviderBareMetal:
		mNodes, err = c.doBareMetalCreateNodes(cluster.Cluster, item)
		if err != nil {
			return nil, err
		}
	case constant.ClusterProviderPlan:
		mNodes, err = c.doPlanCreateNodes(cluster.Cluster, item)
		if err != nil {
			return nil, err
		}
		nodes, _ := c.NodeRepo.List(cluster.Name)
		var hosts []*model.Host
		for i, _ := range nodes {
			hosts = append(hosts, &nodes[i].Host)
		}
		c.doCreateHosts(ch, cluster.Cluster, hosts)
	}
	var nodes []dto.Node
	for _, n := range mNodes {
		c.doCreate(ch, *n, cluster.Cluster)
		nodes = append(nodes, dto.Node{ClusterNode: *n})
	}
	return nodes, nil
}

func (c clusterNodeService) doBareMetalCreateNodes(cluster model.Cluster, item dto.NodeBatch) ([]*model.ClusterNode, error) {
	var hosts []*model.Host
	for _, h := range item.Hosts {
		host, err := c.HostRepo.Get(h)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, &host)
	}
	return c.createNodes(cluster, hosts)

}

func (c clusterNodeService) doPlanCreateNodes(cluster model.Cluster, item dto.NodeBatch) ([]*model.ClusterNode, error) {
	hosts, err := c.createPlanHosts(cluster, item.Increase)
	var hs []*model.Host
	if err != nil {
		return nil, err
	}
	for _, h := range hosts {
		host, err := c.HostRepo.Get(h.Name)
		if err != nil {
			return nil, err
		}
		hs = append(hs, &host)
	}
	return c.createNodes(cluster, hs)
}

func (c clusterNodeService) createPlanHosts(cluster model.Cluster, increase int) ([]*model.Host, error) {
	var hosts []*model.Host
	hash := map[string]interface{}{}
	for _, node := range cluster.Nodes {
		hosts = append(hosts, &node.Host)
		hash[node.Host.Name] = nil
	}
	var newHosts []*model.Host
	for i := 0; i < increase; i++ {
		var name string
		for k := 0; k < increase+len(hosts); k++ {
			n := fmt.Sprintf("%s-worker-%d", cluster.Name, k+1)
			if _, ok := hash[n]; !ok {
				name = n
				break
			}
		}
		newHosts = append(newHosts, &model.Host{
			Name: name,
			Port: 22,
		})
	}
	group := allocateZone(cluster.Plan.Zones, newHosts)
	var selectedIps []string
	for k, v := range group {
		providerVars := map[string]interface{}{}
		providerVars["provider"] = cluster.Plan.Region.Provider
		_ = json.Unmarshal([]byte(cluster.Plan.Region.Vars), &providerVars)
		cloudClient := client.NewCloudClient(providerVars)
		err := allocateIpAddr(cloudClient, *k, v, selectedIps)
		if err != nil {
			return nil, err
		}
	}
	_ = c.HostRepo.BatchSave(newHosts)
	return newHosts, nil
}

func (c clusterNodeService) doCreateHosts(ch chan int, cluster model.Cluster, hosts []*model.Host) {
	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	_ = doInit(k, cluster.Plan, hosts)
}

func (c clusterNodeService) createNodes(cluster model.Cluster, hosts []*model.Host) ([]*model.ClusterNode, error) {
	var mNodes []*model.ClusterNode
	ns, err := c.NodeRepo.List(cluster.Name)
	if err != nil {
		return nil, err
	}
	hash := map[string]interface{}{}
	for _, n := range ns {
		hash[n.Name] = nil
	}
	for _, host := range hosts {
		host.ClusterID = cluster.ID
		err = c.HostRepo.Save(host)
		if err != nil {
			return nil, err
		}
		var name string
		for i := 1; i < len(ns)+len(hosts); i++ {
			name = fmt.Sprintf("%s-%d", constant.NodeRoleNameWorker, i)
			if _, ok := hash[name]; ok {
				continue
			}
			break
		}
		n := model.ClusterNode{
			Name:      name,
			ClusterID: cluster.ID,
			HostID:    host.ID,
			Role:      constant.NodeRoleNameWorker,
			Status:    constant.ClusterWaiting,
			Host:      *host,
		}
		mNodes = append(mNodes, &n)
	}
	err = c.NodeRepo.BatchSave(mNodes)
	if err != nil {
		return nil, err
	}
	return mNodes, err
}

const deleteWorkerPlaybook = "96-remove-worker.yml"

func (c clusterNodeService) doDelete(worker model.ClusterNode, clusterName string) {
	cluster, _ := c.ClusterService.Get(clusterName)
	worker.Status = constant.ClusterTerminating
	_ = c.NodeRepo.Save(&worker)
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

func (c clusterNodeService) doCreate(ch chan int, worker model.ClusterNode, cluster model.Cluster) {
	cluster.Nodes, _ = c.NodeRepo.List(cluster.Name)
	worker.Status = constant.ClusterInitializing
	_ = c.NodeRepo.Save(&worker)
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
