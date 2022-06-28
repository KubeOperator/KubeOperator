package hook

import (
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

func init() {
	BeforeApplicationStart.AddFunc(recoverClusterTask)
}

var stableStatus = []string{constant.StatusRunning, constant.StatusFailed, constant.StatusNotReady, constant.StatusLost}
var statleTaskStatus = []string{constant.TaskLogStatusSuccess, constant.TaskLogStatusFailed}
var stableDetailStatus = []string{constant.StatusRunning, constant.TaskLogStatusSuccess, constant.StatusFailed, constant.StatusNotReady, constant.StatusLost, constant.TaskDetailStatusFalse, constant.TaskDetailStatusTrue}

// cluster
func recoverClusterTask() error {
	logger.Log.Info("Update status to failed caused by task cancel")
	tx := db.DB.Begin()
	if err := db.DB.Model(&model.Cluster{}).Where("status not in (?)", stableStatus).Updates(map[string]interface{}{
		"status":  constant.StatusFailed,
		"message": constant.TaskCancel,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := db.DB.Model(&model.TaskLog{}).Where("phase not in (?)", statleTaskStatus).Updates(map[string]interface{}{
		"phase":    constant.TaskLogStatusFailed,
		"message":  constant.TaskCancel,
		"end_time": time.Now().Unix(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.TaskLogDetail{}).Where("status = ?", constant.TaskDetailStatusUnknown).Updates(map[string]interface{}{
		"status":   constant.TaskDetailStatusFalse,
		"message":  constant.TaskCancel,
		"end_time": time.Now().Unix(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.TaskLogDetail{}).Where("status not in (?) ", stableDetailStatus).Updates(map[string]interface{}{
		"status":   constant.StatusFailed,
		"message":  constant.TaskCancel,
		"end_time": time.Now().Unix(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Host{}).Where("status != ? AND status != ?", constant.StatusRunning, constant.StatusFailed).Updates(map[string]interface{}{
		"status":  constant.StatusFailed,
		"message": constant.TaskCancel,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.ClusterStorageProvisioner{}).Where("status not in (?)", stableStatus).Updates(map[string]interface{}{
		"status":  constant.StatusFailed,
		"message": constant.TaskCancel,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.ClusterGpu{}).Where("status = ? OR status = ? OR status = ?", constant.StatusInitializing, constant.StatusTerminating, constant.StatusWaiting).Updates(map[string]interface{}{
		"status":  constant.StatusFailed,
		"message": constant.TaskCancel,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	var nodes []model.ClusterNode
	if err := db.DB.Where("status not in (?) AND status != ''", stableStatus).Find(&nodes).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, statu := range nodes {
		statu.Status = constant.StatusFailed
		statu.Message = constant.TaskCancel
		if err := tx.Save(&statu).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	logger.Log.Info("update status successful !")
	return nil
}
