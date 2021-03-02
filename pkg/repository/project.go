package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ProjectRepository interface {
	Get(name string) (model.Project, error)
	List() ([]model.Project, error)
	Page(num, size int, userId string) (int, []model.Project, error)
	Save(project *model.Project) error
	Batch(operation string, items []model.Project) error
	Delete(name string) error
}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

type projectRepository struct {
}

func (p projectRepository) Get(name string) (model.Project, error) {
	var project model.Project
	if err := db.DB.Where("name = ?", name).First(&project).Error; err != nil {
		return project, err
	}
	return project, nil
}

func (p projectRepository) List() ([]model.Project, error) {
	var projects []model.Project
	err := db.DB.Find(&projects).Error
	return projects, err
}

func (p projectRepository) Page(num, size int, userId string) (int, []model.Project, error) {
	var total int
	var projects []model.Project
	if userId == "" {
		err := db.DB.Model(&model.Project{}).Count(&total).Order("name").Offset((num - 1) * size).Limit(size).Find(&projects).Error
		return total, projects, err
	} else {
		var projectResources []model.ProjectMember
		err := db.DB.Where("user_id = ?", userId).Find(&projectResources).Error
		if err != nil {
			return total, nil, err
		}
		var projectIds []string
		for _, pm := range projectResources {
			projectIds = append(projectIds, pm.ProjectID)
		}
		err = db.DB.Model(&model.Project{}).Order("name").Where("id in (?)", projectIds).Count(&total).Offset((num - 1) * size).Limit(size).Find(&projects).Error
		if err != nil {
			return total, nil, err
		}
		return total, projects, err
	}
}

func (p projectRepository) Save(project *model.Project) error {
	if db.DB.NewRecord(project) {
		return db.DB.Create(&project).Error
	} else {
		return db.DB.Save(&project).Error
	}
}

func (p projectRepository) Batch(operation string, items []model.Project) error {
	switch operation {
	case constant.BatchOperationDelete:
		var ids []string
		for _, item := range items {
			ids = append(ids, item.ID)
		}
		tx := db.DB.Begin()
		err := tx.Where("id in (?)", ids).Delete(&items).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Where("project_id in (?)", ids).Delete(&model.ProjectResource{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	default:
		return constant.NotSupportedBatchOperation
	}
	return nil
}

func (p projectRepository) Delete(name string) error {
	project, err := p.Get(name)
	if err != nil {
		return err
	}
	tx := db.DB.Begin()
	err = tx.Delete(&project).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Where("project_id = ?", project.ID).Delete(&model.ProjectResource{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}
