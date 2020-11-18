package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type SystemLogRepository interface {
	Page(num, size int) (int, []model.SystemLog, error)
	List() ([]model.SystemLog, error)
	Save(item *model.SystemLog) error
}

type systemLogRepository struct {
}

func NewSystemLogRepository() SystemLogRepository {
	return &systemLogRepository{}
}

func (u systemLogRepository) Page(num, size int) (int, []model.SystemLog, error) {
	var total int
	var systemLogs []model.SystemLog
	err := db.DB.Model(model.SystemLog{}).Order("updated_at DESC").Count(&total).Find(&systemLogs).Offset((num - 1) * size).Limit(size).Error
	return total, systemLogs, err
}

func (u systemLogRepository) Save(item *model.SystemLog) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}

func (u systemLogRepository) List() ([]model.SystemLog, error) {
	var logs []model.SystemLog
	err := db.DB.Model(model.SystemLog{}).Order("updated_at DESC").Find(&logs).Error
	return logs, err
}
