package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type SystemRegistryRepository interface {
	Get(arch string) (model.SystemRegistry, error)
	List() ([]model.SystemRegistry, error)
	Save(registry *model.SystemRegistry) error
}

type systemRegistryRepository struct {
}

func NewSystemRegistryRepository() SystemRegistryRepository {
	return &systemRegistryRepository{}
}

func (s systemRegistryRepository) Get(arch string) (model.SystemRegistry, error) {
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
