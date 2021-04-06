package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterResourceService interface {
	Page(num, size int, clusterName, resourceType string) (*page.Page, error)
}

type clusterResourceService struct {
}

func NewClusterResourceService() ClusterResourceService {
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
	if len(resourceIds) > 0 {
		switch resourceType {
		case constant.ResourceHost:
			var hosts []model.Host
			if err := db.DB.Where("id in (?)", resourceIds).Preload("Cluster").Preload("Zone").Find(&hosts).Error; err != nil {
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
			if err := db.DB.Where("id in (?)", resourceIds).Find(&result).Error; err != nil {
				return nil, err
			}
			p.Items = result
		case constant.ResourceBackupAccount:
			var result []model.BackupAccount
			if err := db.DB.Where("id in (?)", resourceIds).Find(&result).Error; err != nil {
				return nil, err
			}
			p.Items = result
		default:
			return nil, nil
		}
	}
	return &p, nil
}
