package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterStatusConditionRepository interface {
	Save(condition *model.ClusterStatusCondition) error
	Delete(id string) error
	List(clusterStatusId string) ([]model.ClusterStatusCondition, error)
}

func NewClusterStatusConditionRepository() ClusterStatusConditionRepository {
	return &clusterStatusConditionRepository{}
}

type clusterStatusConditionRepository struct{}

func (c clusterStatusConditionRepository) List(clusterStatusId string) ([]model.ClusterStatusCondition, error) {
	var clusterStatusConditions []model.ClusterStatusCondition
	if err := db.DB.Where(&model.ClusterStatusCondition{ClusterStatusID: clusterStatusId}).
		Find(&clusterStatusConditions).Error; err != nil {
		return nil, err
	}
	return clusterStatusConditions, nil
}

func (c clusterStatusConditionRepository) Delete(id string) error {
	condition := model.ClusterStatusCondition{
		ID: id,
	}
	if err := db.DB.First(&condition).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&condition).Error; err != nil {
		return err
	}
	return nil
}

func (c clusterStatusConditionRepository) Save(condition *model.ClusterStatusCondition) error {
	if db.DB.NewRecord(condition) {
		var temp model.ClusterStatusCondition
		if db.DB.Where(&model.ClusterStatusCondition{ClusterStatusID: condition.ClusterStatusID, Name: condition.Name}).
			First(&temp).
			RecordNotFound() {
			if err := db.DB.Create(condition).Error; err != nil {
				return err
			}
		} else {
			condition.ID = temp.ID
			if err := db.DB.Save(condition).Error; err != nil {
				return err
			}
		}
	} else {
		if err := db.DB.Save(condition).Error; err != nil {
			return err
		}
	}
	return nil
}
