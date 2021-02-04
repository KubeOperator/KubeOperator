package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

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
	Batch(clusterName string, batch dto.NodeBatch) error
	Page(num, size int, clusterName string) (*dto.NodePage, error)
}

var log = logger.Default

func NewClusterNodeService() ClusterNodeService {
	return &clusterNodeService{
		ClusterService:      NewClusterService(),
		NodeRepo:            repository.NewClusterNodeRepository(),
		HostRepo:            repository.NewHostRepository(),
		systemSettingRepo:   repository.NewSystemSettingRepository(),
		projectResourceRepo: repository.NewProjectResourceRepository(),
		messageService:      NewMessageService(),
		vmConfigRepo:        repository.NewVmConfigRepository(),
		hostService:         NewHostService(),
		planService:         NewPlanService(),
	}
}

type clusterNodeService struct {
	ClusterService      ClusterService
	NodeRepo            repository.ClusterNodeRepository
	HostRepo            repository.HostRepository
	planService         PlanService
	systemSettingRepo   repository.SystemSettingRepository
	projectResourceRepo repository.ProjectResourceRepository
	messageService      MessageService
	vmConfigRepo        repository.VmConfigRepository
	hostService         HostService
}

func (c *clusterNodeService) Get(clusterName, name string) (*dto.Node, error) {
	var n model.ClusterNode
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}

	err = db.DB.Where(&model.ClusterNode{ClusterID: cluster.ID, Name: name}).Find(&n).Error
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
				if cluster.Source == constant.ClusterSourceExternal {
					for _, addr := range kn.Status.Addresses {
						if addr.Type == "InternalIP" {
							n.Ip = addr.Address
						}
					}
				}
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
	mNodes, err := c.NodeRepo.List(cluster.Name)
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
				if cluster.Source == constant.ClusterSourceExternal {
					for _, addr := range kn.Status.Addresses {
						if addr.Type == "InternalIP" {
							n.Ip = addr.Address
						}
					}
				}
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
	return nodes, nil
}

func (c *clusterNodeService) Batch(clusterName string, item dto.NodeBatch) error {
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return fmt.Errorf("can not found %s", clusterName)
	}
	var currentNodes []model.ClusterNode
	if err := db.DB.Where(&model.ClusterNode{ClusterID: cluster.ID}).Preload("Host").Preload("Host.Credential").Preload("Host.Zone").Find(&currentNodes).Error; err != nil {
		return fmt.Errorf("can not read cluster %s current nodes %s", cluster.Name, err.Error())
	}
	for _, node := range currentNodes {
		if !node.Dirty && (node.Status == constant.StatusCreating || node.Status == constant.StatusInitializing || node.Status == constant.StatusWaiting) {
			return errors.New("NODE_ALREADY_RUNNING_TASK")
		}
	}
	switch item.Operation {
	case constant.BatchOperationCreate:
		return c.batchCreate(&cluster.Cluster, currentNodes, item)
	case constant.BatchOperationDelete:
		return c.batchDelete(&cluster.Cluster, currentNodes, item)
	}
	return nil
}

// 脏节点只删除数据库数据，正常节点集群中删除节点然后删数据库
func (c clusterNodeService) batchDelete(cluster *model.Cluster, currentNodes []model.ClusterNode, item dto.NodeBatch) error {
	var (
		nodesForDelete []model.ClusterNode
		notDirtyNodes  []model.ClusterNode
		nodeIDs        []string
		hostIDs        []string
	)
	if err := db.DB.Model(&model.ClusterNode{}).Where("name in (?)", item.Nodes).
		Preload("Host").
		Preload("Host.Credential").
		Preload("Host.Zone").
		Find(&nodesForDelete).Error; err != nil {
		return fmt.Errorf("can not find nodes reason %s", err.Error())
	}

	log.Infof("start delete nodes")
	for _, node := range nodesForDelete {
		hostIDs = append(hostIDs, node.Host.ID)
		nodeIDs = append(nodeIDs, node.ID)
		if !node.Dirty {
			notDirtyNodes = append(notDirtyNodes, node)
		}
	}
	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", nodeIDs).
		Updates(map[string]interface{}{"Status": constant.StatusTerminating}).Error; err != nil {
		log.Errorf("can not update node status %s", err.Error())
		return err
	}

	go c.removeNodes(cluster, currentNodes, notDirtyNodes, hostIDs, nodeIDs)
	return nil
}

