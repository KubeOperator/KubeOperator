package service

import (
	"context"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
		ClusterService: NewClusterService(),
		clusterRepo:    repository.NewClusterRepository(),
		NodeRepo:       repository.NewClusterNodeRepository(),
		taskLogService: NewTaskLogService(),
		HostRepo:       repository.NewHostRepository(),
		planService:    NewPlanService(),
		vmConfigRepo:   repository.NewVmConfigRepository(),
		ntpServerRepo:  repository.NewNtpServerRepository(),
		messageService: NewMessageService(),
		hostService:    NewHostService(),
	}
}

type clusterNodeService struct {
	ClusterService ClusterService
	clusterRepo    repository.ClusterRepository
	NodeRepo       repository.ClusterNodeRepository
	taskLogService TaskLogService
	HostRepo       repository.HostRepository
	planService    PlanService
	vmConfigRepo   repository.VmConfigRepository
	ntpServerRepo  repository.NtpServerRepository
	messageService MessageService
	hostService    HostService
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
	count, mNodes, err := c.NodeRepo.Page(num, size, clusterName)
	if err != nil {
		return nil, err
	}

	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return nil, err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return nil, err
	}

	kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
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
	mNodes, err := c.NodeRepo.List(clusterName)
	if err != nil {
		return nil, err
	}
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return nil, err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return nil, err
	}
	kubeNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodesAfterSync := syncNodeStatus(mNodes, kubeNodes, cluster.Source, "true")
	return nodesAfterSync, nil
}

func (c *clusterNodeService) Batch(clusterName string, item dto.NodeBatch) error {
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Plan", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential", "Nodes.Host.Zone", "MultiClusterRepositories"})
	if err != nil {
		return fmt.Errorf("can not found %s", clusterName)
	}
	var currentNodes []model.ClusterNode
	if err := db.DB.Where("cluster_id = ?", cluster.ID).Preload("Host").Preload("Host.Credential").Preload("Host.Zone").Find(&currentNodes).Error; err != nil {
		return fmt.Errorf("can not read cluster %s current nodes %s", cluster.Name, err.Error())
	}
	for _, node := range currentNodes {
		if node.Status == constant.StatusCreating || node.Status == constant.StatusInitializing || node.Status == constant.StatusWaiting {
			return errors.New("NODE_ALREADY_RUNNING_TASK")
		}
	}
	isON := c.taskLogService.IsTaskOn(clusterName)
	if isON {
		return errors.New("TASK_IN_EXECUTION")
	}
	switch item.Operation {
	case constant.BatchOperationCreate:
		return c.batchCreate(&cluster, currentNodes, item)
	case constant.BatchOperationDelete:
		if err := db.DB.Model(&model.ClusterNode{}).Where("name in (?)", item.Nodes).
			Updates(map[string]interface{}{"Status": constant.StatusTerminating, "Message": ""}).Error; err != nil {
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
		tasklog, err := c.taskLogService.NewTerminalTask(cluster.ID, constant.TaskLogTypeClusterNodeShrink)
		if err != nil {
			return err
		}
		cluster.TaskLog = *tasklog
		cluster.CurrentTaskID = tasklog.ID
		_ = c.clusterRepo.Save(&cluster)

		go c.removeNodes(&cluster, item, currentNodes, nodesForDelete)
		return nil
	}
	return nil
}

