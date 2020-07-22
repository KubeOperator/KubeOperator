package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ProjectResourceRepository interface {
	PageByProjectIdAndType(num, size int, projectId string, resourceType string) (int, []model.ProjectResource, error)
	Batch(operation string, items []model.ProjectResource) error
}

func NewProjectResourceRepository() ProjectResourceRepository {
	return &projectResourceRepository{}
}

type projectResourceRepository struct {
}

func (p projectResourceRepository) PageByProjectIdAndType(num, size int, projectId string, resourceType string) (int, []model.ProjectResource, error) {
	var total int
	var projectResources []model.ProjectResource
	err := db.DB.
		Model(model.ProjectResource{}).
		Where(model.ProjectResource{ProjectID: projectId, ResourceType: resourceType}).
		Count(&total).
		Find(&projectResources).
		Offset((num - 1) * size).
		Limit(size).Error
	return total, projectResources, err
}

func (p projectResourceRepository) Batch(operation string, items []model.ProjectResource) error {
	switch operation {
	case constant.BatchOperationDelete:
		var ids []string
		for _, item := range items {
			ids = append(ids, item.ID)
		}
		err := db.DB.Where("id in (?)", ids).Delete(&items).Error
		if err != nil {
			return err
		}
	case constant.BatchOperationCreate:
		tx := db.DB.Begin()
		for i, _ := range items {
			if err := tx.Model(model.ProjectResource{}).Create(&items[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
	default:
		return constant.NotSupportedBatchOperation
	}
	return nil
}
