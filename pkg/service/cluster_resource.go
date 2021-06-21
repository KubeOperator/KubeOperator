package service

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/errorf"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/jinzhu/gorm"
)

type ClusterResourceService interface {
	Page(num, size int, clusterName, resourceType string) (*page.Page, error)
	List(clusterName string, resourceType string) (interface{}, error)
	Create(clusterName string, request dto.ClusterResourceCreate) ([]dto.ClusterResource, error)
	GetResources(resourceType, projectName, clusterName string) (interface{}, error)
	Delete(name, resourceType, clusterName string) error
}

type clusterResourceService struct {
}

func NewClusterResourceService() *clusterResourceService {
	return &clusterResourceService{}
}

func (c clusterResourceService) Page(num, size int, clusterName, resourceType string) (*page.Page, error) {
	var (
		p                page.Page
		cluster          model.Cluster
		clusterResources []model.ClusterResource
		resourceIds      []string
	)
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Model(&model.ClusterResource{}).
		Where("cluster_id = ? AND resource_type= ?", cluster.ID, resourceType).Count(&p.Total).
		Offset((num - 1) * size).
		Limit(size).
		Find(&clusterResources).Error; err != nil {
		return nil, err
	}
	for _, mo := range clusterResources {
		resourceIds = append(resourceIds, mo.ResourceID)
	}

	if len(resourceIds) == 0 {
		resourceIds = append(resourceIds, "1")
	}

	switch resourceType {
	case constant.ResourceHost:
		var hosts []model.Host
		if err := db.DB.Where("id in (?)", resourceIds).Preload("Cluster").Preload("Zone").Order("created_at desc").Find(&hosts).Error; err != nil {
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
		p.Items = result
	case constant.ResourcePlan:
		var result []model.Plan
		if err := db.DB.Where("id in (?)", resourceIds).Order("created_at desc").Find(&result).Error; err != nil {
			return nil, err
		}
		p.Items = result
	case constant.ResourceBackupAccount:
		var result []model.BackupAccount
		if err := db.DB.Where("id in (?)", resourceIds).Order("created_at desc").Find(&result).Error; err != nil {
			return nil, err
		}
		p.Items = result
	default:
		return nil, nil
	}
	return &p, nil
}

func (c clusterResourceService) List(clusterName string, resourceType string) (interface{}, error) {
	var (
		cluster          model.Cluster
		clusterResources []model.ClusterResource
		resourceIds      []string
		resources        interface{}
	)
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Model(&model.ClusterResource{}).
		Where("cluster_id = ? AND resource_type= ?", cluster.ID, resourceType).
		Find(&clusterResources).Error; err != nil {
		return nil, err
	}
	for _, mo := range clusterResources {
		resourceIds = append(resourceIds, mo.ResourceID)
	}

	if len(resourceIds) == 0 {
		resourceIds = append(resourceIds, "1")
	}

	switch resourceType {
	case constant.ResourceHost:
		var hosts []model.Host
		if err := db.DB.Where("id in (?)", resourceIds).Preload("Cluster").Preload("Zone").Order("created_at desc").Find(&hosts).Error; err != nil {
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
		resources = result
	case constant.ResourcePlan:
		var result []model.Plan
		if err := db.DB.Where("id in (?)", resourceIds).Order("created_at desc").Find(&result).Error; err != nil {
			return nil, err
		}
		resources = result
	case constant.ResourceBackupAccount:
		var result []model.BackupAccount
		if err := db.DB.Where("id in (?)", resourceIds).Order("created_at desc").Find(&result).Error; err != nil {
			return nil, err
		}
		resources = result
	default:
		return nil, nil
	}
	return resources, nil
}

func (c clusterResourceService) Create(clusterName string, request dto.ClusterResourceCreate) ([]dto.ClusterResource, error) {

	if err := createCheck(clusterName, request); err != nil {
		return nil, err
	}
	var (
		cluster model.Cluster
		errs    errorf.CErrFs
		result  []dto.ClusterResource
	)
	if err := db.DB.Model(&model.Cluster{}).Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return nil, err
	}
	for _, name := range request.Names {
		var resourceId string
		if request.ResourceType == constant.ResourceHost {
			var host model.Host
			if err := db.DB.Model(model.Host{}).Where("name = ?", name).Find(&host).Error; err != nil {
				errs = errs.Add(errorf.New("HOST_IS_NOT_FOUND", name))
				continue
			} else {
				resourceId = host.ID
			}
		} else if request.ResourceType == constant.ResourcePlan {
			var plan model.Plan
			if err := db.DB.Model(model.Plan{}).Where("name = ?", name).Find(&plan).Error; err != nil {
				errs = errs.Add(errorf.New("PLAN_IS_NOT_FOUND", name))
				continue
			} else {
				resourceId = plan.ID
			}
		} else if request.ResourceType == constant.ResourceBackupAccount {
			var backupAccount model.BackupAccount
			if err := db.DB.Model(model.BackupAccount{}).Where("name = ?", name).Find(&backupAccount).Error; err != nil {
				errs = errs.Add(errorf.New("BACKUP_ACCOUNT_IS_NOT_FOUND", name))
				continue
			} else {
				resourceId = backupAccount.ID
			}
		}
		if resourceId != "" {
			var oldCr model.ClusterResource
			if err := db.DB.Model(model.ClusterResource{}).Where("resource_id = ? AND cluster_id = ?", resourceId, cluster.ID).Find(&oldCr).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
				errs = errs.Add(errorf.New(err.Error()))
				continue
			}
			if oldCr.ID != "" {
				errs = errs.Add(errorf.New("RESOURCE_IS_ADDED", name))
				continue
			}
			cr := model.ClusterResource{
				ResourceID:   resourceId,
				ClusterID:    cluster.ID,
				ResourceType: request.ResourceType,
			}
			if err := db.DB.Create(&cr).Error; err != nil {
				errs = errs.Add(errorf.New(err.Error()))
			}
			result = append(result, dto.ClusterResource{
				ClusterResource: cr,
				ResourceName:    name,
			})
		}
	}
	if len(errs) > 0 {
		return result, errs
	} else {
		return result, nil
	}
}

