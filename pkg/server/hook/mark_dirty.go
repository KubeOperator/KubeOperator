package hook

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

func init() {
	BeforeApplicationStart.AddFunc(recoverClusterTask)
}

var stableStatus = []string{constant.StatusRunning, constant.StatusFailed, constant.StatusNotReady, constant.StatusLost}

// cluster
func recoverClusterTask() error {
	var taskList []model.TaskLog

	logger.Log.Info("Update status to failed caused by task cancel")
	tx := db.DB.Begin()
	if err := db.DB.Where("phase not in (?)", stableStatus).Find(&taskList).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, task := range taskList {
		task.PrePhase = task.Phase
		task.Phase = constant.StatusFailed
		task.Message = constant.TaskCancel
		if err := tx.Save(&task).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Model(&model.TaskLogDetail{}).Where("status not in (?)", stableStatus).Updates(map[string]interface{}{
		"status":  constant.StatusFailed,
		"message": constant.TaskCancel,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&model.TaskLogDetail{}).Where("status = ?", constant.ConditionUnknown).Updates(map[string]interface{}{
		"status":  constant.ConditionFalse,
		"message": constant.TaskCancel,
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
		statu.PreStatus = statu.Status
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
