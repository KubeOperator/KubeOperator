package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterNodeService interface {
	Get(clusterName, name string) (*dto.Node, error)
	List(clusterName string) ([]dto.Node, error)
	Batch(clusterName string, batch dto.NodeBatch) error
	Recreate(clusterName string, batch dto.NodeBatch) error
	Page(num, size int, isPolling, clusterName string) (*dto.NodePage, error)
}

func NewClusterNodeService() ClusterNodeService {
	return &clusterNodeService{
		ClusterService:      NewClusterService(),
		clusterRepo:         repository.NewClusterRepository(),
		NodeRepo:            repository.NewClusterNodeRepository(),
		taskLogService:      NewTaskLogService(),
		HostRepo:            repository.NewHostRepository(),
		ntpServerRepo:       repository.NewNtpServerRepository(),
		projectResourceRepo: repository.NewProjectResourceRepository(),
		messageService:      NewMessageService(),
		vmConfigRepo:        repository.NewVmConfigRepository(),
		hostService:         NewHostService(),
		planService:         NewPlanService(),
	}
}

type clusterNodeService struct {
	ClusterService      ClusterService
	clusterRepo         repository.ClusterRepository
	NodeRepo            repository.ClusterNodeRepository
	taskLogService      TaskLogService
	HostRepo            repository.HostRepository
	planService         PlanService
	ntpServerRepo       repository.NtpServerRepository
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

	err = db.DB.Where("cluster_id = ? AND name = ?", cluster.ID, name).Find(&n).Error
	if err != nil {
		return nil, err
	}
	return &dto.Node{
		ClusterNode: n,
	}, nil
}

func (c clusterNodeService) Page(num, size int, isPolling, clusterName string) (*dto.NodePage, error) {
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
	nodesAfterSync := syncNodeStatus(mNodes, kubeNodes, cluster.Source, isPolling)
	return &dto.NodePage{
		Items: nodesAfterSync,
		Total: count,
	}, nil
}

func (c clusterNodeService) List(clusterName string) ([]dto.Node, error) {
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
	nodesAfterSync := syncNodeStatus(mNodes, kubeNodes, cluster.Source, "true")
	return nodesAfterSync, nil
}

func (c *clusterNodeService) Batch(clusterName string, item dto.NodeBatch) error {
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return fmt.Errorf("can not found %s", clusterName)
	}
	var currentNodes []model.ClusterNode
	if err := db.DB.Where("cluster_id = ?", cluster.ID).Preload("Host").Preload("Host.Credential").Preload("Host.Zone").Find(&currentNodes).Error; err != nil {
		return fmt.Errorf("can not read cluster %s current nodes %s", cluster.Name, err.Error())
	}
	for _, node := range currentNodes {
		if !node.Dirty && (node.Status == constant.StatusCreating || node.Status == constant.StatusInitializing || node.Status == constant.StatusWaiting) {
			return errors.New("NODE_ALREADY_RUNNING_TASK")
		}
	}
	switch item.Operation {
	case constant.BatchOperationCreate:
		tasklog := model.TaskLog{
			ClusterID: cluster.ID,
			Type:      constant.TaskLogTypeClusterNodeExtend,
			Phase:     constant.ClusterWaiting,
		}
		if err := c.taskLogService.Start(&tasklog); err != nil {
			return err
		}
		cluster.TaskLog = tasklog
		return c.batchCreate(&cluster.Cluster, currentNodes, item)
	case constant.BatchOperationDelete:
		if err := db.DB.Model(&model.ClusterNode{}).Where("name in (?)", item.Nodes).
			Updates(map[string]interface{}{"Status": constant.StatusTerminating, "PreStatus": constant.StatusFailed, "Message": ""}).Error; err != nil {
			logger.Log.Errorf("can not update node status %s", err.Error())
			return err
		}
		var nodesForDelete []model.ClusterNode
		if err := db.DB.Where("name in (?)", item.Nodes).
			Preload("Host").
			Preload("Host.Credential").
			Preload("Host.Zone").
			Find(&nodesForDelete).Error; err != nil {
			return err
		}
		tasklog := model.TaskLog{
			ClusterID: cluster.ID,
			Type:      constant.TaskLogTypeClusterNodeShrink,
			PrePhase:  constant.StatusFailed,
			Phase:     constant.StatusTerminating,
		}
		if err := c.taskLogService.Start(&tasklog); err != nil {
			return err
		}
		cluster.TaskLog = tasklog

		go c.removeNodes(&cluster.Cluster, item, currentNodes, nodesForDelete)
		return nil
	}
	return nil
}

