package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
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
		clusterService:    NewClusterService(),
		clusterStatusRepo: repository.NewClusterStatusRepository(),
		messageService:    NewMessageService(),
	}
}

type clusterUpgradeService struct {
	clusterService    ClusterService
	clusterStatusRepo repository.ClusterStatusRepository
	messageService    MessageService
}

func (c *clusterUpgradeService) Upgrade(upgrade dto.ClusterUpgrade) error {
	cluster, err := c.clusterService.Get(upgrade.ClusterName)
	if err != nil {
		return fmt.Errorf("can not get cluster %s error %s", upgrade.ClusterName, err.Error())
	}
	if !(cluster.Source == constant.ClusterSourceLocal) {
		return errors.New("CLUSTER_IS_NOT_LOCAL")
	}
	if cluster.Status != constant.StatusRunning && cluster.Status != constant.StatusFailed {
		return fmt.Errorf("cluster status error %s", cluster.Status)
	}

	tx := db.DB.Begin()
	//从错误后继续
	if cluster.Cluster.Status.Phase == constant.StatusFailed && cluster.Cluster.Status.PrePhase == constant.StatusUpgrading {
		if err := tx.Model(model.ClusterStatusCondition{}).Where(model.ClusterStatusCondition{ClusterStatusID: cluster.StatusID, Status: constant.ConditionFalse}).Updates(map[string]interface{}{
			"Status":  constant.ConditionUnknown,
			"Message": "",
		}).Error; err != nil {
			return fmt.Errorf("reset status error %s", err.Error())
		}
	} else {
		if err := tx.Delete(&model.ClusterStatusCondition{}, "cluster_status_id = ?", cluster.StatusID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("reset contidion err %s", err.Error())
		}
	}
	// 修改状态
	cluster.Cluster.Status.PrePhase = cluster.Status
	cluster.Cluster.Status.Phase = constant.StatusUpgrading
	if err := tx.Save(&cluster.Cluster.Status).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("change status err %s", err.Error())
	}
	// 创建日志
	logId, writer, err := ansible.CreateAnsibleLogWriter(cluster.Name)
	if err != nil {
		return fmt.Errorf("create log error %s", err.Error())
	}
	cluster.LogId = logId
	if err := tx.Save(&cluster.Cluster).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("save cluster error %s", err.Error())
	}
	cluster.Spec.UpgradeVersion = upgrade.Version
	if err := tx.Save(&cluster.Spec).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("save cluster spec error %s", err.Error())
	}
	tx.Commit()
	go c.do(&cluster.Cluster, writer)
	return nil
}

func (c *clusterUpgradeService) do(cluster *model.Cluster, writer io.Writer) {

	status, err := c.clusterService.GetStatus(cluster.Name)
	if err != nil {
		log.Errorf("can not get current cluster status, error: %s", err.Error())
	}
	cluster.Status = status.ClusterStatus
	ctx, cancel := context.WithCancel(context.Background())
	admCluster := adm.NewCluster(*cluster, writer)
	statusChan := make(chan adm.Cluster)
	go c.doUpgrade(ctx, *admCluster, statusChan)
	for {
		cluster := <-statusChan
		// 保存进度
		_ = c.clusterStatusRepo.Save(&cluster.Status)
		switch cluster.Status.Phase {
		case constant.StatusRunning:
			cluster.Spec.Version = cluster.Spec.UpgradeVersion
			db.DB.Save(&cluster.Spec)
			cancel()
			return
		case constant.StatusFailed:
			cancel()
			return
		}
	}
}
func (c clusterUpgradeService) doUpgrade(ctx context.Context, cluster adm.Cluster, statusChan chan adm.Cluster) {
	ad := adm.NewClusterAdm()
	for {
		resp, err := ad.OnUpgrade(cluster)
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
