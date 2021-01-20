package hook

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
)

func init() {
	BeforeApplicationStart.AddFunc(recoverClusterTask)
	BeforeApplicationStart.AddFunc(markClusterNodeDirtyData)
}

var clusterService = service.NewClusterService()

// cluster
func recoverClusterTask() error {
	clusters, err := clusterService.List()
	if err != nil {
		return err
	}

	tx := db.DB.Begin()
	for _, cluster := range clusters {
		if cluster.Status == constant.StatusCreating || cluster.Status == constant.StatusTerminating || cluster.Status == constant.StatusInitializing {
			var status model.ClusterStatus
			if err := db.DB.Where(model.ClusterStatus{ID: cluster.StatusID}).First(&status).Error; err != nil {
				return err
			}
			status.PrePhase = status.Phase
			status.Phase = constant.StatusFailed
			if err := tx.Save(&status).Error; err != nil {
				tx.Rollback()
				return err
			}
			var conditions []model.ClusterStatusCondition
			if err := db.DB.Where(model.ClusterStatusCondition{ClusterStatusID: status.ID}).Order("last_probe_time asc").Find(&conditions).Error; err != nil {
				return err
			}
			if len(conditions) > 0 {
				for i := range conditions {
					if conditions[i].Status == constant.ConditionUnknown {
						conditions[i].Status = constant.ConditionFalse
						conditions[i].Message = "task cancel"
					}
					if err := tx.Save(&conditions[i]).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}
	}
	tx.Commit()
	return nil
}

// cluster node
func markClusterNodeDirtyData() error {
	var status = []string{constant.StatusTerminating, constant.StatusInitializing, constant.StatusCreating, constant.StatusWaiting}
	if err := db.DB.Model(&model.ClusterNode{}).Where("status in (?)", status).Updates(map[string]interface{}{"dirty": 1}).Error; err != nil {
		return err
	}
	return nil
}
