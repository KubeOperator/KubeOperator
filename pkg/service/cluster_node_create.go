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
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
)

func (c clusterNodeService) batchCreate(cluster *model.Cluster, currentNodes []model.ClusterNode, item dto.NodeBatch) error {
	tasklog := model.TaskLog{
		ClusterID: cluster.ID,
		Type:      constant.TaskLogTypeClusterNodeExtend,
		Phase:     constant.StatusWaiting,
	}
	if err := c.taskLogService.Start(&tasklog); err != nil {
		return err
	}
	cluster.TaskLog = tasklog
	cluster.CurrentTaskID = tasklog.ID
	_ = c.clusterRepo.Save(cluster)

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
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, cluster.TaskLog.ID)
	if err != nil {
		return err
	}
	cluster.Nodes = append(cluster.Nodes, newNodes...)

	go c.addWorkInit(cluster, newNodes, writer, "")
	return nil
}

func (c *clusterNodeService) addWorkInit(cluster *model.Cluster, nodes []model.ClusterNode, writer io.Writer, operation string) {
	var (
		newNodeIDs []string
		newHostIDs []string
	)
	for _, n := range nodes {
		newNodeIDs = append(newNodeIDs, n.ID)
		newHostIDs = append(newHostIDs, n.Host.ID)
	}

	if operation != "recreate" && cluster.Provider == constant.ClusterProviderPlan {
		logger.Log.Info("cluster-plan start add hosts, update hosts status and infos")
		if err := c.updataHostInfo(cluster, newNodeIDs, newHostIDs); err != nil {
			c.updateNodeStatus(cluster, constant.ClusterAddWorker, constant.StatusFailed, newNodeIDs, err)
			return
		}
	}

	cluster.TaskLog.Phase = constant.TaskLogStatusRunning
	cluster.TaskLog.CreatedAt = time.Now()
	_ = c.taskLogService.Save(&cluster.TaskLog)

	logger.Log.Info("start binding nodes to cluster")
	var (
		nodeIds   []string
		nodeNames []string
	)
	for _, n := range nodes {
		nodeIds = append(nodeIds, n.ID)
	}

	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", nodeIds).
		Updates(map[string]interface{}{"Status": constant.StatusInitializing, "CurrentTaskID": cluster.TaskLog.ID}).Error; err != nil {
		c.updateNodeStatus(cluster, constant.ClusterAddWorker, constant.StatusFailed, newNodeIDs, err)
	}

	for _, n := range nodes {
		nodeIds = append(nodeIds, n.ID)
		nodeNames = append(nodeNames, n.Name)
	}
	admCluster := adm.NewAnsibleHelperWithNewWorker(*cluster, nodeNames, writer)
	statusChan := make(chan adm.AnsibleHelper)
	ctx, cancel := context.WithCancel(context.Background())

	go c.doCreate(ctx, *admCluster, statusChan)
	for {
		result := <-statusChan
		cluster.TaskLog.Phase = result.Status
		cluster.TaskLog.Message = result.Message
		cluster.TaskLog.Details = result.LogDetail
		_ = c.taskLogService.Save(&cluster.TaskLog)
		// 保存进度
		switch result.Status {
		case constant.TaskLogStatusSuccess:
			cancel()
			c.updateNodeStatus(cluster, constant.ClusterAddWorker, constant.StatusRunning, nodeIds, fmt.Errorf(result.Message))
			cluster.CurrentTaskID = ""
			_ = c.clusterRepo.Save(cluster)
			return
		case constant.TaskLogStatusFailed:
			cancel()
			c.updateNodeStatus(cluster, constant.ClusterAddWorker, constant.StatusFailed, nodeIds, fmt.Errorf(result.Message))
			return
		}
	}
}

func (c clusterNodeService) doCreate(ctx context.Context, aHelper adm.AnsibleHelper, statusChan chan adm.AnsibleHelper) {
	ad := adm.NewClusterAdm()
	for {
		if err := ad.OnAddWorker(&aHelper); err != nil {
			aHelper.Message = err.Error()
		}
		select {
		case <-ctx.Done():
			return
		case statusChan <- aHelper:
		}
		time.Sleep(5 * time.Second)
	}
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
			Status:    constant.StatusWaiting,
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
		newNodes[i].CurrentTaskID = cluster.CurrentTaskID
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
			Status: constant.StatusCreating,
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

func (c *clusterNodeService) updateNodeStatus(cluster *model.Cluster, operation, status string, nodeIDs []string, errMsg error) {
	errmsg, taskSuccess := "", true
	if errMsg != nil {
		errmsg = errMsg.Error()
	}
	if status == constant.StatusFailed {
		taskSuccess = false
		_ = c.messageService.SendMessage(constant.System, false, GetContent(operation, false, errmsg), cluster.Name, operation)
	}
	_ = c.taskLogService.End(&cluster.TaskLog, taskSuccess, errmsg)

	logger.Log.Infof("update node status to %s, msg: %v", status, errMsg)
	if err := db.DB.Model(&model.ClusterNode{}).Where("id in (?)", nodeIDs).
		Updates(map[string]interface{}{"Status": status, "Message": errmsg}).Error; err != nil {
		logger.Log.Errorf("can not update node status %s", err.Error())
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
				"Status":  constant.StatusFailed,
				"Message": fmt.Errorf("can not create hosts reason %s", err.Error()),
				"HostID":  "",
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
