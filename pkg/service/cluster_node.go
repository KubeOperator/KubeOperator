package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterNodeService interface {
	Get(clusterName, name string) (*dto.Node, error)
	List(clusterName string) ([]dto.Node, error)
	Batch(clusterName string, batch dto.NodeBatch) ([]dto.Node, error)
	Page(num, size int, clusterName string) (*dto.NodePage, error)
}

var log = logger.Default

func NewClusterNodeService() ClusterNodeService {
	return &clusterNodeService{
		ClusterService:      NewClusterService(),
		NodeRepo:            repository.NewClusterNodeRepository(),
		HostRepo:            repository.NewHostRepository(),
		PlanRepo:            repository.NewPlanRepository(),
		systemSettingRepo:   repository.NewSystemSettingRepository(),
		projectResourceRepo: repository.NewProjectResourceRepository(),
		messageService:      NewMessageService(),
		vmConfigRepo:        repository.NewVmConfigRepository(),
	}
}

type clusterNodeService struct {
	ClusterService      ClusterService
	NodeRepo            repository.ClusterNodeRepository
	HostRepo            repository.HostRepository
	PlanRepo            repository.PlanRepository
	systemSettingRepo   repository.SystemSettingRepository
	projectResourceRepo repository.ProjectResourceRepository
	messageService      MessageService
	vmConfigRepo        repository.VmConfigRepository
}

func (c *clusterNodeService) Get(clusterName, name string) (*dto.Node, error) {
	var n model.ClusterNode
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}

	err = db.DB.Where(model.ClusterNode{ClusterID: cluster.ID, Name: name}).Find(&n).Error
	if err != nil {
		return nil, err
	}
	return &dto.Node{
		ClusterNode: n,
	}, nil
}

func (c clusterNodeService) Page(num, size int, clusterName string) (*dto.NodePage, error) {
	var nodes []dto.Node
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	count, mNodes, err := c.NodeRepo.Page(num, size, cluster.Name)
	if err != nil {
		return nil, err
	}

	secret, err := c.ClusterService.GetSecrets(clusterName)
	if err != nil {
		return nil, err
	}

	endpoints, err := c.ClusterService.GetApiServerEndpoints(clusterName)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Hosts: endpoints,
		Token: secret.KubernetesToken,
	})
	if err != nil {
		return nil, err
	}
	kubeNodes, err := kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
