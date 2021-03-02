package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type VmConfigRepository interface {
	Page(num, size int) (int, []model.VmConfig, error)
	List() ([]model.VmConfig, error)
	Save(item *model.VmConfig) error
	Batch(operation string, items []model.VmConfig) error
	Get(name string) (model.VmConfig, error)
}

func NewVmConfigRepository() VmConfigRepository {
	return &vmConfigRepository{}
}

type vmConfigRepository struct {
}

func (v vmConfigRepository) Page(num, size int) (int, []model.VmConfig, error) {
	var total int
	var configs []model.VmConfig
	err := db.DB.Model(&model.VmConfig{}).Count(&total).Order("cpu").Offset((num - 1) * size).Limit(size).Find(&configs).Error
	return total, configs, err
}

func (v vmConfigRepository) List() ([]model.VmConfig, error) {
	var configs []model.VmConfig
	err := db.DB.Find(&configs).Error
	return configs, err
}

func (v vmConfigRepository) Save(item *model.VmConfig) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}

func (v vmConfigRepository) Get(name string) (model.VmConfig, error) {
	var vmConfig model.VmConfig
	vmConfig.Name = name
	if err := db.DB.Where("name = ?", name).First(&vmConfig).Error; err != nil {
		return vmConfig, err
	}
	return vmConfig, nil
}

func (v vmConfigRepository) Batch(operation string, items []model.VmConfig) error {
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, item := range items {
			var config model.VmConfig
			err := db.DB.Where("name = ?", item.Name).First(&config).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			if err := db.DB.Delete(&config).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
	default:
		return constant.NotSupportedBatchOperation
	}
	return nil
}
