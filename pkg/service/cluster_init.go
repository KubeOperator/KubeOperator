package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
)

type ClusterInitService interface {
	Init(cluster model.Cluster, writer io.Writer)
	GatherKubernetesToken(cluster model.Cluster) error
}

func NewClusterInitService() ClusterInitService {
	return &clusterInitService{
		clusterRepo:        repository.NewClusterRepository(),
		clusterNodeRepo:    repository.NewClusterNodeRepository(),
		clusterSecretRepo:  repository.NewClusterSecretRepository(),
		clusterSpecRepo:    repository.NewClusterSpecRepository(),
		taskLogService:     NewTaskLogService(),
		clusterIaasService: NewClusterIaasService(),
		msgService:         NewMsgService(),
	}
}

type clusterInitService struct {
	clusterRepo        repository.ClusterRepository
	clusterNodeRepo    repository.ClusterNodeRepository
	clusterSecretRepo  repository.ClusterSecretRepository
	clusterSpecRepo    repository.ClusterSpecRepository
	taskLogService     TaskLogService
	clusterIaasService ClusterIaasService
	msgService         MsgService
}

func (c clusterInitService) Init(cluster model.Cluster, writer io.Writer) {
	cluster.TaskLog.Phase = constant.TaskLogStatusWaiting
	_ = c.taskLogService.Save(&cluster.TaskLog)

	if cluster.Provider == constant.ClusterProviderPlan {
		if err := c.clusterIaasService.LoadPlanNodes(&cluster); err != nil {
			_ = c.taskLogService.End(&cluster.TaskLog, false, err.Error())
			cluster.Status = constant.StatusFailed
			cluster.Message = err.Error()
			_ = c.clusterRepo.Save(&cluster)
			logger.Log.Errorf("init cluster resource for create failed: %s", err.Error())
			_ = c.msgService.SendMsg(constant.ClusterInstall, constant.System, cluster, false, map[string]string{"errMsg": err.Error()})
			return
		}
	}

	cluster.Status = constant.StatusInitializing
	cluster.CurrentTaskID = cluster.TaskLog.ID
	_ = c.clusterRepo.Save(&cluster)
	cluster.TaskLog.Phase = constant.TaskLogStatusRunning
	cluster.TaskLog.CreatedAt = time.Now()
	_ = c.taskLogService.Save(&cluster.TaskLog)
	cluster.Nodes, _ = c.clusterNodeRepo.List(cluster.Name)
	ctx, cancel := context.WithCancel(context.Background())
	statusChan := make(chan adm.AnsibleHelper)

	admCluster := adm.NewAnsibleHelper(cluster, writer)
	admCluster.Kobe.SetVar(facts.ComponentOptionFactName, "cluster")
	go c.doCreate(ctx, *admCluster, statusChan)
	for {
		result := <-statusChan
		switch cluster.TaskLog.Phase {
		case constant.TaskLogStatusFailed:
			if err := c.taskLogService.End(&cluster.TaskLog, false, result.Message); err != nil {
				logger.Log.Infof("save task failed %v", err)
			}
			cancel()
			cluster.Status = constant.StatusFailed
			cluster.Message = result.Message
			_ = c.clusterRepo.Save(&cluster)
			logger.Log.Errorf("cluster install failed: %s", cluster.TaskLog.Message)
			_ = c.msgService.SendMsg(constant.ClusterInstall, constant.System, cluster, false, map[string]string{"errMsg": cluster.TaskLog.Message})
			return
		case constant.TaskLogStatusSuccess:
			if err := c.taskLogService.End(&cluster.TaskLog, true, ""); err != nil {
				logger.Log.Infof("save task failed %v", err)
			}
			logger.Log.Infof("cluster %s install successful!", cluster.Name)
			cluster.Status = constant.StatusRunning
			cluster.Message = result.Message
			cluster.CurrentTaskID = ""
			_ = c.msgService.SendMsg(constant.ClusterInstall, constant.System, cluster, true, map[string]string{})
			firstMasterIP := ""
			for i := range cluster.Nodes {
				if cluster.Nodes[i].Role == constant.NodeRoleNameMaster && len(firstMasterIP) == 0 {
					firstMasterIP = cluster.Nodes[i].Host.Ip
				}
				cluster.Nodes[i].Status = constant.StatusRunning
				_ = c.clusterNodeRepo.Save(&cluster.Nodes[i])
			}
			cluster.SpecConf.KubeRouter = firstMasterIP
			if cluster.SpecConf.LbMode == constant.LbModeInternal {
				cluster.SpecConf.LbKubeApiserverIp = firstMasterIP
			}
			_ = c.clusterSpecRepo.SaveConf(&cluster.SpecConf)

			logger.Log.Infof("start to load tools ...")
			if err := c.loadTools(&cluster); err != nil {
				logger.Log.Infof("load tool failed, err: %v!", err)
			}
			cancel()
			err := c.GatherKubernetesToken(cluster)
			if err != nil {
				cluster.Status = constant.ClusterNotConnected
				cluster.Message = err.Error()
			}
			_ = c.clusterRepo.Save(&cluster)
			return
		default:
			cluster.TaskLog.Phase = result.Status
			cluster.TaskLog.Message = result.Message
			cluster.TaskLog.Details = result.LogDetail
			if err := c.taskLogService.Save(&cluster.TaskLog); err != nil {
				logger.Log.Infof("save task failed %v", err)
			}
		}
	}
}

func (c clusterInitService) doCreate(ctx context.Context, aHelper adm.AnsibleHelper, statusChan chan adm.AnsibleHelper) {
	ad := adm.NewClusterAdm()
	for {
		if err := ad.OnInitialize(&aHelper); err != nil {
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

func (c clusterInitService) loadTools(cluster *model.Cluster) error {
	var (
		manifest model.ClusterManifest
		toolVars []model.VersionHelp
	)
	tx := db.DB.Begin()
	if err := tx.Where("name = ?", cluster.Version).First(&manifest).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("can find manifest version: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(manifest.ToolVars), &toolVars); err != nil {
		tx.Rollback()
		return fmt.Errorf("unmarshal manifest.toolvar error %s", err.Error())
	}
	for _, tool := range cluster.PrepareTools() {
		for _, item := range toolVars {
			if tool.Name == item.Name {
				tool.Version = item.Version
				break
			}
		}
		err := tx.Create(&tool).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("can not prepare cluster tool %s reason %s", tool.Name, err.Error())
		}
	}
	tx.Commit()

	return nil
}

func (c clusterInitService) GatherKubernetesToken(cluster model.Cluster) error {
	secret, err := c.clusterSecretRepo.Get(cluster.SecretID)
	if err != nil {
		return err
	}
	master, err := c.clusterNodeRepo.FirstMaster(cluster.ID)
	if err != nil {
		return err
	}
	sshConfig := master.ToSSHConfig()
	client, err := ssh.New(&sshConfig)
	if err != nil {
		return err
	}
	token, err := clusterUtil.GetClusterToken(client)
	if err != nil {
		return err
	}
	secret.KubernetesToken = token
	return c.clusterSecretRepo.Save(&secret)
}
