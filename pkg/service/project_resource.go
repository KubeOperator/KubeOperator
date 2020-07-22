package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type ProjectResourceService interface {
	Batch(op dto.ProjectResourceOp) error
	PageByProjectIdAndType(num, size int, projectId string, resourceType string) (page.Page, error)
	GetResources(resourceType string) (interface{}, error)
}

type projectResourceService struct {
	projectResourceRepo repository.ProjectResourceRepository
}

func NewProjectResourceService() ProjectResourceService {
	return &projectResourceService{
		projectResourceRepo: repository.NewProjectResourceRepository(),
	}
}

func (p projectResourceService) PageByProjectIdAndType(num, size int, projectId string, resourceType string) (page.Page, error) {

	var page page.Page
	var projectResourceDTOS []interface{}

	total, mos, err := p.projectResourceRepo.PageByProjectIdAndType(num, size, projectId, resourceType)
	if err != nil {
		return page, err
	}
	var resourceIds []string
	for _, mo := range mos {
		resourceIds = append(resourceIds, mo.ResourceId)
	}

	if len(resourceIds) > 0 {
		var tableName string
		switch resourceType {
		case constant.ResourceHost:
			tableName = model.Host{}.TableName()
			break
		case constant.ResourcePlan:
			tableName = model.Plan{}.TableName()
			break
		default:
			return page, err
		}
		err := db.DB.Table(tableName).Where("id in (?)", resourceIds).Find(&projectResourceDTOS).Error
		if err != nil {
			return page, err
		}
	}

	page.Total = total
	page.Items = projectResourceDTOS
	return page, err
}

func (p projectResourceService) Batch(op dto.ProjectResourceOp) error {
	var opItems []model.ProjectResource
	for _, item := range op.Items {
		opItems = append(opItems, model.ProjectResource{
			BaseModel:    common.BaseModel{},
			ID:           item.ID,
			ResourceId:   item.ResourceId,
			ResourceType: item.ResourceType,
			ProjectID:    item.ProjectID,
		})
	}
	err := p.projectResourceRepo.Batch(op.Operation, opItems)
	if err != nil {
		return err
	}
	return nil
}

func (p projectResourceService) GetResources(resourceType string) (interface{}, error) {
	var result interface{}
	var resourceIds []string
	err := db.DB.Debug().Table(model.ProjectResource{}.TableName()).Select("resource_id").Where("resource_type = ?", resourceType).Find(&resourceIds).Error
	if err != nil {
		return result, err
	}
	resourceIds = append(resourceIds, "id")

	switch resourceType {
	case constant.ResourceHost:
		var result []model.Host
		err = db.DB.Debug().Table(model.Host{}.TableName()).
			Where("id not  in (?) and cluster_id = ''", resourceIds).
			Find(&result).Error
		if err != nil {
			return result, err
		}
		return result, nil
	}
	return result, nil
}