func (c *clusterNodeService) removeNodes(cluster *model.Cluster, item dto.NodeBatch, currentNodes, nodesForDelete []model.ClusterNode) {
	var (
		nodeIDs []string
		hostIDs []string
		hostIPs []string
	)
	for _, node := range nodesForDelete {
		hostIDs = append(hostIDs, node.Host.ID)
		hostIPs = append(hostIPs, node.Host.Ip)
		nodeIDs = append(nodeIDs, node.ID)
	}

	if cluster.Provider == constant.ClusterProviderPlan {
		var p model.Plan
		if err := db.DB.Where("id = ?", cluster.PlanID).First(&p).Error; err != nil {
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
			return
		}
		planDTO, err := c.planService.Get(p.Name)
		if err != nil {
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
			return
		}
		cluster.Plan = planDTO.Plan

		logger.Log.Info("start run removeWorkerPlaybook")
		if err := c.runDeleteWorkerPlaybook(cluster, nodesForDelete, removeWorkerPlaybook); err != nil {
			if !item.IsForce {
				c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
				return
			} else {
				logger.Log.Errorf("run runDeleteWorkerPlaybook failed, err: %s", err.Error())
			}
		}
		if err := c.destroyHosts(cluster, currentNodes, nodeIDs); err != nil {
			if !item.IsForce {
				c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
				return
			} else {
				logger.Log.Errorf("destroy host failed, err: %s", err.Error())
			}
		}
		logger.Log.Info("delete all nodes successful! now start updata cluster datas")

		tx := db.DB.Begin()
		if err := tx.Where("id in (?)", hostIDs).Delete(&model.Host{}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
			return
		}
		if err := tx.Where("resource_id in (?) AND resource_type = ?", hostIDs, constant.ResourceHost).
			Delete(&model.ProjectResource{}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
			return
		}
		if err := tx.Where("resource_id in (?) AND resource_type = ?", hostIDs, constant.ResourceHost).
			Delete(&model.ClusterResource{}).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
			return
		}
		if err := tx.Model(&model.Ip{}).Where("address in (?)", hostIPs).
			Update("status", constant.IpAvailable).Error; err != nil {
			tx.Rollback()
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
			return
		}
		tx.Commit()
	} else {
		logger.Log.Info("start run removeWorkerPlaybook")
		if err := c.runDeleteWorkerPlaybook(cluster, nodesForDelete, removeWorkerPlaybook); err != nil {
			if !item.IsForce {
				c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
				return
			} else {
				logger.Log.Errorf("run removeWorkerPlaybook failed, err: %s", err.Error())
			}
		}
		logger.Log.Info("start run resetWorkerPlaybook")
		if err := c.runDeleteWorkerPlaybook(cluster, nodesForDelete, resetWorkerPlaybook); err != nil {
			if !item.IsForce {
				c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
				return
			} else {
				logger.Log.Errorf("run resetWorkerPlaybook failed, err: %s", err.Error())
			}
		}
		logger.Log.Info("delete all nodes successful! now start updata cluster datas")
		if err := db.DB.Model(&model.Host{}).Where("id in (?)", hostIDs).Update(map[string]interface{}{"ClusterID": ""}).Error; err != nil {
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
			return
		}
		if err := db.DB.Where("resource_id in (?) AND resource_type = ?", hostIDs, constant.ResourceHost).
			Delete(&model.ClusterResource{}).Error; err != nil {
			c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
			return
		}
	}
	if err := db.DB.Where("id in (?)", nodeIDs).Delete(&model.ClusterNode{}).Error; err != nil {
		c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusFailed, nodeIDs, err)
		return
	}
	c.updateNodeStatus(cluster, constant.ClusterRemoveWorker, constant.StatusRunning, nodeIDs, nil)
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

const removeWorkerPlaybook = "96-remove-worker.yml"
const resetWorkerPlaybook = "97-reset-worker.yml"

func (c *clusterNodeService) runDeleteWorkerPlaybook(cluster *model.Cluster, nodes []model.ClusterNode, playbookName string) error {
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
	k.SetVar(facts.NtpServerFactName, ntps)
	if err = phases.RunPlaybookAndGetResult(k, playbookName, "", writer); err != nil {
		return err
	}
	return nil
}

func (c *clusterNodeService) Recreate(clusterName string, batch dto.NodeBatch) error {
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	tasklog, err := c.taskLogService.GetByID(cluster.CurrentTaskID)
	if err != nil {
		return err
	}
	cluster.TaskLog = tasklog
	if err := c.taskLogService.RestartTask(&cluster, constant.TaskLogTypeClusterNodeExtend); err != nil {
		return err
	}

	var nodes []model.ClusterNode
	if err := db.DB.Where("current_task_id = ?", cluster.CurrentTaskID).Find(&nodes).Error; err != nil {
		return err
	}
	if err := db.DB.Model(&model.ClusterNode{}).Where("current_task_id = ?", batch.StatusID).
		Updates(map[string]interface{}{"Status": constant.StatusInitializing}).Error; err != nil {
		return fmt.Errorf("can not update cluster status %s", err.Error())
	}

	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, cluster.TaskLog.ID)
	if err != nil {
		return err
	}

	logger.Log.WithFields(logrus.Fields{
		"log_id": cluster.TaskLog.ID,
	}).Debugf("get ansible writer log of cluster %s successful, now start to init the cluster", cluster.Name)

	go c.addWorkInit(&cluster, nodes, writer, "recreate")
	return nil
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
