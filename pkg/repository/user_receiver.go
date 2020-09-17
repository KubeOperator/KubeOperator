package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type UserReceiverRepository interface {
	Get(userId string) (model.UserReceiver, error)
	Save(userReceiver *model.UserReceiver) error
}

func NewUserReceiverRepository() UserReceiverRepository {
	return &userReceiverRepository{}
}

type userReceiverRepository struct {
}

func (u userReceiverRepository) Get(userId string) (model.UserReceiver, error) {
	var userReceiver model.UserReceiver
	userReceiver.UserId = userId
	if err := db.DB.Where(userReceiver).First(&userReceiver).Error; err != nil {
		return userReceiver, err
	}
	return userReceiver, nil
}

func (u userReceiverRepository) Save(userReceiver *model.UserReceiver) error {
	if db.DB.NewRecord(userReceiver) {
		return db.DB.Create(&userReceiver).Error
	} else {
		return db.DB.Model(&userReceiver).Updates(&userReceiver).Error
	}
}
