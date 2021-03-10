package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type UserNotificationConfigRepository interface {
	Get(userId string) ([]model.UserNotificationConfig, error)
	Save(notificationConfig *model.UserNotificationConfig) error
	GetByType(userId string, mType string) (*model.UserNotificationConfig, error)
}

func NewUserNotificationConfigRepository() UserNotificationConfigRepository {
	return &userNotificationConfigRepository{}
}

type userNotificationConfigRepository struct {
}

func (u userNotificationConfigRepository) Get(userId string) ([]model.UserNotificationConfig, error) {
	var notificationConfigs []model.UserNotificationConfig
	if err := db.DB.Where("user_id = ?", userId).Find(&notificationConfigs).Error; err != nil {
		return nil, err
	}
	return notificationConfigs, nil
}

func (u userNotificationConfigRepository) GetByType(userId string, mType string) (*model.UserNotificationConfig, error) {
	var notificationConfig model.UserNotificationConfig
	if err := db.DB.Where("user_id = ? AND type = ?", userId, mType).First(&notificationConfig).Error; err != nil {
		return nil, err
	}
	return &notificationConfig, nil
}
func (u userNotificationConfigRepository) Save(notificationConfig *model.UserNotificationConfig) error {
	if db.DB.NewRecord(notificationConfig) {
		return db.DB.Create(&notificationConfig).Error
	} else {
		return db.DB.Model(&notificationConfig).Updates(&notificationConfig).Error
	}
}
