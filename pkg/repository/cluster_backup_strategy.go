package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/jinzhu/gorm"
)

type ClusterBackupStrategyRepository interface {
	Get(clusterName string) (*model.ClusterBackupStrategy, error)
	Save(clusterBackupStrategy *model.ClusterBackupStrategy) error
	List() ([]model.ClusterBackupStrategy, error)
}

type clusterBackupStrategyRepository struct {
	clusterRepository ClusterRepository
}

func NewClusterBackupStrategyRepository() ClusterBackupStrategyRepository {
	return &clusterBackupStrategyRepository{
		clusterRepository: NewClusterRepository(),
	}
}

func (c clusterBackupStrategyRepository) Get(clusterName string) (*model.ClusterBackupStrategy, error) {
	var clusterBackupStrategy model.ClusterBackupStrategy
	cluster, err := c.clusterRepository.Get(clusterName)
	if err != nil {
		return nil, err
	}
	if err := db.DB.Where("cluster_id = ?", cluster.ID).Preload("BackupAccount").First(&clusterBackupStrategy).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &clusterBackupStrategy, nil
		} else {
			return nil, err
		}
	}
	return &clusterBackupStrategy, nil
}

func (c clusterBackupStrategyRepository) Save(clusterBackupStrategy *model.ClusterBackupStrategy) error {
	if db.DB.NewRecord(clusterBackupStrategy) {
		return db.DB.Create(clusterBackupStrategy).Error
	} else {
		return db.DB.Save(&clusterBackupStrategy).Error
	}
}

func (c clusterBackupStrategyRepository) List() ([]model.ClusterBackupStrategy, error) {
	var clusterBackupStrategies []model.ClusterBackupStrategy
	err := db.DB.Order("created_at desc").Find(&clusterBackupStrategies).Error
	if err != nil {
		return nil, err
	}
	return clusterBackupStrategies, err
}
