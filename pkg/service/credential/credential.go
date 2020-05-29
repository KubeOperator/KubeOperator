package credential

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	credentialModel "github.com/KubeOperator/KubeOperator/pkg/model/credential"
)

func Page(num, size int) (credentials []credentialModel.Credential, total int, err error) {
	err = db.DB.Model(credentialModel.Credential{}).
		Count(&total).
		Find(&credentials).
		Offset((num - 1) * size).
		Limit(size).
		Error
	return
}

func List() (credentials []credentialModel.Credential, err error) {
	err = db.DB.Model(credentialModel.Credential{}).Find(&credentials).Error
	return
}

func Get(name string) (credentialModel.Credential, error) {
	var result credentialModel.Credential
	result.Name = name
	if err := db.DB.Where(result).First(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

func GetById(id string) (credentialModel.Credential, error) {
	var result credentialModel.Credential
	result.ID = id
	if err := db.DB.Where(result).First(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

func Save(item *credentialModel.Credential) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}

func Delete(name string) error {
	var c credentialModel.Credential
	c.Name = name
	return db.DB.Delete(&c).Error
}

func Batch(operation string, items []credentialModel.Credential) ([]credentialModel.Credential, error) {
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, item := range items {
			err := db.DB.Model(credentialModel.Credential{}).First(&item).Delete(&item).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		tx.Commit()
	default:
		return nil, constant.NotSupportedBatchOperation
	}
	return items, nil
}
