package service

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/jinzhu/gorm"
)

type ProjectResourceService interface {
	Batch(op dto.ProjectResourceOp) error
	Page(num, size int, projectName string, resourceType string) (*page.Page, error)
	GetResources(resourceType, projectName string) (interface{}, error)
}

type projectResourceService struct {
	projectResourceRepo repository.ProjectResourceRepository
	projectRepo         repository.ProjectRepository
}

func NewProjectResourceService() ProjectResourceService {
	return &projectResourceService{
		projectResourceRepo: repository.NewProjectResourceRepository(),
		projectRepo:         repository.NewProjectRepository(),
	}
}

func (p projectResourceService) Page(num, size int, projectName string, resourceType string) (*page.Page, error) {
	var page page.Page
	pj, err := p.projectRepo.Get(projectName)
	if err != nil {
		return nil, err
	}
	total, mos, err := p.projectResourceRepo.PageByProjectIdAndType(num, size, pj.ID, resourceType)
	if err != nil {
		return nil, err
	}
	var resourceIds []string
	for _, mo := range mos {
		resourceIds = append(resourceIds, mo.ResourceID)
	}

	if len(resourceIds) > 0 {
		switch resourceType {
		case constant.ResourceHost:
			var hosts []model.Host
			err = db.DB.Model(model.Host{}).Where("id in (?)", resourceIds).Preload("Cluster").Preload("Zone").Find(&hosts).Error
			if err != nil {
				return nil, err
			}

			var result []dto.Host
			for _, mo := range hosts {
				hostDTO := dto.Host{
					Host:        mo,
					ClusterName: mo.Cluster.Name,
					ZoneName:    mo.Zone.Name,
				}
				result = append(result, hostDTO)
			}
			page.Items = result
		case constant.ResourcePlan:
			var result []model.Plan
			err = db.DB.Model(model.Plan{}).Where("id in (?)", resourceIds).Find(&result).Error
			if err != nil {
				return nil, err
			}
			page.Items = result
		case constant.ResourceBackupAccount:
			var result []model.BackupAccount
			err = db.DB.Model(model.BackupAccount{}).Where("id in (?)", resourceIds).Find(&result).Error
			if err != nil {
				return nil, err
			}
			page.Items = result
		default:
			return nil, err
		}

		page.Total = total
	}

	return &page, err
}

func (p projectResourceService) Batch(op dto.ProjectResourceOp) error {
	var opItems []model.ProjectResource
	for _, item := range op.Items {

		var resourceId string
		switch item.ResourceType {
		case constant.ResourceHost:
			host, err := NewHostService().Get(item.ResourceName)
			if err != nil {
				return err
			}
			resourceId = host.ID
			if host.ClusterID != "" {
				return errors.New("DELETE_HOST_FAILED_BY_CLUSTER")
			}
		case constant.ResourcePlan:
			plan, err := NewPlanService().Get(item.ResourceName)
			if err != nil {
				return err
			}
			resourceId = plan.ID
		case constant.ResourceBackupAccount:
			plan, err := NewBackupAccountService().Get(item.ResourceName)
			if err != nil {
				return err
			}
			resourceId = plan.ID
		}

		var itemId string
		if op.Operation == constant.BatchOperationDelete {
			var p model.ProjectResource
			err := db.DB.Model(model.ProjectResource{}).
				Where("project_id = ? AND resource_type = ? AND resource_id = ?", item.ProjectID, item.ResourceType, resourceId).First(&p).Error
			if err != nil {
				return err
			}
			itemId = p.ID

			if item.ResourceType == constant.ResourceBackupAccount {
				var clusterResources []model.ProjectResource
				err = db.DB.Where(model.ProjectResource{ProjectID: item.ProjectID, ResourceType: constant.ResourceCluster}).Find(&clusterResources).Error
				if err != nil && !gorm.IsRecordNotFoundError(err) {
					return err
				}
				if len(clusterResources) > 0 {
					for _, clusterResource := range clusterResources {
						var backupStrategy model.ClusterBackupStrategy
						err = db.DB.Where(model.ClusterBackupStrategy{BackupAccountID: resourceId, ClusterID: clusterResource.ResourceID}).First(&backupStrategy).Error
						if err != nil && !gorm.IsRecordNotFoundError(err) {
							return err
						}
						if backupStrategy.ID != "" {
							var backupFiles []model.ClusterBackupFile
							err = db.DB.Where(model.ClusterBackupFile{ClusterBackupStrategyID: backupStrategy.ID, ClusterID: clusterResource.ResourceID}).Find(&backupFiles).Error
							if err != nil && !gorm.IsRecordNotFoundError(err) {
								return err
							}
							if len(backupFiles) > 0 {
								return errors.New("DELETE_FAILED_BY_BACKUP_FILE")
							}
						}
					}
				}
			}

			if item.ResourceType == constant.ResourceHost {
				var clusterResources []model.ProjectResource
				err = db.DB.Where(model.ProjectResource{ResourceID: item.ResourceID, ResourceType: constant.ResourceHost}).Find(&clusterResources).Error
				if err != nil && !gorm.IsRecordNotFoundError(err) {
					return err
				}
				if len(clusterResources) > 0 {
					continue
				}
			}
		}

		opItems = append(opItems, model.ProjectResource{
			BaseModel:    common.BaseModel{},
			ID:           itemId,
			ResourceID:   resourceId,
			ResourceType: item.ResourceType,
			ProjectID:    item.ProjectID,
		})
	}
	return p.projectResourceRepo.Batch(op.Operation, opItems)
}

func (p projectResourceService) GetResources(resourceType, projectName string) (interface{}, error) {
	var result interface{}
	var projectResources []model.ProjectResource
	var resourceIds []string
	if resourceType == constant.ResourcePlan || resourceType == constant.ResourceBackupAccount {
		project, err := p.projectRepo.Get(projectName)
		if err != nil {
			return nil, err
		}
		err = db.DB.Model(model.ProjectResource{}).Select("resource_id").Where(model.ProjectResource{ProjectID: project.ID, ResourceType: resourceType}).Find(&projectResources).Error
		if err != nil {
			return nil, err
		}
	}
	if resourceType == constant.ResourceHost {
		err := db.DB.Model(model.ProjectResource{}).Select("resource_id").Where(model.ProjectResource{ResourceType: resourceType}).Find(&projectResources).Error
		if err != nil {
			return nil, err
		}
	}
	for _, pr := range projectResources {
		resourceIds = append(resourceIds, pr.ResourceID)
	}
	if len(resourceIds) == 0 {
		resourceIds = append(resourceIds, "1")
	}

	switch resourceType {
	case constant.ResourceHost:
		var result []model.Host
		err := db.DB.Model(model.Host{}).
			Where("id not  in (?) and cluster_id = ''", resourceIds).
			Find(&result).Error
		if err != nil {
			return nil, err
		}
		return result, nil

	case constant.ResourcePlan:
		var result []model.Plan
		resourceIds = append(resourceIds, "1")
		err := db.DB.Model(model.Plan{}).Where("id not in (?)", resourceIds).Preload("Zones").Preload("Region").Find(&result).Error
		if err != nil {
			return nil, err
		}
		return result, nil

	case constant.ResourceBackupAccount:
		var result []model.BackupAccount
		resourceIds = append(resourceIds, "1")
		err := db.DB.Model(model.BackupAccount{}).Where("id not in (?)", resourceIds).Find(&result).Error
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return result, nil
}
