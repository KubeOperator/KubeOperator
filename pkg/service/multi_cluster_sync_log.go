package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type MultiClusterSyncLogService interface {
	Page(name string, num, size int) (*page.Page, error)
	Detail(name, logId string) (*dto.MultiClusterSyncLogDetail, error)
}
type multiClusterSyncLogService struct {
}

func (m multiClusterSyncLogService) Page(name string, num, size int) (*page.Page, error) {
	var p page.Page
	var repository model.MultiClusterRepository
	if err := db.DB.
		Where(model.MultiClusterRepository{Name: name}).
		First(&repository).Error; err != nil {
		return nil, err
	}
	var syncLogMos []model.MultiClusterSyncLog
	if err := db.DB.Model(model.MultiClusterSyncLog{}).Where(model.MultiClusterSyncLog{MultiClusterRepositoryID: repository.ID}).
		Count(&p.Total).
		Offset((num - 1) * size).
		Limit(size).
		Order("created_at desc").
		Find(&syncLogMos).Error; err != nil {
		return nil, err
	}
	var items []dto.MultiClusterSyncLog
	for _, mo := range syncLogMos {
		items = append(items, dto.MultiClusterSyncLog{MultiClusterSyncLog: mo})
	}
	p.Items = items
	return &p, nil
}

func (m multiClusterSyncLogService) Detail(name, logId string) (*dto.MultiClusterSyncLogDetail, error) {
	var repository model.MultiClusterRepository
	if err := db.DB.
		Where(model.MultiClusterRepository{Name: name}).
		First(&repository).Error; err != nil {
		return nil, err
	}
	var clusterLogs []model.MultiClusterSyncClusterLog
	var syncLog model.MultiClusterSyncLog
	if err := db.DB.Where(model.MultiClusterSyncLog{ID: logId}).First(&syncLog).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Where(model.MultiClusterSyncClusterLog{MultiClusterSyncLogID: logId}).Find(&clusterLogs).Error; err != nil {
		return nil, err
	}
	var item dto.MultiClusterSyncLogDetail
	item.MultiClusterSyncLog = syncLog
	for _, cl := range clusterLogs {
		var clDto dto.MultiClusterSyncClusterLog
		clDto.MultiClusterSyncClusterLog = cl
		var resourceLogs []model.MultiClusterSyncClusterResourceLog
		if err := db.DB.Where(model.MultiClusterSyncClusterResourceLog{MultiClusterSyncClusterLogID: cl.ID}).Find(&resourceLogs).Error; err != nil {
			return nil, err
		}
		var cluster model.Cluster
		if err := db.DB.Where(model.Cluster{ID: cl.ClusterID}).First(&cluster).Error; err != nil {
			return nil, err
		}
		clDto.ClusterName = cluster.Name
		clDto.MultiClusterSyncClusterResourceLogs = append(clDto.MultiClusterSyncClusterResourceLogs, resourceLogs...)
		item.MultiClusterSyncClusterLogs = append(item.MultiClusterSyncClusterLogs, clDto)
	}
	return &item, nil
}

func NewMultiClusterSyncLogService() MultiClusterSyncLogService {
	return &multiClusterSyncLogService{}
}
