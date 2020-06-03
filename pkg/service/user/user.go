package user

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
	userModel "github.com/KubeOperator/KubeOperator/pkg/model/user"
)

func Page(num, size int) (users []userModel.User, total int, err error) {
	err = db.DB.Model(userModel.User{}).Count(&total).Find(&users).Offset((num - 1) * size).Limit(size).Error
	return
}

func List() (users []userModel.User, err error) {
	err = db.DB.Model(hostModel.Host{}).Find(&users).Error
	return
}

func Get(name string) (userModel.User, error) {
	var result userModel.User
	result.Name = name
	if err := db.DB.Where(result).First(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

func Save(item *userModel.User) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}

func Delete(name string) error {
	var u userModel.User
	u.Name = name
	return db.DB.Delete(&u).Error
}

func Batch(operation string, items []userModel.User) ([]userModel.User, error) {
	var deleteItems []userModel.User
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, item := range items {
			err := db.DB.Model(userModel.User{}).First(&item).Delete(&item).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			deleteItems = append(deleteItems, item)
			tx.Commit()
		}
	default:
		return nil, constant.NotSupportedBatchOperation
	}
	return deleteItems, nil
}
