package host

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
)

func Page(num, size int) (host []hostModel.Host, total int, err error) {
	err = db.DB.Model(hostModel.Host{}).
		Find(&host).
		Offset((num - 1) * size).
		Limit(size).
		Count(&total).
		Error
	return
}

func List() (host []hostModel.Host, err error) {
	err = db.DB.Model(hostModel.Host{}).Find(&host).Error
	return
}

func Get(name string) (*hostModel.Host, error) {
	var result hostModel.Host
	err := db.DB.Model(hostModel.Host{}).Where(&result).First(&result).Error
	return &result, err
}

func Save(item *hostModel.Host) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}

func Delete(name string) error {
	var h hostModel.Host
	h.Name = name
	return db.DB.Delete(&h).Error
}

func Batch(operation string, items []hostModel.Host) ([]hostModel.Host, error) {
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, item := range items {
			err := db.DB.Model(hostModel.Host{}).Delete(&item).Error
			if err != nil {
				tx.Rollback()
			}
		}
		tx.Commit()
	default:
		return nil, constant.NotSupportedBatchOperation
	}
	return items, nil
}