func prepareRemove(item dto.NodeBatch, nodesForDelete []model.ClusterNode) ([]model.ClusterNode, []string, []string, []string, []string, error) {
	var (
		notDirtyNodes []model.ClusterNode
		dirtyNodeIDs  []string
		nodeIDs       []string
		hostIDs       []string
		hostIPs       []string
	)

	// notDirtyNodes 待执行脚本的节点（非强制删除的所有、强制删除时 运行状态或非脏节点）
	logger.Log.Infof("start delete nodes")
	for _, node := range nodesForDelete {

		hostIDs = append(hostIDs, node.Host.ID)
		hostIPs = append(hostIPs, node.Host.Ip)
		nodeIDs = append(nodeIDs, node.ID)
		if item.IsForce {
			if !node.Dirty || node.Status == constant.StatusRunning {
				notDirtyNodes = append(notDirtyNodes, node)
			} else {
				dirtyNodeIDs = append(dirtyNodeIDs, node.ID)
			}
		} else {
			notDirtyNodes = append(notDirtyNodes, node)
		}
	}
	return notDirtyNodes, dirtyNodeIDs, nodeIDs, hostIDs, hostIPs, nil
}

func (c *clusterNodeService) removeNodes(cluster *model.Cluster, item dto.NodeBatch, currentNodes, nodesForDelete []model.ClusterNode) {
	notDirtyNodes, dirtyNodeIDs, nodeIDs, hostIDs, hostIPs, err := prepareRemove(item, nodesForDelete)
	if err != nil {
		c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, false)
	}

	if cluster.Provider == constant.ClusterProviderPlan {
		var p model.Plan
		if err := db.DB.Where("id = ?", cluster.PlanID).First(&p).Error; err != nil {
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, false)
			return
		}
		planDTO, err := c.planService.Get(p.Name)
		if err != nil {
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, false)
			return
		}
		cluster.Plan = planDTO.Plan

		if err := c.runDeleteWorkerPlaybook(cluster, nodesForDelete, removeWorkerPlaybook); err != nil {
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, true)
			return
		}
		if err := c.destroyHosts(cluster, currentNodes, nodeIDs); err != nil {
			if !item.IsForce {
				c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, true)
				return
			} else {
				logger.Log.Errorf("destroy host failed, err: %s", err.Error())
			}
		}
		logger.Log.Info("delete all nodes successful! now start updata cluster datas")

		tx := db.DB.Begin()
		if err := tx.Where("id in (?)", hostIDs).Delete(&model.Host{}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, false)
			return
		}
		if err := tx.Where("resource_id in (?) AND resource_type = ?", hostIDs, constant.ResourceHost).
			Delete(&model.ProjectResource{}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, false)
			return
		}
		if err := tx.Where("resource_id in (?) AND resource_type = ?", hostIDs, constant.ResourceHost).
			Delete(&model.ClusterResource{}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, false)
			return
		}
		if err := tx.Model(&model.Ip{}).Where("address in (?)", hostIPs).
			Update("status", constant.IpAvailable).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, false)
			return
		}
		tx.Commit()
	} else {
		if err := c.runDeleteWorkerPlaybook(cluster, nodesForDelete, removeWorkerPlaybook); err != nil {
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, true)
			return
		}
		tx := db.DB.Begin()
		if len(notDirtyNodes) != 0 {
			if err := c.runDeleteWorkerPlaybook(cluster, notDirtyNodes, resetWorkerPlaybook); err != nil {
				// 未执行 reset 的脏节点直接删除
				if err := tx.Where("id in (?)", dirtyNodeIDs).Delete(&model.ClusterNode{}).Error; err != nil {
					logger.Log.Errorf("delete node failed, err: %v", err)
				}
				// 执行 reset 失败的节点，返回错误信息
				var notDirtyNodeIDs []string
				for _, node := range notDirtyNodes {
					notDirtyNodeIDs = append(notDirtyNodeIDs, node.ID)
				}
				tx.Rollback()
				c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, notDirtyNodeIDs, err, true)
				return
			}
		}
		logger.Log.Info("delete all nodes successful! now start updata cluster datas")
		if err := tx.Model(&model.Host{}).Where("id in (?)", hostIDs).Update(map[string]interface{}{"ClusterID": ""}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, false)
			return
		}
		if err := tx.Where("resource_id in (?) AND resource_type = ?", hostIDs, constant.ResourceHost).
			Delete(&model.ClusterResource{}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, false)
			return
		}
		tx.Commit()
	}
	if err := db.DB.Where("id in (?)", nodeIDs).Delete(&model.ClusterNode{}).Error; err != nil {
		c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, constant.StatusTerminating, nodeIDs, err, false)
		return
	}
	_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterRemoveWorker, true, ""), cluster.Name, constant.ClusterRemoveWorker)
	logger.Log.Info("delete node successful!")
}

