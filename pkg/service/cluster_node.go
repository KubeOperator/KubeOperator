package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider/client"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
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
	"sync"
	"time"
)

type ClusterNodeService interface {
	List(clusterName string) ([]dto.Node, error)
	Batch(clusterName string, batch dto.NodeBatch) ([]dto.Node, error)
}

var log = logger.Default

func NewClusterNodeService() ClusterNodeService {
	return &clusterNodeService{
		ClusterService:      NewClusterService(),
		NodeRepo:            repository.NewClusterNodeRepository(),
		HostRepo:            repository.NewHostRepository(),
		PlanRepo:            repository.NewPlanRepository(),
		systemSettingRepo:   repository.NewSystemSettingRepository(),
		clusterLogService:   NewClusterLogService(),
		projectResourceRepo: repository.NewProjectResourceRepository(),
	}
}

type clusterNodeService struct {
	ClusterService      ClusterService
	NodeRepo            repository.ClusterNodeRepository
	HostRepo            repository.HostRepository
	PlanRepo            repository.PlanRepository
	systemSettingRepo   repository.SystemSettingRepository
	clusterLogService   ClusterLogService
	projectResourceRepo repository.ProjectResourceRepository
}

type nodeMessage struct {
	node    *model.ClusterNode
	message string
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
	// 判断是否存在正在运行的节点变更任务
	nodes, _ := c.NodeRepo.List(clusterName)
	for _, node := range nodes {
		if node.Status != constant.ClusterRunning {
			return nil, errors.New("NODE_ALREADY_RUNNING_TASK")
		}
	}
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
	if cluster.Spec.Provider == constant.ClusterProviderPlan {
		cluster.Cluster.Plan, err = c.PlanRepo.GetById(cluster.PlanID)
		if err != nil {
			return nil, err
		}
	}
	var needDeleteNodes []*model.ClusterNode
	for _, nodeName := range item.Nodes {
		n, err := c.NodeRepo.Get(clusterName, nodeName)
		nodes = append(nodes, dto.Node{ClusterNode: n})
		if err != nil {
			return nil, err
		}
		if n.Status == constant.ClusterRunning {
			needDeleteNodes = append(needDeleteNodes, &n)
		} else {
			_ = c.NodeRepo.Delete(n.ID)
		}
		go c.doDelete(&cluster.Cluster, needDeleteNodes)

	}
	return nodes, nil
}

func (c *clusterNodeService) doDelete(cluster *model.Cluster, nodes []*model.ClusterNode) {
	var clog model.ClusterLog
	clog.Type = constant.ClusterLogTypeDeleteNode
	clog.StartTime = time.Now()
	clog.EndTime = time.Now()
	err := c.clusterLogService.Save(cluster.Name, &clog)
	if err != nil {
		log.Error(err)
	}
	err = c.clusterLogService.Start(&clog)
	if err != nil {
		log.Error(err)
	}
	wg := sync.WaitGroup{}
	for i := range nodes {
		nodes[i].Status = constant.ClusterTerminating
		db.DB.Save(&nodes[i])
		go c.doSingleDelete(&wg, cluster, nodes[i])
		wg.Add(1)
	}
	wg.Wait()
	if cluster.Spec.Provider == constant.ClusterProviderPlan {
		err := c.destroyHosts(cluster, nodes)
		if err != nil {
			log.Debug(err)
		}
	}
	for i := range nodes {
		if cluster.Spec.Provider == constant.ClusterProviderBareMetal {
			nodes[i].Host.ClusterID = ""
			_ = c.HostRepo.Save(&nodes[i].Host)
		}
		if cluster.Spec.Provider == constant.ClusterProviderPlan {
			db.DB.Delete(model.ClusterNode{ID: nodes[i].ID})
			db.DB.Delete(model.Host{ID: nodes[i].HostID})
			hostResources, err := c.projectResourceRepo.ListByResourceIdAndType(nodes[i].HostID, constant.ResourceHost)
			if err != nil {
				log.Error(err)
			}
			if len(hostResources) > 0 {
				db.DB.Delete(model.ProjectResource{ID: hostResources[0].ID})
			}
		}
		_ = c.NodeRepo.Delete(nodes[i].ID)
	}
	e := c.clusterLogService.End(&clog, true, "")
	if e != nil {
		log.Error(e)
	}
}

