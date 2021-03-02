package repository

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

var (
	DeleteRegionError = "DELETE_REGION_FAILED_RESOURCE"
)

type RegionRepository interface {
	Get(name string) (model.Region, error)
	List() ([]model.Region, error)
	Page(num, size int) (int, []model.Region, error)
	Save(region *model.Region) error
	Delete(name string) error
	Batch(operation string, items []model.Region) error
}

func NewRegionRepository() RegionRepository {
	return &regionRepository{}
}

type regionRepository struct {
}

func (r regionRepository) Get(name string) (model.Region, error) {
	var region model.Region
	if err := db.DB.Where("name = ?", name).First(&region).Error; err != nil {
		return region, err
	}
	return region, nil
}

func (r regionRepository) List() ([]model.Region, error) {
	var regions []model.Region
	err := db.DB.Find(&regions).Error
	return regions, err
}

func (r regionRepository) Page(num, size int) (int, []model.Region, error) {
	var total int
	var regions []model.Region
	err := db.DB.Model(&model.Region{}).Count(&total).Offset((num - 1) * size).Limit(size).Find(&regions).Error
	return total, regions, err
}

func (r regionRepository) Save(region *model.Region) error {
	if db.DB.NewRecord(region) {
		return db.DB.Create(&region).Error
	} else {
		return db.DB.Save(&region).Error
	}
}

func (r regionRepository) Delete(name string) error {
	region, err := r.Get(name)
	if err != nil {
		return err
	}
	return db.DB.Delete(&region).Error
}

func (r regionRepository) Batch(operation string, items []model.Region) error {
	switch operation {
	case constant.BatchOperationDelete:
		var regionIds []string
		for _, item := range items {
			regionIds = append(regionIds, item.ID)
		}
		var zones []model.Zone
		err := db.DB.Where("region_id in (?)", regionIds).Find(&zones).Error
		if err != nil {
			return err
		}
		if len(zones) > 0 {
			return errors.New(DeleteRegionError)
		}
		var ids []string
		for _, item := range items {
			ids = append(ids, item.ID)
		}
		err = db.DB.Where("id in (?)", ids).Delete(&items).Error
		if err != nil {
			return err
		}
	default:
		return constant.NotSupportedBatchOperation
	}
	return nil
}
