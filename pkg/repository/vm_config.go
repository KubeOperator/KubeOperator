package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type VmConfigRepository interface {
	Page(num, size int) (int, []model.VmConfig, error)
	List() ([]model.VmConfig, error)
}

func NewVmConfigRepository() VmConfigRepository {
	return &vmConfigRepository{}
}

type vmConfigRepository struct {
}

func (v vmConfigRepository) Page(num, size int) (int, []model.VmConfig, error) {
	var total int
	var configs []model.VmConfig
	err := db.DB.Model(model.VmConfig{}).
		Count(&total).Find(&configs).
		Offset((num - 1) * size).
		Limit(size).
		Error
	return total, configs, err
}

func (v vmConfigRepository) List() ([]model.VmConfig, error) {
	var configs []model.VmConfig
	err := db.DB.Model(model.VmConfig{}).Find(&configs).Error
	return configs, err
}
