package repository

import (
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterEventRepository interface {
	List(clusterId string) ([]model.ClusterEvent, error)
	Save(event *model.ClusterEvent) error
	ListLimitOneDay(clusterId string) ([]model.ClusterEvent, error)
	ListByUidAndClusterId(uid, clusterId string) ([]model.ClusterEvent, error)
}

func NewClusterEventRepository() ClusterEventRepository {
	return &clusterEventRepository{}
}

type clusterEventRepository struct {
}

func (c clusterEventRepository) List(clusterId string) ([]model.ClusterEvent, error) {
	var events []model.ClusterEvent
	err := db.DB.Where("cluster_id = ?", clusterId).Find(&events).Error
	return events, err
}

func (c clusterEventRepository) ListLimitOneDay(clusterId string) ([]model.ClusterEvent, error) {
	var events []model.ClusterEvent
	day := time.Now().Add(time.Hour * -24)
	err := db.DB.Where("cluster_id = ? AND created_at > ?", clusterId, day).Find(&events).Error
	return events, err
}

func (c clusterEventRepository) Save(event *model.ClusterEvent) error {
	if db.DB.NewRecord(event) {
		return db.DB.Create(&event).Error
	} else {
		return db.DB.Save(&event).Error
	}
}

func (c clusterEventRepository) ListByUidAndClusterId(uid, clusterId string) ([]model.ClusterEvent, error) {
	var events []model.ClusterEvent
	day := time.Now().Add(time.Hour * -24)
	err := db.DB.Where("cluster_id = ? AND uid = ? AND created_at > ?", clusterId, uid, day).Find(&events).Error
	return events, err
}
