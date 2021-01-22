package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"time"
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
	err := db.DB.Where(&model.ClusterEvent{ClusterID: clusterId}).Find(&events).Error
	return events, err
}
func (c clusterEventRepository) ListLimitOneDay(clusterId string) ([]model.ClusterEvent, error) {
	var events []model.ClusterEvent
	day := time.Now().Add(time.Hour * -24).Format("2006-01-02 15:04:05")
	err := db.DB.Where(&model.ClusterEvent{ClusterID: clusterId}).
		Where("created_at > (?)", day).
		Find(&events).Error
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
	day := time.Now().Add(time.Hour * -24).Format("2006-01-02 15:04:05")
	err := db.DB.Where(&model.ClusterEvent{ClusterID: clusterId, UID: uid}).
		Where("created_at > (?)", day).
		Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, err
}
