package repository

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type ClusterStatusRepository interface {
	Get(id string) (model.ClusterStatus, error)
	Save(status *model.ClusterStatus) error
	Delete(id string) error
}

func NewClusterStatusRepository() ClusterStatusRepository {
	return &clusterStatusRepository{
		conditionRepo: NewClusterStatusConditionRepository(),
	}
}

type clusterStatusRepository struct {
	conditionRepo ClusterStatusConditionRepository
}

func (c clusterStatusRepository) Get(id string) (model.ClusterStatus, error) {
	status := model.ClusterStatus{
		ID: id,
	}
	if err := db.DB.
		First(&status).
		Order("last_probe_time asc").
		Related(&status.ClusterStatusConditions).
		Error; err != nil {
		return status, err
	}
	return status, nil
}

func (c clusterStatusRepository) Save(status *model.ClusterStatus) error {
	tx := db.DB.Begin()
	if tx.NewRecord(status) {
		if err := tx.Create(&status).Error; err != nil {
			return err
		}
	} else {
		var oldStatus model.ClusterStatus
		db.DB.First(&oldStatus)
		if status.Phase != oldStatus.Phase {
			status.PrePhase = oldStatus.Phase
		}
		if err := db.DB.Save(&status).Error; err != nil {
			return err
		}
	}
	if err := tx.Delete(&model.ClusterStatusCondition{}, "cluster_status_id = ?", status.ID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("reset contidion err %s", err.Error())
	}
	for i := range status.ClusterStatusConditions {
		status.ClusterStatusConditions[i].ClusterStatusID = status.ID
		if tx.NewRecord(status.ClusterStatusConditions[i]) {
			var temp model.ClusterStatusCondition
			if tx.Where(model.ClusterStatusCondition{ClusterStatusID: status.ClusterStatusConditions[i].ClusterStatusID, Name: status.ClusterStatusConditions[i].Name}).
				First(&temp).
				RecordNotFound() {
				status.ClusterStatusConditions[i].CreatedAt = time.Now()
				status.ClusterStatusConditions[i].UpdatedAt = time.Now()
				status.ClusterStatusConditions[i].ID = uuid.NewV4().String()
				if err := tx.Create(status.ClusterStatusConditions[i]).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				status.ClusterStatusConditions[i].ID = temp.ID
				status.ClusterStatusConditions[i].UpdatedAt = time.Now()
				if err := tx.Save(status.ClusterStatusConditions[i]).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		} else {
			if err := tx.Save(status.ClusterStatusConditions[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (c clusterStatusRepository) Delete(id string) error {
	if err := db.DB.
		First(&model.Cluster{ID: id}).
		Delete(model.Cluster{}).Error; err != nil {
		return err
	}
	return nil
}