exit:
	for _, node := range mNodes {
		n := dto.Node{
			ClusterNode: node,
			Ip:          node.Host.Ip,
		}
		for _, kn := range kubeNodes.Items {
			if node.Name == kn.Name {
				n.Info = kn
				nodes = append(nodes, n)
				continue exit
			}
		}
		if n.Status == constant.StatusRunning {
			n.Status = constant.StatusLost
			go func() {
				if err := db.DB.Save(&n.ClusterNode).Error; err != nil {
					log.Error(err)
				}
			}()
		}
		nodes = append(nodes, n)
	}
	return &dto.NodePage{
		Items: nodes,
		Total: count,
	}, nil
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
	secret, err := c.ClusterService.GetSecrets(clusterName)
	if err != nil {
		return nil, err
	}
	endpoints, err := c.ClusterService.GetApiServerEndpoints(clusterName)
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Hosts: endpoints,
		Token: secret.KubernetesToken,
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
		if node.Status != constant.ClusterRunning && node.Status != constant.ClusterFailed && node.Status != constant.StatusLost {
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
	for i := range nodes {
		nodes[i].Status = constant.ClusterTerminating
		db.DB.Save(&nodes[i])
	}
	err := c.doNodeDelete(cluster, nodes)
	if err != nil {
		log.Error(err.Error())
	}
	if cluster.Spec.Provider == constant.ClusterProviderPlan {
		err := c.destroyHosts(cluster, nodes)
		if err != nil {
			log.Error(err)
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
			hostResources, err := c.projectResourceRepo.ListByResourceIDAndType(nodes[i].HostID, constant.ResourceHost)
			if err != nil {
				log.Error(err)
			}
			if len(hostResources) > 0 {
				db.DB.Delete(model.ProjectResource{ID: hostResources[0].ID})
			}
		}
		_ = c.NodeRepo.Delete(nodes[i].ID)
	}
	_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterRemoveWorker, true, ""), cluster.Name, constant.ClusterRemoveWorker)
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
	var mNodes []*model.ClusterNode
	switch cluster.Spec.Provider {
	case constant.ClusterProviderBareMetal:
		mNodes, err = c.doBareMetalCreateNodes(cluster.Cluster, item)
		if err != nil {
			return nil, err
		}
	case constant.ClusterProviderPlan:
		cluster.Cluster.Plan, err = c.PlanRepo.GetById(cluster.PlanID)
		if err != nil {
			return nil, err
		}
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
	switch cluster.Spec.Provider {
	case constant.ClusterProviderPlan:
		allNodes, _ := c.NodeRepo.List(cluster.Name)
		var allHosts []*model.Host
		for i := range allNodes {
			allHosts = append(allHosts, &allNodes[i].Host)
		}
		err := c.doCreateHosts(cluster, allHosts)
		// 处理创建主机错误
		if err != nil {
			// 删除主机和节点
			for i := range nodes {
				db.DB.Delete(model.ClusterNode{ID: nodes[i].ID})
				db.DB.Delete(model.Host{ID: nodes[i].HostID})
			}
			_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterAddWorker, false, err.Error()), cluster.Name, constant.ClusterAddWorker)
			return
		}
		for i := range nodes {
			nodes[i].Host.Status = constant.ClusterRunning
			err := db.DB.Save(&nodes[i].Host).Error
			if err != nil {
				log.Error(err)
			}
			clusterResources, err := c.projectResourceRepo.ListByResourceIDAndType(cluster.ID, constant.ResourceCluster)
			if err != nil {
				log.Error(err)
			}
			if len(clusterResources) > 0 {
				if err := db.DB.Create(&model.ProjectResource{
					ResourceID:   nodes[i].Host.ID,
					ResourceType: constant.ResourceHost,
					ProjectID:    clusterResources[0].ProjectID,
				}).Error; err != nil {
					log.Error(err.Error())
				}
			}
		}
	}
	// 初始化init状态
	for i := range nodes {
		nodes[i].Status = constant.ClusterInitializing
		if err := c.NodeRepo.Save(nodes[i]); err != nil {
			log.Error(err.Error())
		}
	}
	if err := c.doNodeCreate(cluster, nodes); err != nil {
		for i := range nodes {
			nodes[i].Status = constant.ClusterFailed
			if err := c.NodeRepo.Save(nodes[i]); err != nil {
				log.Error(err.Error())
			}
		}
		for i := range nodes {
			nodes[i].Status = constant.ClusterTerminating
			if err := c.NodeRepo.Save(nodes[i]); err != nil {
				log.Error(err.Error())
			}
		}
		if err := c.doNodeDelete(cluster, nodes); err != nil {
			log.Error(err.Error())
		}
		switch cluster.Spec.Provider {
		case constant.ClusterProviderBareMetal:
			for i := range nodes {
				db.DB.Model(model.Host{}).Where(model.Host{ID: nodes[i].HostID}).Updates(map[string]interface{}{"ClusterID": ""})
				db.DB.Delete(nodes[i])
			}
		case constant.ClusterProviderPlan:
			nos, _ := c.NodeRepo.List(cluster.Name)
			cluster.Nodes = nos
			if e := c.destroyHosts(cluster, nodes); e != nil {
				log.Error(e.Error())
			}
			for i := range nodes {
				db.DB.Delete(model.ClusterNode{ID: nodes[i].ID})
				db.DB.Delete(model.Host{ID: nodes[i].HostID})
			}
		}
		if e := c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterAddWorker, false, err.Error()), cluster.Name, constant.ClusterAddWorker); e != nil {
			log.Error(e.Error())
		}
		return
	}
	for i := range nodes {
		nodes[i].Status = constant.ClusterRunning
		if err := c.NodeRepo.Save(nodes[i]); err != nil {
			log.Error(err.Error())
		}
	}
	if e := c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterAddWorker, true, ""), cluster.Name, constant.ClusterAddWorker); e != nil {
		log.Error(e.Error())
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
		newHost := &model.Host{
			Name:   name,
			Port:   22,
			Status: constant.ClusterCreating,
		}
		if cluster.Plan.Region.Provider != constant.OpenStack {
			planVars := map[string]string{}
			_ = json.Unmarshal([]byte(cluster.Plan.Vars), &planVars)
			role := getHostRole(newHost.Name)
			workerConfig, err := c.vmConfigRepo.Get(planVars[fmt.Sprintf("%sModel", role)])
			if err != nil {
				return nil, err
			}
			newHost.CpuCore = workerConfig.Cpu
			newHost.Memory = workerConfig.Memory * 1024
		}
		newHosts = append(newHosts, newHost)
	}
	group := allocateZone(cluster.Plan.Zones, newHosts)
	for k, v := range group {
		providerVars := map[string]interface{}{}
		providerVars["provider"] = cluster.Plan.Region.Provider
		_ = json.Unmarshal([]byte(cluster.Plan.Region.Vars), &providerVars)
		cloudClient := cloud_provider.NewCloudClient(providerVars)
		err := allocateIpAddr(cloudClient, *k, v, cluster.ID)
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

func (c *clusterNodeService) doNodeDelete(cluster *model.Cluster, nodes []*model.ClusterNode) error {
	logId, writer, err := ansible.CreateAnsibleLogWriter(cluster.Name)
	if err != nil {
		log.Error(err)
	}
	cluster.LogId = logId
	db.DB.Save(cluster)
	cluster.Nodes, _ = c.NodeRepo.List(cluster.Name)
	inventory := cluster.ParseInventory()
	for i := range inventory.Groups {
		if inventory.Groups[i].Name == "del-worker" {
			for _, n := range nodes {
				inventory.Groups[i].Hosts = append(inventory.Groups[i].Hosts, n.Name)
			}
		}
	}
	k := kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})
	for i := range facts.DefaultFacts {
		k.SetVar(i, facts.DefaultFacts[i])
	}
	clusterVars := cluster.GetKobeVars()
	for j, v := range clusterVars {
		k.SetVar(j, v)
	}
	k.SetVar(facts.ClusterNameFactName, cluster.Name)
	val, _ := c.systemSettingRepo.Get("ip")
	k.SetVar(facts.LocalHostnameFactName, val.Value)
	err = phases.RunPlaybookAndGetResult(k, deleteWorkerPlaybook, writer)
	if err != nil {
		return err
	}
	return nil
}

