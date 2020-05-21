package credential

import (
	"ko3-gin/pkg/constant"
	"ko3-gin/pkg/db"
	credentialModel "ko3-gin/pkg/model/credential"
)

func Page(num, size int) (credentials []credentialModel.Credential, total int, err error) {
	err = db.DB.Model(credentialModel.Credential{}).
		Find(&credentials).
		Offset((num - 1) * size).
		Limit(size).
		Count(&total).
		Error
	return
}

func List() (credentials []credentialModel.Credential, err error) {
	err = db.DB.Model(credentialModel.Credential{}).Find(&credentials).Error
	return
}

func Get(name string) (*credentialModel.Credential, error) {
	var result credentialModel.Credential
	err := db.DB.Model(credentialModel.Credential{}).Where(&result).First(&result).Error
	return &result, err
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
			err := db.DB.Model(credentialModel.Credential{}).Delete(&item).Error
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