func (c *clusterNodeService) destroyHosts(cluster *model.Cluster, currentNodes []model.ClusterNode, deleteNodeIDs []string) error {
	var aliveHosts []*model.Host
	for i := range currentNodes {
		alive := true
		for k := range deleteNodeIDs {
			if currentNodes[i].ID == deleteNodeIDs[k] {
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
	hostNames = append(hostNames, item.Hosts...)

	logger.Log.Info("start create cluster nodes")
	switch cluster.Provider {
	case constant.ClusterProviderBareMetal:
		var hosts []model.Host
		if err := db.DB.Where("name in (?)", hostNames).
			Preload("Volumes").
			Preload("Credential").
			Find(&hosts).Error; err != nil {
			return fmt.Errorf("get hosts failed: %v", err)
		}
		ns, err := c.createNodeModels(cluster, currentNodes, hosts)
		if err != nil {
			return fmt.Errorf("create node model failed: %v", err)
		}
		newNodes = ns
	case constant.ClusterProviderPlan:
		var plan model.Plan
		if err := db.DB.Where("id = ?", cluster.PlanID).First(&plan).
			Preload("Zones").
			Preload("Region").Find(&plan).Error; err != nil {
			return fmt.Errorf("load plan failed: %v", err)
		}
		cluster.Plan = plan
		hosts, err := c.createHostModels(cluster, item.Increase)
		if err != nil {
			return fmt.Errorf("create host model failed: %v", err)
		}
		ns, err := c.createNodeModels(cluster, currentNodes, hosts)
		if err != nil {
			return fmt.Errorf("create node model failed: %v", err)
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

	if cluster.Provider == constant.ClusterProviderPlan {
		logger.Log.Info("cluster-plan start add hosts, update hosts status and infos")
		if err := c.updataHostInfo(cluster, newNodeIDs, newHostIDs); err != nil {
			c.updateNodeStatus(cluster, constant.ClusterAddWorker, constant.StatusFailed, constant.StatusCreating, newNodeIDs, err, false)
			return
		}
	}

	if err := c.AddWorkInit(cluster.Name, newNodes); err != nil {
		c.updateNodeStatus(cluster, constant.ClusterAddWorker, constant.StatusFailed, constant.StatusInitializing, newNodeIDs, err, false)
		return
	}
}

func (c *clusterNodeService) updateNodeStatus(cluster *model.Cluster, operation, status, preStatus string, nodeIDs []string, errMsg error, isDirty bool) {
	errmsg := ""
	if errMsg != nil {
		errmsg = errMsg.Error()
	}
	if status == constant.ClusterFailed {
		_ = c.messageService.SendMessage(constant.System, false, GetContent(operation, false, errmsg), cluster.Name, operation)
	}
	logger.Log.Infof("change node statu %s to %s, msg: %v", preStatus, status, errMsg)
	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", nodeIDs).
		Updates(map[string]interface{}{"Status": status, "PreStatus": preStatus, "Message": errmsg, "Dirty": isDirty}).Error; err != nil {
		logger.Log.Errorf("can not update node status %s", err.Error())
	}
	if err := db.DB.Model(&model.TaskLog{}).Where("id = ?", cluster.TaskLog.ID).
		Updates(map[string]interface{}{"Phase": status, "PrePhase": preStatus, "Message": errmsg}).Error; err != nil {
		logger.Log.Errorf("can not update task log status %s", err.Error())
	}
}

// 添加主机、修改主机状态及相关信息
func (c clusterNodeService) updataHostInfo(cluster *model.Cluster, newNodeIDs, newHostIDs []string) error {
	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", newNodeIDs).
		Updates(map[string]interface{}{"Status": constant.StatusCreating}).Error; err != nil {
		logger.Log.Errorf("can not update node status reason %s", err.Error())
		return err
	}
	var allNodes []model.ClusterNode
	if err := db.DB.Where("cluster_id = ?", cluster.ID).
		Preload("Host").
		Preload("Host.Credential").
		Preload("Host.Zone").Find(&allNodes).Error; err != nil {
		logger.Log.Errorf("can not load all nodes %s", err.Error())
		return err
	}

	var allHosts []*model.Host
	for i := range allNodes {
		allHosts = append(allHosts, &allNodes[i].Host)
	}

	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	if err := doInit(k, cluster.Plan, allHosts); err != nil {
		if err := db.DB.Where("id in (?)", newHostIDs).Delete(&model.Host{}).Error; err != nil {
			logger.Log.Errorf("can not delete hosts reason %s", err.Error())
		}
		if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", newNodeIDs).
			Updates(map[string]interface{}{
				"Status":    constant.StatusFailed,
				"PreStatus": constant.StatusCreating,
				"Message":   fmt.Errorf("can not create hosts reason %s", err.Error()),
				"HostID":    "",
			}).Error; err != nil {
			logger.Log.Errorf("can not update node status reason %s", err.Error())
		}
		return err
	}
	wg := sync.WaitGroup{}
	for _, h := range allHosts {
		wg.Add(1)
		go func(ho *model.Host) {
			_, err := c.hostService.Sync(ho.Name)
			if err != nil {
				logger.Log.Errorf("sync host %s status error %s", ho.Name, err.Error())
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
		switch cluster.NodeNameRule {
		case constant.NodeNameRuleDefault:
			for i := 1; i < len(currentNodes)+len(hosts); i++ {
				name = fmt.Sprintf("%s-%s-%d", cluster.Name, constant.NodeRoleNameWorker, i)
				if _, ok := hash[name]; ok {
					continue
				}
				hash[name] = nil
				break
			}
		case constant.NodeNameRuleIP:
			name = host.Ip
			if _, ok := hash[name]; ok {
				continue
			}
			hash[name] = nil
		case constant.NodeNameRuleHostName:
			name = host.Name
			if _, ok := hash[name]; ok {
				continue
			}
			hash[name] = nil
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

	var projectResource model.ProjectResource
	if err := db.DB.Where("resource_id = ? AND resource_type = ?", cluster.ID, constant.ResourceCluster).First(&projectResource).Error; err != nil {
		return nil, fmt.Errorf("can not find project resource %s", err.Error())
	}

	tx := db.DB.Begin()
	for i := range newHosts {
		if err := tx.Create(newHosts[i]).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not save host %s reasone %s", newHosts[i].Name, err.Error())
		}
		var ip model.Ip
		if err := tx.Where("address = ?", newHosts[i].Ip).First(&ip).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not save host %s reasone %s", newHosts[i].Name, err.Error())
		}
		if ip.ID != "" {
			ip.Status = constant.IpUsed
			ip.ClusterID = cluster.ID
			tx.Save(&ip)
		}
		hostProjectResource := model.ProjectResource{
			ResourceType: constant.ResourceHost,
			ResourceID:   newHosts[i].ID,
			ProjectID:    projectResource.ProjectID,
		}
		if err := tx.Create(&hostProjectResource).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not create peroject resource host %s ", newHosts[i].Name)
		}
		clusterResource := model.ClusterResource{
			ResourceType: constant.ResourceHost,
			ResourceID:   newHosts[i].ID,
			ClusterID:    cluster.ID,
		}
		if err := tx.Create(&clusterResource).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not create cluster resource host %s ", newHosts[i].Name)
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

const removeWorkerPlaybook = "96-remove-worker.yml"
const resetWorkerPlaybook = "97-reset-worker.yml"

func (c *clusterNodeService) runDeleteWorkerPlaybook(cluster *model.Cluster, nodes []model.ClusterNode, playbookName string) error {
	detail := model.TaskLogDetail{
		Task:   playbookName,
		TaskID: cluster.TaskLog.ID,
	}
	if err := c.taskLogService.StartDetail(&detail); err != nil {
		return err
	}

	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, cluster.TaskLog.ID)
	if err != nil {
		logger.Log.Error(err)
	}
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
	ntps, _ := c.ntpServerRepo.GetAddressStr()
	k.SetVar(facts.NtpServerName, ntps)
	if err = phases.RunPlaybookAndGetResult(k, playbookName, "", writer); err != nil {
		_ = c.taskLogService.EndDetail(&detail, constant.StatusFailed, err.Error())
		return err
	}
	_ = c.taskLogService.EndDetail(&detail, constant.StatusSuccess, err.Error())
	return nil
}

func (c *clusterNodeService) Recreate(clusterName string, batch dto.NodeBatch) error {
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"TaskLog", "TaskLog.Details", "SpecConf", "SpecNetwork", "SpecRuntime"})
	if err != nil {
		return err
	}

	cluster.TaskLog.Phase = constant.StatusInitializing
	cluster.TaskLog.PrePhase = constant.ClusterFailed

	if len(cluster.TaskLog.Details) > 0 {
		for i := range cluster.TaskLog.Details {
			if cluster.TaskLog.Details[i].Status == constant.ConditionFalse {
				cluster.TaskLog.Details[i].Status = constant.ConditionUnknown
				cluster.TaskLog.Details[i].Message = ""
			}
		}
	}
	if err := c.taskLogService.Save(&cluster.TaskLog); err != nil {
		return err
	}

	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, cluster.TaskLog.ID)
	if err != nil {
		return err
	}

	var nodes []model.ClusterNode
	if err := db.DB.Where("status_id = ?", cluster.TaskLog.ID).Find(&nodes).Error; err != nil {
		return err
	}
	if err := db.DB.Model(&model.ClusterNode{}).Where("status_id = ?", batch.StatusID).
		Updates(map[string]interface{}{"Status": constant.StatusInitializing, "PreStatus": constant.StatusFailed}).Error; err != nil {
		return fmt.Errorf("can not update cluster status %s", err.Error())
	}
	go c.doBindNodeToCluster(&cluster, nodes, writer)
	return nil
}

func (c *clusterNodeService) AddWorkInit(clusterName string, nodes []model.ClusterNode) error {
	logger.Log.Info("start binding nodes to cluster")
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime"})
	if err != nil {
		return err
	}
	var nodeIds []string
	for _, n := range nodes {
		nodeIds = append(nodeIds, n.ID)
	}

	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", nodeIds).
		Updates(map[string]interface{}{"Status": constant.StatusInitializing, "PreStatus": constant.StatusCreating, "status_id": cluster.TaskLog.ID}).Error; err != nil {
		return fmt.Errorf("can not update cluster status %s", err.Error())
	}

	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, cluster.TaskLog.ID)
	if err != nil {
		return err
	}

	go c.doBindNodeToCluster(&cluster, nodes, writer)
	return nil
}

func (c *clusterNodeService) doBindNodeToCluster(cluster *model.Cluster, nodes []model.ClusterNode, writer io.Writer) {
	var nodeIds []string
	k := &adm.AnsibleHelper{
		Writer: writer,
	}
	inventory := cluster.ParseInventory()
	for i := range inventory.Groups {
		if inventory.Groups[i].Name == "new-worker" {
			for _, n := range nodes {
				nodeIds = append(nodeIds, n.ID)
				inventory.Groups[i].Hosts = append(inventory.Groups[i].Hosts, n.Name)
			}
		}
	}
	k.Kobe = kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})
	for i := range facts.DefaultFacts {
		k.Kobe.SetVar(i, facts.DefaultFacts[i])
	}
	clusterVars := cluster.GetKobeVars()
	for j, v := range clusterVars {
		k.Kobe.SetVar(j, v)
	}
	k.Kobe.SetVar(facts.ClusterNameFactName, cluster.Name)
	ntps, _ := c.ntpServerRepo.GetAddressStr()
	k.Kobe.SetVar(facts.NtpServerName, ntps)
	maniFest, _ := adm.GetManiFestBy(cluster.Version)
	if maniFest.Name != "" {
		vars := maniFest.GetVars()
		for j, v := range vars {
			k.Kobe.SetVar(j, v)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())

	statusChan := make(chan adm.AnsibleHelper)
	go c.doHandleResp(ctx, *k, statusChan)
	for {
		result := <-statusChan
		cluster.TaskLog.Phase = result.Status
		cluster.TaskLog.Message = result.Message
		cluster.TaskLog.Details = result.LogDetail
		_ = c.taskLogService.Save(&cluster.TaskLog)
		// 保存进度
		switch result.Status {
		case constant.StatusRunning:
			c.updateNodeStatus(cluster, constant.ClusterAddWorker, constant.ClusterRunning, constant.ClusterInitializing, nodeIds, fmt.Errorf(result.Message), false)
			cancel()
			return
		case constant.StatusFailed:
			c.updateNodeStatus(cluster, constant.ClusterAddWorker, constant.ClusterFailed, constant.ClusterInitializing, nodeIds, fmt.Errorf(result.Message), false)
			cancel()
			return
		}
	}
}

func (c clusterNodeService) doHandleResp(ctx context.Context, aHelper adm.AnsibleHelper, statusChan chan adm.AnsibleHelper) {
	ad := adm.NewClusterAdm()
	for {
		resp, err := ad.OnAddWorker(aHelper)
		if err != nil {
			aHelper.Message = err.Error()
		}
		aHelper.Status = resp.Status
		select {
		case <-ctx.Done():
			return
		case statusChan <- aHelper:
		}
		time.Sleep(5 * time.Second)
	}
}

// db 存在，cluster 不存在  ====>  失联
func syncNodeStatus(nodesInDB []model.ClusterNode, kubeNodes *v1.NodeList, source, isPolling string) []dto.Node {
	var (
		runningList  []string
		notReadyList []string
		lostedList   []string
		failedList   []string
		nodes        []dto.Node
	)
	for _, node := range nodesInDB {
		n := dto.Node{
			ClusterNode: node,
			Ip:          node.Host.Ip,
		}
		hasNode := false
		for _, kn := range kubeNodes.Items {
			if kn.ObjectMeta.Name == node.Name {
				hasNode = true
				if source == constant.ClusterSourceExternal {
					for _, addr := range kn.Status.Addresses {
						if addr.Type == "InternalIP" {
							n.Ip = addr.Address
							break
						}
					}
				}
			}
			if hasNode {
				n.Info = kn
				if isPolling != "true" {
					if node.Status == constant.StatusRunning || node.Status == constant.StatusNotReady || node.Status == constant.StatusLost {
						for _, condition := range kn.Status.Conditions {
							if condition.Type == "Ready" && condition.Status == "True" {
								if node.Status != constant.StatusRunning {
									runningList = append(runningList, node.ID)
								}
								n.Status = constant.StatusRunning
							}
							if condition.Type == "Ready" && (condition.Status == "False" || condition.Status == "Unknown") {
								if node.Status != constant.StatusNotReady {
									notReadyList = append(notReadyList, node.ID)
								}
								n.Status = constant.StatusNotReady
								n.Message = condition.Message
							}
						}
					}
				}
				nodes = append(nodes, n)
				break
			}
		}
		if hasNode {
			continue
		}
		if node.Status == constant.StatusRunning {
			lostedList = append(lostedList, node.ID)
		}
		nodes = append(nodes, n)
	}
	go func() {
		if len(runningList) != 0 {
			if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", runningList).Updates(map[string]interface{}{"status": constant.StatusRunning}).Error; err != nil {
				logger.Log.Errorf("Change node(%v) status into running failed, err is %s", runningList, err.Error())
			}
		}
		if len(failedList) != 0 {
			if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", failedList).Updates(map[string]interface{}{"status": constant.StatusFailed}).Error; err != nil {
				logger.Log.Errorf("Change node(%v) status into failed failed, err is %s", failedList, err.Error())
			}
		}
		if len(notReadyList) != 0 {
			if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", notReadyList).Updates(map[string]interface{}{"status": constant.StatusNotReady}).Error; err != nil {
				logger.Log.Errorf("Change node(%v) status into not ready failed, err is %s", notReadyList, err.Error())
			}
		}
		if len(lostedList) != 0 {
			if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", lostedList).Updates(map[string]interface{}{"status": constant.StatusLost}).Error; err != nil {
				logger.Log.Errorf("Change node(%v) status into losted failed, err is %s", lostedList, err.Error())
			}
		}
	}()
	return nodes
}