func (c clusterResourceService) Delete(name, resourceType, clusterName string) error {

	if err := deleteCheck(name, resourceType, clusterName); err != nil {
		return err
	}

	var (
		cluster    model.Cluster
		cr         model.ClusterResource
		resourceId string
	)
	if resourceType == constant.ResourceHost {
		var host model.Host
		if err := db.DB.Model(model.Host{}).Where("name = ?", name).Find(&host).Error; err != nil {
			return err
		} else {
			if host.ClusterID != "" {
				return errors.New("DELETE_HOST_FAILED")
			}
			resourceId = host.ID
		}
	} else if resourceType == constant.ResourceBackupAccount {
		var backupAccount model.BackupAccount
		if err := db.DB.Model(model.BackupAccount{}).Where("name = ?", name).Find(&backupAccount).Error; err != nil {
			return err
		} else {
			resourceId = backupAccount.ID
		}
	}
	if err := db.DB.Model(model.Cluster{}).Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return err
	}
	if err := db.DB.Model(model.ProjectResource{}).Where("cluster_id = ? AND resource_id = ?", cluster.ID, resourceId).Find(&cr).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&cr).Error; err != nil {
		return err
	}
	return nil
}

func (c clusterResourceService) GetResources(resourceType, projectName, clusterName string) (interface{}, error) {
	var (
		result  interface{}
		project model.Project
		cluster model.Cluster
	)

	if err := db.DB.Model(&model.Project{}).Where("name = ?", projectName).Find(&project).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Model(&model.Cluster{}).Where("name = ?", clusterName).Find(&cluster).Error; err != nil {
		return nil, err
	}

	if resourceType == constant.ResourceHost {
		var hosts []model.Host
		if err := db.DB.Raw("SELECT * FROM ko_host WHERE cluster_id='' AND id IN (SELECT resource_id FROM ko_project_resource WHERE resource_type='HOST' AND project_id= ? AND resource_id NOT IN (SELECT resource_id FROM ko_cluster_resource WHERE resource_type='HOST'))", project.ID).Scan(&hosts).Error; err != nil {
			return nil, err
		}
		result = hosts
	}

	if resourceType == constant.ResourceBackupAccount {
		var backupAccounts []model.BackupAccount
		if err := db.DB.Raw("SELECT * FROM ko_backup_account WHERE id in (SELECT resource_id FROM ko_project_resource WHERE resource_type='BACKUP_ACCOUNT' AND project_id= ? AND resource_id NOT IN (SELECT resource_id FROM ko_cluster_resource WHERE resource_type='BACKUP_ACCOUNT' AND cluster_id  =?) )", project.ID, cluster.ID).Scan(&backupAccounts).Error; err != nil {
			return nil, err
		}
		result = backupAccounts
	}

	return result, nil
}

func createCheck(clusterName string, request dto.ClusterResourceCreate) error {

	resourceTypes := []string{constant.ResourceHost, constant.ResourceBackupAccount}
	result := false
	for _, resourceType := range resourceTypes {
		if resourceType == request.ResourceType {
			result = true
			break
		}
	}
	if !result {
		return errors.New("RESOURCE_TYPE_ERROR")
	}

	var cluster model.Cluster
	if err := db.DB.Model(&model.Cluster{}).Preload("Spec").
		Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return err
	}
	if cluster.Spec.Provider == constant.ClusterProviderPlan && request.ResourceType == constant.ResourceHost {
		return errors.New("CLUSTER_PROVIDER_ERROR")
	}

	return nil
}

func deleteCheck(resourceName, resourceType, clusterName string) error {

	resourceTypes := []string{constant.ResourceHost, constant.ResourceBackupAccount}
	result := false
	for _, res := range resourceTypes {
		if res == resourceType {
			result = true
			break
		}
	}
	if !result {
		return errors.New("RESOURCE_TYPE_ERROR")
	}

	var cluster model.Cluster
	if err := db.DB.Model(&model.Cluster{}).Preload("Spec").
		Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return err
	}
	if cluster.Spec.Provider == constant.ClusterProviderPlan && resourceType == constant.ResourceHost {
		return errors.New("CLUSTER_PROVIDER_ERROR")
	}

	if resourceType == constant.ResourceBackupAccount {
		var backupAccount model.BackupAccount
		if err := db.DB.Model(&model.BackupAccount{}).Where("name = ?", resourceName).Find(&backupAccount).Error; err != nil {
			return err
		}
		var backupStrategy model.ClusterBackupStrategy
		if err := db.DB.Where("cluster_id = ? AND backup_account_id = ? AND status = ?", cluster.ID, backupAccount.ID, "ENABLE").First(&backupStrategy).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
			return err
		}
		if backupStrategy.ID != "" {
			return errors.New("DELETE_BACKUP_ACCOUNT_FAILED_BY_CLUSTER")
		}
	}

	return nil
}
