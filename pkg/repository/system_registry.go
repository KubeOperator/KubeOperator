package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type SystemRegistryRepository interface {
	Get(id string) (model.SystemRegistry, error)
	GetByArch(arch string) (model.SystemRegistry, error)
	List() ([]model.SystemRegistry, error)
	Save(registry *model.SystemRegistry) error
	Page(num, size int) (int, []model.SystemRegistry, error)
	Batch(operation string, items []model.SystemRegistry) error
	Delete(id string) error
}

type systemRegistryRepository struct {
}

func NewSystemRegistryRepository() SystemRegistryRepository {
	return &systemRegistryRepository{}
}

func (s systemRegistryRepository) Get(id string) (model.SystemRegistry, error) {
	var registry model.SystemRegistry
	if err := db.DB.Where("id = ?", id).First(&registry).Error; err != nil {
		return registry, err
	}
	return registry, nil
}

func (s systemRegistryRepository) GetByArch(arch string) (model.SystemRegistry, error) {
	var registry model.SystemRegistry
	if err := db.DB.Where("architecture = ?", arch).First(&registry).Error; err != nil {
		return registry, err
	}
	return registry, nil
}

func (s systemRegistryRepository) List() ([]model.SystemRegistry, error) {
	var registry []model.SystemRegistry
	if err := db.DB.Find(&registry).Error; err != nil {
		return registry, err
	}
	return registry, nil
}

func (s systemRegistryRepository) Save(registry *model.SystemRegistry) error {
	if db.DB.NewRecord(registry) {
		return db.DB.Create(&registry).Error
	} else {
		return db.DB.Save(&registry).Error
	}
}

func (s systemRegistryRepository) Page(num, size int) (int, []model.SystemRegistry, error) {
	var total int
	var registry []model.SystemRegistry
	err := db.DB.Model(&model.SystemRegistry{}).Order("architecture").Count(&total).Find(&registry).Offset((num - 1) * size).Limit(size).Error
	return total, registry, err
}

func (s systemRegistryRepository) Batch(operation string, items []model.SystemRegistry) error {
	switch operation {
	case constant.BatchOperationDelete:
		var ids []string
		for _, item := range items {
			ids = append(ids, item.ID)
		}
		err := db.DB.Where("id in (?)", ids).Delete(&items).Error
		if err != nil {
			return err
		}
	default:
		return constant.NotSupportedBatchOperation
	}
	return nil
}

func (s systemRegistryRepository) Delete(id string) error {
	registryItem, err := s.Get(id)
	if err != nil {
		return err
	}
	return db.DB.Delete(&registryItem).Error
}
