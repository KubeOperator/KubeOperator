package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"time"
)

type ClusterLogRepository interface {
	Save(clusterName string, log *model.ClusterLog) error
	List(clusterName string) ([]model.ClusterLog, error)
	GetRunningLogWithClusterNameAndType(clusterName string, logType string) (model.ClusterLog, error)
}

func NewClusterLogRepository() ClusterLogRepository {
	return &clusterLogRepository{}
}

type clusterLogRepository struct {
}

func (c *clusterLogRepository) Save(clusterName string, log *model.ClusterLog) error {
	var cluster model.Cluster
	if err := db.DB.Where(&model.Cluster{Name: clusterName}).First(&cluster).Error; err != nil {
		return err
	}
	log.ClusterID = cluster.ID
	if db.DB.NewRecord(log) {
		return db.DB.Create(log).Error
	} else {
		return db.DB.Save(log).Error
	}
}

func (c *clusterLogRepository) List(clusterName string) ([]model.ClusterLog, error) {
	var cluster model.Cluster
	if err := db.DB.Where(&model.Cluster{Name: clusterName}).First(&cluster).Error; err != nil {
		return nil, err
	}
	var items []model.ClusterLog
	if err := db.DB.Where(&model.ClusterLog{ClusterID: cluster.ID}).
		Order("created_at desc").
		Find(&items).
		Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (c *clusterLogRepository) GetRunningLogWithClusterNameAndType(clusterName string, logType string) (model.ClusterLog, error) {
	var item model.ClusterLog
	var cluster model.Cluster
	if err := db.DB.Where(&model.Cluster{Name: clusterName}).First(&cluster).Error; err != nil {
		return item, err
	}
	now := time.Now()
	h, _ := time.ParseDuration("-12h")
	halfDayAgo := now.Add(h)
	if err := db.DB.Where("cluster_id = ? AND type = ? AND status = ? AND created_at BETWEEN ? AND ?", cluster.ID, logType, constant.ClusterRunning, halfDayAgo, now).
		Find(&item).
		Error; err != nil {
		return item, err
	}
	return item, nil
}
