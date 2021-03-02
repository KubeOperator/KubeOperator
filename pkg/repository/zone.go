package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ZoneRepository interface {
	Get(name string) (model.Zone, error)
	List() ([]model.Zone, error)
	Page(num, size int) (int, []model.Zone, error)
	Save(zone *model.Zone) error
	Delete(name string) error
	Batch(operation string, items []model.Zone) error
	ListByRegionId(id string) ([]model.Zone, error)
}

func NewZoneRepository() ZoneRepository {
	return &zoneRepository{}
}

type zoneRepository struct {
}

func (z zoneRepository) Get(name string) (model.Zone, error) {
	var zone model.Zone
	if err := db.DB.Where("name = ?", name).First(&zone).Error; err != nil {
		return zone, err
	}
	return zone, nil
}

func (z zoneRepository) ListByRegionId(id string) ([]model.Zone, error) {
	var zones []model.Zone
	err := db.DB.Where("region_id = ? AND status = ?", id, constant.Ready).Find(&zones).Error
	if err != nil {
		return zones, err
	}
	return zones, nil
}

func (z zoneRepository) List() ([]model.Zone, error) {
	var zones []model.Zone
	err := db.DB.Preload("IpPool").Find(&zones).Error
	return zones, err
}

func (z zoneRepository) Page(num, size int) (int, []model.Zone, error) {
	var total int
	var zones []model.Zone
	err := db.DB.Model(&model.Zone{}).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Preload("Region").
		Preload("IpPool").
		Preload("IpPool.Ips").
		Find(&zones).
		Error
	return total, zones, err
}

func (z zoneRepository) Save(zone *model.Zone) error {
	if db.DB.NewRecord(zone) {
		return db.DB.Create(&zone).Error
	} else {
		return db.DB.Save(&zone).Error
	}
}

func (z zoneRepository) Delete(name string) error {
	zone, err := z.Get(name)
	if err != nil {
		return err
	}
	return db.DB.Delete(&zone).Error
}

func (z zoneRepository) Batch(operation string, items []model.Zone) error {
	switch operation {
	case constant.BatchOperationDelete:
		//var zoneIds []string
		//for _, item := range items {
		//	zoneIds = append(zoneIds, item.ID)
		//}
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