const addWorkerPlaybook = "91-add-worker.yml"

func (c clusterNodeService) doNodeCreate(cluster *model.Cluster, nodes []*model.ClusterNode) error {
	logId, writer, err := ansible.CreateAnsibleLogWriter(cluster.Name)
	if err != nil {
		log.Error(err)
	}
	cluster.LogId = logId
	db.DB.Save(cluster)
	cluster.Nodes, _ = c.NodeRepo.List(cluster.Name)
	inventory := cluster.ParseInventory()
	for i := range inventory.Groups {
		if inventory.Groups[i].Name == "new-worker" {
			for _, n := range nodes {
				inventory.Groups[i].Hosts = append(inventory.Groups[i].Hosts, n.Name)
			}
		}
	}
	k := kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})
	for i := range facts.DefaultFacts {
		k.SetVar(i, facts.DefaultFacts[i])
	}
	clusterVars := cluster.GetKobeVars()
	for j, v := range clusterVars {
		k.SetVar(j, v)
	}
	k.SetVar(facts.ClusterNameFactName, cluster.Name)
	val, _ := c.systemSettingRepo.Get("ip")
	k.SetVar(facts.LocalHostnameFactName, val.Value)
	err = phases.RunPlaybookAndGetResult(k, addWorkerPlaybook, writer)
	if err != nil {
		return err
	}
	return nil
}
