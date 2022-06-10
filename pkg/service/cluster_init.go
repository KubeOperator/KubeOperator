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
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
)

type ClusterInitService interface {
	Init(cluster model.Cluster, writer io.Writer)
	GatherKubernetesToken(cluster model.Cluster) error
}

func NewClusterInitService() ClusterInitService {
	return &clusterInitService{
		clusterRepo:                repository.NewClusterRepository(),
		clusterNodeRepo:            repository.NewClusterNodeRepository(),
		clusterStatusRepo:          repository.NewClusterStatusRepository(),
		clusterSecretRepo:          repository.NewClusterSecretRepository(),
		clusterStatusConditionRepo: repository.NewClusterStatusConditionRepository(),
		clusterSpecRepo:            repository.NewClusterSpecRepository(),
		messageService:             NewMessageService(),
		clusterCreateHelper:        NewClusterCreateHelper(),
	}
}

type clusterInitService struct {
	clusterRepo                repository.ClusterRepository
	clusterNodeRepo            repository.ClusterNodeRepository
	clusterStatusRepo          repository.ClusterStatusRepository
	clusterSecretRepo          repository.ClusterSecretRepository
	clusterStatusConditionRepo repository.ClusterStatusConditionRepository
	clusterSpecRepo            repository.ClusterSpecRepository
	messageService             MessageService
	clusterCreateHelper        ClusterCreateHelper
}

func (c clusterInitService) Init(cluster model.Cluster, writer io.Writer) {
	if cluster.Provider == constant.ClusterProviderPlan {
		cluster.Status.Phase = constant.ClusterCreating
		_ = c.clusterStatusRepo.Save(&cluster.Status)
		if err := c.clusterCreateHelper.LoadPlanNodes(&cluster); err != nil {
			cluster.Status.Phase = constant.ClusterFailed
			cluster.Status.Message = err.Error()
			_ = c.clusterStatusRepo.Save(&cluster.Status)
			logger.Log.Errorf("init cluster resource for create failed: %s", err.Error())
			_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterInstall, false, err.Error()), cluster.Name, constant.ClusterInstall)
			return
		}
	}

	cluster.Nodes, _ = c.clusterNodeRepo.List(cluster.Name)
	ctx, cancel := context.WithCancel(context.Background())
	statusChan := make(chan adm.Cluster)
	cluster.Status.Phase = constant.ClusterInitializing
	_ = c.clusterStatusRepo.Save(&cluster.Status)

	admCluster := adm.NewCluster(cluster, writer)
	go c.doCreate(ctx, *admCluster, statusChan)
	for {
		cluster := <-statusChan
		_ = c.clusterStatusRepo.Save(&cluster.Status)
		switch cluster.Status.Phase {
		case constant.ClusterFailed:
			cancel()
			logger.Log.Errorf("cluster install failed: %s", cluster.Status.Message)
			_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterInstall, false, cluster.Status.Message), cluster.Name, constant.ClusterInstall)
			return
		case constant.ClusterRunning:
			logger.Log.Infof("cluster %s install successful!", cluster.Name)
			_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterInstall, true, ""), cluster.Name, constant.ClusterInstall)
			firstMasterIP := ""
			for i := range cluster.Nodes {
				if cluster.Nodes[i].Role == constant.NodeRoleNameMaster && len(firstMasterIP) == 0 {
					firstMasterIP = cluster.Nodes[i].Host.Ip
				}
				cluster.Nodes[i].Status = constant.ClusterRunning
				_ = c.clusterNodeRepo.Save(&cluster.Nodes[i])
			}
			cluster.SpecConf.KubeRouter = firstMasterIP
			if cluster.SpecConf.LbMode == constant.LbModeInternal {
				cluster.SpecConf.LbKubeApiserverIp = firstMasterIP
			}
			_ = c.clusterSpecRepo.SaveConf(&cluster.SpecConf)

			logger.Log.Infof("start to load tools ...")
			if err := c.loadTools(&cluster.Cluster); err != nil {
				logger.Log.Infof("load tool failed, err: %v!", err)
			} else {
				logger.Log.Infof("load tool successful !")
			}
			cancel()
			err := c.GatherKubernetesToken(cluster.Cluster)
			if err != nil {
				cluster.Status.Phase = constant.ClusterNotConnected
				cluster.Status.Message = err.Error()
			}
			return
		}
	}
}

func (c clusterInitService) doCreate(ctx context.Context, cluster adm.Cluster, statusChan chan adm.Cluster) {
	ad := adm.NewClusterAdm()
	for {
		resp, err := ad.OnInitialize(cluster)
		if err != nil {
			cluster.Status.Message = err.Error()
		}
		cluster.Status = resp.Status
		select {
		case <-ctx.Done():
			return
		case statusChan <- cluster:
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

	if cluster.Architectures == "amd64" {
		for _, istio := range cluster.PrepareIstios() {
			err := tx.Create(&istio).Error
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("can not prepare cluster istio %s reason %s", istio.Name, err.Error())
			}
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
