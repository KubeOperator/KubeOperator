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
	"sync"
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
	if err := db.DB.Where(model.ClusterNode{ClusterID: cluster.ID}).Preload("Host").Preload("Host.Credential").Preload("Host.Zone").Find(&currentNodes).Error; err != nil {
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

func (c clusterNodeService) batchDelete(cluster *model.Cluster, currentNodes []model.ClusterNode, item dto.NodeBatch) error {
	var nodesForDelete []model.ClusterNode
	if err := db.DB.Model(model.ClusterNode{}).Where("name in (?)", item.Nodes).Preload("Host").Preload("Host.Credential").Preload("Host.Zone").Find(&nodesForDelete).Error; err != nil {
		return fmt.Errorf("can not find nodes reason %s", err.Error())
	}
	var notDirtyNodes []model.ClusterNode
	tx := db.DB.Begin()
	for i := range nodesForDelete {
		if nodesForDelete[i].Dirty {
			if err := tx.Model(model.Host{}).Where(model.Host{ID: nodesForDelete[i].HostID}).Updates(map[string]interface{}{
				"ClusterID": "",
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Delete(&nodesForDelete[i]).Error; err != nil {
				tx.Rollback()
				return err
			}

		} else {
			notDirtyNodes = append(notDirtyNodes, nodesForDelete[i])
		}
	}
	tx.Commit()

	needDeleteNodeIds := func() []string {
		var names []string
		for i := range notDirtyNodes {
			names = append(names, notDirtyNodes[i].ID)
		}
		return names
	}()

	if cluster.Spec.Provider == constant.ClusterProviderPlan {
		var p model.Plan
		if err := db.DB.Where(model.Plan{ID: cluster.PlanID}).First(&p).Error; err != nil {
			return fmt.Errorf("can not load plan err %s", err.Error())
		}
		planDTO, err := c.planService.Get(p.Name)
		if err != nil {
			return fmt.Errorf("can not load plan err %s", err.Error())
		}
		cluster.Plan = planDTO.Plan
	}
	go func() {
		if err := db.DB.Model(model.ClusterNode{}).Where("id in (?)", needDeleteNodeIds).Updates(map[string]interface{}{
			"Status": constant.StatusTerminating,
		}).Error; err != nil {
			log.Errorf("can not update node status %s", err.Error())
		}
		if err := c.runDeleteWorkerPlaybook(cluster, notDirtyNodes); err != nil {
			log.Errorf("delete node failed error %s", err.Error())
		}
		if cluster.Spec.Provider == constant.ClusterProviderPlan {
			if err := c.destroyHosts(cluster, currentNodes, notDirtyNodes); err != nil {
				log.Error(err)
			}
		}
		for i := range notDirtyNodes {
			if err := db.DB.Model(model.Host{}).Where(model.Host{ID: notDirtyNodes[i].HostID}).Updates(map[string]interface{}{
				"ClusterID": "",
			}).Error; err != nil {
				log.Error("can not update host clusterID %s reason %s", notDirtyNodes[i].Host.Name, err.Error())
			}
			if err := db.DB.Delete(&notDirtyNodes[i]).Error; err != nil {
				log.Error("can not delete node %s reason %s", notDirtyNodes[i].Name, err.Error())
			}

			if cluster.Spec.Provider == constant.ClusterProviderPlan {
				var projectResource model.ProjectResource
				if err := db.DB.Where(model.ProjectResource{ResourceType: constant.ResourceHost, ResourceID: notDirtyNodes[i].HostID}).First(&projectResource).Error; err != nil {
					log.Errorf("can not find project resource reason %s", err.Error())
				}
				if err := db.DB.Delete(&projectResource).Error; err != nil {
					log.Error("can not delete project resource reason %s", err.Error())
				}
				notDirtyNodes[i].Host.ClusterID = ""
				if err := db.DB.Delete(&notDirtyNodes[i].Host).Error; err != nil {
					log.Error("can not delete host %s reason %s", notDirtyNodes[i].Host.Name, err.Error())
				}
			}
		}
	}()
	return nil
}

func (c *clusterNodeService) destroyHosts(cluster *model.Cluster, currentNodes []model.ClusterNode, nodes []model.ClusterNode) error {
	var aliveHosts []*model.Host
exit:
	for i := range currentNodes {
		for k := range nodes {
			if cluster.Nodes[i].Name == nodes[k].Name {
				continue exit
			}
		}
		aliveHosts = append(aliveHosts, &currentNodes[i].Host)
	}
	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	return doInit(k, cluster.Plan, aliveHosts)
}

func (c clusterNodeService) batchCreate(cluster *model.Cluster, currentNodes []model.ClusterNode, item dto.NodeBatch) error {
	var newNodes []model.ClusterNode
	switch cluster.Spec.Provider {
	case constant.ClusterProviderBareMetal:
		var hosts []model.Host
		for _, hostName := range item.Hosts {
			h, err := c.hostService.Get(hostName)
			if err != nil {
				return fmt.Errorf("can not find host %s reason %s", hostName, err.Error())
			}
			hosts = append(hosts, h.Host)
		}
		ns, err := c.createNodeModels(cluster, currentNodes, hosts)
		if err != nil {
			return err
		}
		newNodes = ns
	case constant.ClusterProviderPlan:
		var p model.Plan
		if err := db.DB.Where(model.Plan{ID: cluster.PlanID}).First(&p).Error; err != nil {
			return fmt.Errorf("can not load plan err %s", err.Error())
		}
		planDTO, err := c.planService.Get(p.Name)
		if err != nil {
			return fmt.Errorf("can not load plan err %s", err.Error())
		}
		cluster.Plan = planDTO.Plan
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
	go func() {
		newNodeIds := func() []string {
			var names []string
			for _, n := range newNodes {
				names = append(names, n.ID)
			}
			return names
		}()
		if cluster.Spec.Provider == constant.ClusterProviderPlan {
			if err := db.DB.Model(model.ClusterNode{}).Where("id in (?)", newNodeIds).Updates(map[string]interface{}{"Status": constant.StatusCreating}).Error; err != nil {
				log.Errorf("can not update node status reason %s", err.Error())
				return
			}
			var allNodes []model.ClusterNode
			if err := db.DB.Where(model.ClusterNode{ClusterID: cluster.ID}).
				Preload("Host").
				Preload("Host.Credential").
				Preload("Host.Zone").Find(&allNodes).Error; err != nil {
				log.Errorf("can not load all nodes %s", err.Error())
				return
			}
			allHosts := func() []*model.Host {
				var hs []*model.Host
				for i := range allNodes {
					hs = append(hs, &allNodes[i].Host)
				}
				return hs
			}()
			if err := c.doCreateHosts(cluster, allHosts); err != nil {
				for i := range newNodes {
					if err := db.DB.Delete(&newNodes[i].Host).Error; err != nil {
						log.Errorf("can not delete host %s reason %s", newNodes[i].Host.Name, err.Error())
					}
				}
				if err := db.DB.Model(model.ClusterNode{}).Where("id in (?)", newNodes).
					Updates(map[string]interface{}{
						"Status":  constant.StatusFailed,
						"Message": fmt.Errorf("can not create hosts reason %s", err.Error()),
						"HostID":  "",
					}).Error; err != nil {
					log.Errorf("can not update node status reason %s", err.Error())
				}
				return
			}
			wg := sync.WaitGroup{}
			for _, h := range allHosts {
				wg.Add(1)
				go func(ho *model.Host) {
					_, err := c.hostService.Sync(ho.Name)
					if err != nil {
						log.Error("sync host %s status error %s", ho.Name, err.Error())
					}
					defer wg.Done()
				}(h)
			}
			wg.Wait()
		}
		if err := db.DB.Model(model.ClusterNode{}).Where("id in (?)", newNodeIds).Updates(map[string]interface{}{"Status": constant.StatusInitializing}).Error; err != nil {
			log.Errorf("can not update node status reason %s", err.Error())
			return
		}
		if err := c.runAddWorkerPlaybook(cluster, newNodes); err != nil {
			if err := db.DB.Model(model.ClusterNode{}).Where("id in (?)", newNodeIds).Updates(map[string]interface{}{"Status": constant.StatusFailed, "Message": err.Error()}).Error; err != nil {
				log.Errorf("can not update node status reason %s", err.Error())
			}
			return
		}
		if err := db.DB.Model(model.ClusterNode{}).Where("id in (?)", newNodeIds).Updates(map[string]interface{}{"Status": constant.StatusRunning}).Error; err != nil {
			log.Errorf("can not update node status reason %s", err.Error())
		}
	}()
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
	if err := db.DB.Where(model.ProjectResource{ResourceID: cluster.ID, ResourceType: constant.ResourceCluster}).First(&clusterResource).Error; err != nil {
		return nil, fmt.Errorf("can not find project resource %s", err.Error())
	}

	tx := db.DB.Begin()
	for i := range newHosts {
		if err := tx.Create(newHosts[i]).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not save host %s reasone %s", newHosts[i].Name, err.Error())
		}
		var ip model.Ip
		tx.Where(model.Ip{Address: newHosts[i].Ip}).First(&ip)
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
	err = phases.RunPlaybookAndGetResult(k, deleteWorkerPlaybook, writer)
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
	val, _ := c.systemSettingRepo.Get("ip")
	k.SetVar(facts.RegistryHostnameFactName, val.Value)
	err = phases.RunPlaybookAndGetResult(k, addWorkerPlaybook, writer)
	if err != nil {
		return err
	}
	return nil
}
