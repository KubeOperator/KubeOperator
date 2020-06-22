package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterMonitorRepository interface {
	Get(id string) (model.ClusterMonitor, error)
	Save(status *model.ClusterMonitor) error
	Delete(id string) error
}

func NewClusterMonitorRepository() ClusterMonitorRepository {
	return &clusterMonitorRepository{}
}

type clusterMonitorRepository struct{}

func (c clusterMonitorRepository) Get(id string) (model.ClusterMonitor, error) {
	monitor := model.ClusterMonitor{
		ID: id,
	}
	if err := db.DB.First(&monitor).Error; err != nil {
		return monitor, err
	}
	return monitor, nil
}

func (c clusterMonitorRepository) Save(monitor *model.ClusterMonitor) error {
	if db.DB.NewRecord(monitor) {
		if err := db.DB.Create(&monitor).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(&monitor).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c clusterMonitorRepository) Delete(id string) error {
	monitor := model.ClusterMonitor{ID: id}
	if err := db.DB.First(&monitor).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&monitor).Error; err != nil {
		return err
	}
	return nil
}