func (c *clusterNodeService) destroyHosts(cluster *model.Cluster, nodes []*model.ClusterNode) error {
	var aliveHosts []*model.Host
	for i := range cluster.Nodes {
		flag := false
		for k := range nodes {
			if cluster.Nodes[i].Name == nodes[k].Name {
				flag = true
			}
		}
		if !flag {
			aliveHosts = append(aliveHosts, &cluster.Nodes[i].Host)
		}
	}
	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	return doInit(k, cluster.Plan, aliveHosts)
}

func (c clusterNodeService) batchCreate(clusterName string, item dto.NodeBatch) ([]dto.Node, error) {
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	if cluster.Spec.Provider == constant.ClusterProviderPlan {
		cluster.Cluster.Plan, err = c.PlanRepo.GetById(cluster.PlanID)
		if err != nil {
			return nil, err
		}
	}
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
	}
	go c.doCreate(&cluster.Cluster, mNodes)
	var nodes []dto.Node
	for _, n := range mNodes {
		nodes = append(nodes, dto.Node{ClusterNode: *n})
	}
	return nodes, nil
}

func (c *clusterNodeService) doCreate(cluster *model.Cluster, nodes []*model.ClusterNode) {
	var clog model.ClusterLog
	clog.Type = constant.ClusterLogTypeAddNode
	clog.StartTime = time.Now()
	clog.EndTime = time.Now()
	err := c.clusterLogService.Save(cluster.Name, &clog)
	if err != nil {
		log.Error(err)
	}
	err = c.clusterLogService.Start(&clog)
	if err != nil {
		log.Error(err)
	}
	if cluster.Spec.Provider == constant.ClusterProviderPlan {
		allNodes, _ := c.NodeRepo.List(cluster.Name)
		var allHosts []*model.Host
		for i, _ := range allNodes {
			allHosts = append(allHosts, &allNodes[i].Host)
		}
		err := c.doCreateHosts(cluster, allHosts)
		if err != nil {
			e := c.clusterLogService.End(&clog, false, err.Error())
			if e != nil {
				log.Error(e)
			}
			for i := range nodes {
				db.DB.Delete(model.ClusterNode{ID: nodes[i].ID})
				db.DB.Delete(model.Host{ID: nodes[i].HostID})
			}
			return
		} else {
			for i := range nodes {
				nodes[i].Host.Status = constant.ClusterRunning
				err := db.DB.Save(&nodes[i].Host).Error
				if err != nil {
					log.Error(err)
				}
				//add project resource
				clusterResources, err := c.projectResourceRepo.ListByResourceIdAndType(cluster.ID, constant.ResourceCluster)
				if err != nil {
					log.Error(err)
				}
				if len(clusterResources) > 0 {
					db.DB.Create(&model.ProjectResource{
						ResourceId:   nodes[i].Host.ID,
						ResourceType: constant.ResourceHost,
						ProjectID:    clusterResources[0].ProjectID,
					})
				}
			}
		}
	}
	var waitGroup sync.WaitGroup
	var nms []*nodeMessage
	for i := range nodes {
		nodes[i].Status = constant.ClusterInitializing
		_ = c.NodeRepo.Save(nodes[i])
		nm := &nodeMessage{
			node: nodes[i],
		}
		nms = append(nms, nm)
		go c.doSingleNodeCreate(&waitGroup, cluster, nm)
		waitGroup.Add(1)
	}
	waitGroup.Wait()
	success := true
	mergedLogMap := make(map[string]string)
	for i := range nms {
		err := db.DB.Save(nms[i].node).Error
		if err != nil {
			log.Error(err)
		}
		if nms[i].node.Status != constant.ClusterRunning {
			success = false
			mergedLogMap[nms[i].node.Name] = nms[i].message
		}
	}
	if success {
		e := c.clusterLogService.End(&clog, true, "")
		if e != nil {
			log.Error(e)
		}
	} else {
		buf, _ := json.Marshal(&mergedLogMap)
		e := c.clusterLogService.End(&clog, false, string(buf))
		if e != nil {
			log.Error(e)
		}
		for i := range nodes {
			if nodes[i].Status == constant.ClusterRunning {
				nodes[i].Status = constant.ClusterTerminating
				_ = c.NodeRepo.Save(nodes[i])
				go c.doSingleDelete(&waitGroup, cluster, nodes[i])
				waitGroup.Add(1)
			}
		}
		waitGroup.Wait()
		if cluster.Spec.Provider == constant.ClusterProviderBareMetal {
			for i := range nodes {
				db.DB.Delete(nodes[i])
			}
		} else {
			nos, _ := c.NodeRepo.List(cluster.Name)
			cluster.Nodes = nos
			_ = c.destroyHosts(cluster, nodes)
			for i := range nodes {
				db.DB.Delete(model.ClusterNode{ID: nodes[i].ID})
				db.DB.Delete(model.Host{ID: nodes[i].HostID})
			}
		}
	}
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
				hash[name] = nil
				break
			}
		}
		newHosts = append(newHosts, &model.Host{
			Name:   name,
			Port:   22,
			Status: constant.ClusterCreating,
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
	err := c.HostRepo.BatchSave(newHosts)
	if err != nil {
		log.Error(err)
	}
	return newHosts, nil
}

func (c clusterNodeService) doCreateHosts(cluster *model.Cluster, hosts []*model.Host) error {
	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	return doInit(k, cluster.Plan, hosts)
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
			name = fmt.Sprintf("%s-%s-%d", cluster.Name, constant.NodeRoleNameWorker, i)
			if _, ok := hash[name]; ok {
				continue
			}
			hash[name] = nil
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

func (c *clusterNodeService) doSingleDelete(wg *sync.WaitGroup, cluster *model.Cluster, worker *model.ClusterNode) {
	defer wg.Done()
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
	clusterVars := cluster.GetKobeVars()
	for j, v := range clusterVars {
		k.SetVar(j, v)
	}
	k.SetVar(facts.ClusterNameFactName, cluster.Name)
	val, _ := c.systemSettingRepo.Get("ip")
	k.SetVar(facts.LocalHostnameFactName, val.Value)
	_ = phases.RunPlaybookAndGetResult(k, deleteWorkerPlaybook)
	worker.Status = constant.ClusterTerminated
}

const addWorkerPlaybook = "91-add-worker.yml"

func (c clusterNodeService) doSingleNodeCreate(waitGroup *sync.WaitGroup, cluster *model.Cluster, nm *nodeMessage) {
	defer waitGroup.Done()
	cluster.Nodes, _ = c.NodeRepo.List(cluster.Name)

	inventory := cluster.ParseInventory()
	for i, _ := range inventory.Groups {
		if inventory.Groups[i].Name == "new-worker" {
			inventory.Groups[i].Hosts = append(inventory.Groups[i].Hosts, nm.node.Name)
		}
	}
	k := kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})
	for name, _ := range facts.DefaultFacts {
		k.SetVar(name, facts.DefaultFacts[name])
	}
	clusterVars := cluster.GetKobeVars()
	for j, v := range clusterVars {
		k.SetVar(j, v)
	}
	k.SetVar(facts.ClusterNameFactName, cluster.Name)
	val, _ := c.systemSettingRepo.Get("ip")
	k.SetVar(facts.LocalHostnameFactName, val.Value)
	err := phases.RunPlaybookAndGetResult(k, addWorkerPlaybook)
	if err != nil {
		nm.node.Status = constant.ClusterFailed
		nm.message = err.Error()
	}
	nm.node.Status = constant.ClusterRunning
}
