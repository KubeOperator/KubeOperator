package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ProjectMemberRepository interface {
	PageByProjectId(num, size int, projectId string) (int, []model.ProjectMember, error)
	Batch(operation string, items []model.ProjectMember) error
	Create(member *model.ProjectMember) error
	ListByProjectId(projectId string) ([]model.ProjectMember, error)
}

type projectMemberRepository struct {
}

func NewProjectMemberRepository() ProjectMemberRepository {
	return &projectMemberRepository{}
}

func (p projectMemberRepository) ListByProjectId(projectId string) ([]model.ProjectMember, error) {
	var projectMembers []model.ProjectMember
	err := db.DB.Where("project_id = ?", projectId).Find(&projectMembers).Error
	return projectMembers, err
}

func (p projectMemberRepository) PageByProjectId(num, size int, projectId string) (int, []model.ProjectMember, error) {
	var total int
	var projectMembers []model.ProjectMember
	err := db.DB.Model(&model.Project{}).
		Where("project_id = ?", projectId).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Preload("User").
		Find(&projectMembers).Error
	return total, projectMembers, err
}

func (p projectMemberRepository) Create(member *model.ProjectMember) error {
	return db.DB.Create(&member).Error
}

func (p projectMemberRepository) Batch(operation string, items []model.ProjectMember) error {
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
		for i := range items {
			if err := tx.Create(&items[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()

	case constant.BatchOperationUpdate:
		tx := db.DB.Begin()
		for i := range items {
			if err := tx.Save(&items[i]).Error; err != nil {
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