func (c *clusterNodeService) removeNodes(cluster *model.Cluster, currentNodes, notDirtyNodes []model.ClusterNode, hostIDs, nodeIDs []string) {
	tx := db.DB.Begin()
	if cluster.Spec.Provider == constant.ClusterProviderPlan {
		var p model.Plan
		if err := tx.Where(&model.Plan{ID: cluster.PlanID}).First(&p).Error; err != nil {
			c.updateNodeStatus(nodeIDs, err.Error(), false)
			log.Errorf("can not load plan err %s", err.Error())
		}
		planDTO, err := c.planService.Get(p.Name)
		if err != nil {
			c.updateNodeStatus(nodeIDs, err.Error(), false)
			log.Errorf("can not load plan err %s", err.Error())
		}
		cluster.Plan = planDTO.Plan

		if err := c.runDeleteWorkerPlaybook(cluster, notDirtyNodes); err != nil {
			c.updateNodeStatus(nodeIDs, err.Error(), false)
			log.Errorf("delete node failed error %s", err.Error())
		}
		if err := c.destroyHosts(cluster, currentNodes, notDirtyNodes); err != nil {
			log.Error(err)
		}
		log.Info("delete all nodes successful! now start updata cluster datas")

		if err := tx.Model(&model.Host{}).Where("id in (?)", hostIDs).
			Delete(&model.Host{}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(nodeIDs, err.Error(), true)
			log.Errorf("can not update hosts clusterID reason %s", err.Error())
		}
		if err := tx.Model(&model.ProjectResource{}).Where("resource_id in (?) AND resource_type = ?", hostIDs, constant.ResourceHost).
			Delete(&model.ProjectResource{}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(nodeIDs, err.Error(), true)
			log.Errorf("can not delete project resource reason %s", err.Error())
		}
	} else {
		if err := c.runDeleteWorkerPlaybook(cluster, notDirtyNodes); err != nil {
			c.updateNodeStatus(nodeIDs, err.Error(), false)
			log.Errorf("delete node failed error %s", err.Error())
		}
		log.Info("delete all nodes successful! now start updata cluster datas")

		if err := tx.Model(&model.Host{}).Where("id in (?)", hostIDs).
			Updates(map[string]interface{}{"ClusterID": ""}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(nodeIDs, err.Error(), true)
			log.Errorf("can not update hosts clusterID reason %s", err.Error())
		}
	}
	if err := tx.Model(&model.ClusterNode{}).Where("id in (?)", nodeIDs).Delete(&model.ClusterNode{}).Error; err != nil {
		tx.Rollback()
		log.Errorf("can not delete nodes reason %s", err.Error())
	}
	tx.Commit()
	log.Info("delete node successful!")
}

func (c *clusterNodeService) updateNodeStatus(notDirtyNodeIDs []string, errMsg string, isDirty bool) {
	log.Errorf("delete node failed，cluster data has been rollback")
	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", notDirtyNodeIDs).
		Updates(map[string]interface{}{"Status": constant.ClusterFailed, "Message": errMsg, "Dirty": isDirty}).Error; err != nil {
		log.Errorf("can not update node status %s", err.Error())
	}
}

func (c *clusterNodeService) destroyHosts(cluster *model.Cluster, currentNodes []model.ClusterNode, deleteNodes []model.ClusterNode) error {
	var aliveHosts []*model.Host
	for i := range currentNodes {
		alive := true
		for k := range deleteNodes {
			if currentNodes[i].Name == deleteNodes[k].Name {
				alive = false
			}
		}
		if alive {
			aliveHosts = append(aliveHosts, &currentNodes[i].Host)
		}
	}
	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	return doInit(k, cluster.Plan, aliveHosts)
}

