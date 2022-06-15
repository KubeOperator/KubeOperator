package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
)

type ClusterUpgradeService interface {
	Upgrade(upgrade dto.ClusterUpgrade) error
}

func NewClusterUpgradeService() ClusterUpgradeService {
	return &clusterUpgradeService{
		clusterService: NewClusterService(),
		messageService: NewMessageService(),
		clusterRepo:    repository.NewClusterRepository(),
		taskLogService: NewTaskLogService(),
	}
}

type clusterUpgradeService struct {
	clusterService ClusterService
	messageService MessageService
	clusterRepo    repository.ClusterRepository
	taskLogService TaskLogService
}

func (c *clusterUpgradeService) Upgrade(upgrade dto.ClusterUpgrade) error {
	loginfo, _ := json.Marshal(upgrade)
	logger.Log.WithFields(logrus.Fields{"cluster_upgrade_info": string(loginfo)}).Debugf("start to upgrade the cluster %s", upgrade.ClusterName)

	cluster, err := c.clusterRepo.GetWithPreload(upgrade.ClusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Nodes", "Nodes.Host", "Nodes.Host.Credential", "Nodes.Host.Zone", "MultiClusterRepositories"})
	if err != nil {
		return fmt.Errorf("can not get cluster %s error %s", upgrade.ClusterName, err.Error())
	}

	if cluster.Source == constant.ClusterSourceExternal {
		return errors.New("CLUSTER_IS_NOT_LOCAL")
	}
	if cluster.Status != constant.StatusRunning && cluster.Status != constant.StatusFailed {
		return fmt.Errorf("cluster status error %s", cluster.Status)
	}

	tx := db.DB.Begin()
	//从错误后继续
	if cluster.TaskLog.Phase == constant.StatusFailed && cluster.TaskLog.Type == constant.TaskLogTypeClusterUpgrade {
		if err := tx.Model(&model.TaskLogDetail{}).
			Where("task_log_id = ? AND status = ?", cluster.TaskLog.ID, constant.ConditionFalse).
			Updates(map[string]interface{}{
				"Status":  constant.ConditionUnknown,
				"Message": "",
			}).Error; err != nil {
			return fmt.Errorf("reset status error %s", err.Error())
		}
	} else {
		cluster.TaskLog = model.TaskLog{
			ClusterID: cluster.ID,
			Type:      constant.TaskLogTypeClusterUpgrade,
		}
	}
	// 修改状态
	cluster.TaskLog.Phase = constant.StatusUpgrading

	if err := c.taskLogService.Save(&cluster.TaskLog); err != nil {
		tx.Rollback()
		return fmt.Errorf("reset contidion err %s", err.Error())
	}
	// 创建日志
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, cluster.TaskLog.ID)
	if err != nil {
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUpgrade, false, err.Error()), cluster.Name, constant.ClusterUpgrade)
		return fmt.Errorf("create log error %s", err.Error())
	}
	cluster.UpgradeVersion = upgrade.Version
	if err := tx.Save(&cluster).Error; err != nil {
		tx.Rollback()
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUpgrade, false, err.Error()), cluster.Name, constant.ClusterUpgrade)
		return fmt.Errorf("save cluster spec error %s", err.Error())
	}
	// 更新工具版本状态
	if err := c.updateToolVersion(tx, upgrade.Version, cluster.ID); err != nil {
		tx.Rollback()
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUpgrade, false, err.Error()), cluster.Name, constant.ClusterUpgrade)
		return err
	}

	tx.Commit()

	logger.Log.Infof("update db data of cluster %s successful, now start to upgrade cluster", cluster.Name)
	go c.do(&cluster, writer)
	return nil
}

func (c *clusterUpgradeService) do(cluster *model.Cluster, writer io.Writer) {
	ctx, cancel := context.WithCancel(context.Background())
	admCluster := adm.NewCluster(*cluster, writer)
	statusChan := make(chan adm.AnsibleHelper)
	go c.doUpgrade(ctx, *admCluster, statusChan)
	for {
		result := <-statusChan
		// 保存进度
		cluster.Status = result.Status
		cluster.Message = result.Message
		_ = c.clusterRepo.Save(cluster)
		switch result.Status {
		case constant.StatusRunning:
			_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterUpgrade, true, ""), cluster.Name, constant.ClusterUpgrade)
			cluster.Version = cluster.UpgradeVersion
			db.DB.Save(&cluster)
			cancel()
			return
		case constant.StatusFailed:
			_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUpgrade, false, result.Message), cluster.Name, constant.ClusterUpgrade)
			cancel()
			return
		}
	}
}

func (c clusterUpgradeService) doUpgrade(ctx context.Context, aHelper adm.AnsibleHelper, statusChan chan adm.AnsibleHelper) {
	ad := adm.NewClusterAdm()
	for {
		resp, err := ad.OnUpgrade(aHelper)
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

func (c clusterUpgradeService) updateToolVersion(tx *gorm.DB, version, clusterID string) error {
	var (
		tools    []model.ClusterTool
		manifest model.ClusterManifest
		toolVars []model.VersionHelp
	)
	if err := tx.Where("name = ?", version).First(&manifest).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("get manifest error %s", err.Error())
	}
	if err := tx.Where("cluster_id = ?", clusterID).Find(&tools).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("get tools error %s", err.Error())
	}
	if err := json.Unmarshal([]byte(manifest.ToolVars), &toolVars); err != nil {
		return fmt.Errorf("unmarshal manifest.toolvar error %s", err.Error())
	}
	for _, tool := range tools {
		for _, item := range toolVars {
			if tool.Name == item.Name {
				if tool.Version != item.Version {
					if tool.Status == constant.ClusterWaiting {
						tool.Version = item.Version
					} else {
						tool.HigherVersion = item.Version
					}
					if err := tx.Save(&tool).Error; err != nil {
						return fmt.Errorf("update tool version error %s", err.Error())
					}
				}
				break
			}
		}
	}
	return nil
}
