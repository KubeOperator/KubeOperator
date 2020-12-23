package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type IpPoolRepository interface {
	Get(name string) (model.IpPool, error)
	Save(ipPool *model.IpPool) error
	Page(num, size int) (int, []model.IpPool, error)
	Batch(operation string, items []model.IpPool) error
}

type ipPoolRepository struct {
}

func NewIpPoolRepository() IpPoolRepository {
	return &ipPoolRepository{}
}

func (i ipPoolRepository) Get(name string) (model.IpPool, error) {
	var ipPool model.IpPool
	ipPool.Name = name
	if err := db.DB.Where(ipPool).Preload("Ips").First(&ipPool).Error; err != nil {
		return ipPool, err
	}
	return ipPool, nil
}

func (i ipPoolRepository) Save(ipPool *model.IpPool) error {
	if db.DB.NewRecord(ipPool) {
		return db.DB.Create(&ipPool).Error
	} else {
		return db.DB.Save(&ipPool).Error
	}
}

func (i ipPoolRepository) Page(num, size int) (int, []model.IpPool, error) {
	var total int
	var ipPools []model.IpPool
	err := db.DB.Model(model.IpPool{}).Count(&total).Find(&ipPools).Offset((num - 1) * size).Limit(size).Error
	if err != nil {
		return 0, nil, err
	}
	return total, ipPools, nil
}

func (i ipPoolRepository) Batch(operation string, items []model.IpPool) error {

	tx := db.DB.Begin()
	switch operation {
	case constant.BatchOperationDelete:
		for i := range items {
			var ipPool model.IpPool
			if err := db.DB.Where(model.IpPool{Name: items[i].Name}).First(&ipPool).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := db.DB.Delete(&ipPool).Error; err != nil {
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
