package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type UserNotificationConfigRepository interface {
	Get(userId string) ([]model.UserNotificationConfig, error)
	Save(notificationConfig *model.UserNotificationConfig) error
}

func NewUserNotificationConfigRepository() UserNotificationConfigRepository {
	return &userNotificationConfigRepository{}
}

type userNotificationConfigRepository struct {
}

func (u userNotificationConfigRepository) Get(userId string) ([]model.UserNotificationConfig, error) {
	var notificationConfigs []model.UserNotificationConfig
	if err := db.DB.Where(model.UserNotificationConfig{UserID: userId}).Find(&notificationConfigs).Error; err != nil {
		return notificationConfigs, err
	}
	return notificationConfigs, nil
}

func (u userNotificationConfigRepository) Save(notificationConfig *model.UserNotificationConfig) error {
	if db.DB.NewRecord(notificationConfig) {
		return db.DB.Create(&notificationConfig).Error
	} else {
		return db.DB.Model(&notificationConfig).Updates(&notificationConfig).Error
	}
}