func (c clusterNodeService) batchCreate(cluster *model.Cluster, currentNodes []model.ClusterNode, item dto.NodeBatch) error {
	var (
		newNodes  []model.ClusterNode
		hostNames []string
	)
	for _, host := range item.Hosts {
		hostNames = append(hostNames, host)
	}

	log.Info("start create cluster nodes")
	switch cluster.Spec.Provider {
	case constant.ClusterProviderBareMetal:
		var hosts []model.Host
		if err := db.DB.Model(&model.Host{}).Where("name in (?)", hostNames).
			Preload("Volumes").
			Preload("Credential").
			Find(&hosts).Error; err != nil {
			return fmt.Errorf("can not find hosts reason %s", err.Error())
		}
		ns, err := c.createNodeModels(cluster, currentNodes, hosts)
		if err != nil {
			return err
		}
		newNodes = ns
	case constant.ClusterProviderPlan:
		var plan model.Plan
		if err := db.DB.Where(&model.Plan{ID: cluster.PlanID}).First(&plan).
			Preload("Zones").
			Preload("Region").Find(&plan).Error; err != nil {
			return fmt.Errorf("can not load plan err %s", err.Error())
		}
		cluster.Plan = plan
		hosts, err := c.createHostModels(cluster, item.Increase)
		if err != nil {
			return fmt.Errorf("can not create host models err %s", err.Error())
		}
		ns, err := c.createNodeModels(cluster, currentNodes, hosts)
		if err != nil {
			return err
		}
		newNodes = ns
	}
	go c.addNodes(cluster, newNodes)
	return nil
}

func (c clusterNodeService) addNodes(cluster *model.Cluster, newNodes []model.ClusterNode) {
	var (
		newNodeIDs []string
		newHostIDs []string
	)
	for _, n := range newNodes {
		newNodeIDs = append(newNodeIDs, n.ID)
		newHostIDs = append(newHostIDs, n.Host.ID)
	}

	if cluster.Spec.Provider == constant.ClusterProviderPlan {
		log.Info("cluster-plan start add hosts, update hosts status and infos")
		c.updataHostInfo(cluster, newNodeIDs, newHostIDs)
	}
	log.Info("start binding nodes to cluster")
	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", newNodeIDs).
		Updates(map[string]interface{}{"Status": constant.StatusInitializing}).Error; err != nil {
		log.Errorf("can not update node status reason %s", err.Error())
		return
	}
	if err := c.runAddWorkerPlaybook(cluster, newNodes); err != nil {
		if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", newNodeIDs).
			Updates(map[string]interface{}{"Status": constant.StatusFailed, "Message": err.Error()}).Error; err != nil {
			log.Errorf("can not update node status reason %s", err.Error())
		}
		return
	}
	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", newNodeIDs).
		Updates(map[string]interface{}{"Status": constant.StatusRunning}).Error; err != nil {
		log.Errorf("can not update node status reason %s", err.Error())
	}
	log.Info("create cluster nodes successful!")
}

// 添加主机、修改主机状态及相关信息
func (c clusterNodeService) updataHostInfo(cluster *model.Cluster, newNodeIDs, newHostIDs []string) error {
	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", newNodeIDs).
		Updates(map[string]interface{}{"Status": constant.StatusCreating}).Error; err != nil {
		log.Errorf("can not update node status reason %s", err.Error())
		return err
	}
	var allNodes []model.ClusterNode
	if err := db.DB.Where(&model.ClusterNode{ClusterID: cluster.ID}).
		Preload("Host").
		Preload("Host.Credential").
		Preload("Host.Zone").Find(&allNodes).Error; err != nil {
		log.Errorf("can not load all nodes %s", err.Error())
		return err
	}

	var allHosts []*model.Host
	for i := range allNodes {
		allHosts = append(allHosts, &allNodes[i].Host)
	}

	if err := c.doCreateHosts(cluster, allHosts); err != nil {
		if err := db.DB.Model(&model.Host{}).Where("id in (?)", newHostIDs).Delete(&model.Host{}).Error; err != nil {
			log.Errorf("can not delete hosts reason %s", err.Error())
		}
		if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", newNodeIDs).
			Updates(map[string]interface{}{
				"Status":  constant.StatusFailed,
				"Message": fmt.Errorf("can not create hosts reason %s", err.Error()),
				"HostID":  "",
			}).Error; err != nil {
			log.Errorf("can not update node status reason %s", err.Error())
		}
		return err
	}
	wg := sync.WaitGroup{}
	for _, h := range allHosts {
		wg.Add(1)
		go func(ho *model.Host) {
			_, err := c.hostService.Sync(ho.Name)
			if err != nil {
				log.Errorf("sync host %s status error %s", ho.Name, err.Error())
			}
			defer wg.Done()
		}(h)
	}
	wg.Wait()
	return nil
}

