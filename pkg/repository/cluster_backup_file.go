package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterBackupFileRepository interface {
	Page(num, size int, clusterId string) (int, []model.ClusterBackupFile, error)
	Save(file *model.ClusterBackupFile) error
	Batch(operation string, items []model.ClusterBackupFile) error
	Get(name string) (model.ClusterBackupFile, error)
	Delete(name string) error
}

type clusterBackupFileRepository struct {
}

func NewClusterBackupFileRepository() ClusterBackupFileRepository {
	return &clusterBackupFileRepository{}
}

func (c clusterBackupFileRepository) Page(num, size int, clusterId string) (int, []model.ClusterBackupFile, error) {
	var total int
	var files []model.ClusterBackupFile
	err := db.DB.Model(&model.ClusterBackupFile{}).
		Where("cluster_id = ?", clusterId).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Find(&files).Error
	return total, files, err
}

func (c clusterBackupFileRepository) Get(name string) (model.ClusterBackupFile, error) {
	var file model.ClusterBackupFile
	if err := db.DB.Where("name = ?", name).
		Preload("ClusterBackupStrategy").
		Preload("ClusterBackupStrategy.BackupAccount").
		Find(&file).Error; err != nil {
		return file, err
	}
	return file, nil
}

func (c clusterBackupFileRepository) Save(file *model.ClusterBackupFile) error {
	if db.DB.NewRecord(file) {
		return db.DB.Create(&file).Error
	} else {
		return db.DB.Updates(&file).Error
	}
}

func (c clusterBackupFileRepository) Batch(operation string, items []model.ClusterBackupFile) error {

	tx := db.DB.Begin()
	switch operation {
	case constant.BatchOperationDelete:
		for i := range items {

			var file model.ClusterBackupFile
			if err := db.DB.Where("name = ?", items[i].Name).First(&file).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := db.DB.Delete(&file).Error; err != nil {
				tx.Rollback()
				return err
			}

		}
	default:
		return constant.NotSupportedBatchOperation
	}
	tx.Commit()
	return nil
}

func (c clusterBackupFileRepository) Delete(name string) error {
	host, err := c.Get(name)
	if err != nil {
		return err
	}
	return db.DB.Delete(&host).Error
}
