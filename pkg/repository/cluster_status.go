package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterStatusRepository interface {
	Get(id string) (model.ClusterStatus, error)
	Save(status *model.ClusterStatus, conditions ...model.ClusterStatusCondition) error
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
		Related(&status.Conditions).Error; err != nil {
		return status, err
	}
	return status, nil
}

func (c clusterStatusRepository) Save(status *model.ClusterStatus, conditions ...model.ClusterStatusCondition) error {
	tx := db.DB.Begin()
	if db.DB.NewRecord(status) {
		if err := db.DB.Create(&status).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(&status).Error; err != nil {
			return err
		}
	}
	for i, _ := range conditions {
		conditions[i].StatusID = status.ID
		err := c.conditionRepo.Save(&conditions[i])
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (c clusterStatusRepository) Delete(id string) error {
	status := model.ClusterStatus{ID: id}
	if err := db.DB.First(&status).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&status).Error; err != nil {
		return err
	}
	return nil
}