func (c clusterNodeService) createNodeModels(cluster *model.Cluster, currentNodes []model.ClusterNode, hosts []model.Host) ([]model.ClusterNode, error) {
	var newNodes []model.ClusterNode
	hash := map[string]interface{}{}
	for _, n := range currentNodes {
		hash[n.Name] = nil
	}
	for _, host := range hosts {
		var name string
		for i := 1; i < len(currentNodes)+len(hosts); i++ {
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
			Host:      host,
		}
		newNodes = append(newNodes, n)
	}
	tx := db.DB.Begin()
	for i := range newNodes {
		newNodes[i].Host.ClusterID = cluster.ID
		if err := tx.Save(&newNodes[i].Host).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not save host %s", newNodes[i].Host.Name)
		}
		if err := tx.Create(&newNodes[i]).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not save node %s", newNodes[i].Name)
		}
	}
	tx.Commit()
	return newNodes, nil
}

func (c clusterNodeService) createHostModels(cluster *model.Cluster, increase int) ([]model.Host, error) {
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
		providerVars["datacenter"] = cluster.Plan.Region.Datacenter
		zoneVars := map[string]interface{}{}
		_ = json.Unmarshal([]byte(k.Vars), &zoneVars)
		providerVars["cluster"] = zoneVars["cluster"]
		_ = json.Unmarshal([]byte(cluster.Plan.Region.Vars), &providerVars)
		cloudClient := cloud_provider.NewCloudClient(providerVars)
		err := allocateIpAddr(cloudClient, *k, v, cluster.ID)
		if err != nil {
			return nil, err
		}
		err = allocateDatastore(cloudClient, *k, v)
		if err != nil {
			return nil, err
		}
	}

	var clusterResource model.ProjectResource
	if err := db.DB.Where(&model.ProjectResource{ResourceID: cluster.ID, ResourceType: constant.ResourceCluster}).First(&clusterResource).Error; err != nil {
		return nil, fmt.Errorf("can not find project resource %s", err.Error())
	}

	tx := db.DB.Begin()
	for i := range newHosts {
		if err := tx.Create(newHosts[i]).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not save host %s reasone %s", newHosts[i].Name, err.Error())
		}
		var ip model.Ip
		tx.Where(&model.Ip{Address: newHosts[i].Ip}).First(&ip)
		if ip.ID != "" {
			ip.Status = constant.IpUsed
			ip.ClusterID = cluster.ID
			tx.Save(&ip)
		}
		hostProjectResource := model.ProjectResource{
			ResourceType: constant.ResourceHost,
			ResourceID:   newHosts[i].ID,
			ProjectID:    clusterResource.ProjectID,
		}
		if err := tx.Create(&hostProjectResource).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not create peroject resource host %s ", newHosts[i].Name)
		}
	}
	tx.Commit()

	res := func() []model.Host {
		var hs []model.Host
		for i := range newHosts {
			hs = append(hs, *newHosts[i])
		}
		return hs
	}()
	return res, nil
}

func (c clusterNodeService) doCreateHosts(cluster *model.Cluster, hosts []*model.Host) error {
	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	return doInit(k, cluster.Plan, hosts)
}

const deleteWorkerPlaybook = "96-remove-worker.yml"

func (c *clusterNodeService) runDeleteWorkerPlaybook(cluster *model.Cluster, nodes []model.ClusterNode) error {
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
	registryIp, _ := c.systemSettingRepo.Get("ip")
	registryProtocol, _ := c.systemSettingRepo.Get("REGISTRY_PROTOCOL")
	k.SetVar(facts.RegistryProtocolFactName, registryProtocol.Value)
	k.SetVar(facts.RegistryHostnameFactName, registryIp.Value)
	err = phases.RunPlaybookAndGetResult(k, deleteWorkerPlaybook, "", writer)
	if err != nil {
		return err
	}
	return nil
}

const addWorkerPlaybook = "91-add-worker.yml"

func (c *clusterNodeService) runAddWorkerPlaybook(cluster *model.Cluster, nodes []model.ClusterNode) error {
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
	registryIp, _ := c.systemSettingRepo.Get("ip")
	registryProtocol, _ := c.systemSettingRepo.Get("REGISTRY_PROTOCOL")
	k.SetVar(facts.RegistryProtocolFactName, registryProtocol.Value)
	k.SetVar(facts.RegistryHostnameFactName, registryIp.Value)
	err = phases.RunPlaybookAndGetResult(k, addWorkerPlaybook, "", writer)
	if err != nil {
		return err
	}
	return nil
}
