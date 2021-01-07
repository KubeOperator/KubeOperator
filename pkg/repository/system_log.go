package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type SystemLogRepository interface {
	Page(num, size int, queryOption, queryInfo string) (int, []model.SystemLog, error)
	Save(item *model.SystemLog) error
}

type systemLogRepository struct {
}

func NewSystemLogRepository() SystemLogRepository {
	return &systemLogRepository{}
}

func (u systemLogRepository) Page(num, size int, queryOption, queryInfo string) (total int, systemLogs []model.SystemLog, err error) {
	if len(queryInfo) != 0 {
		switch queryOption {
		case "name":
			err = db.DB.Model(model.SystemLog{}).Where("name LIKE ?", "%"+queryInfo+"%").Order("updated_at DESC").Count(&total).Find(&systemLogs).Offset((num - 1) * size).Limit(size).Error
		case "operationUnit":
			err = db.DB.Model(model.SystemLog{}).Where("operation_unit LIKE ?", "%"+queryInfo+"%").Order("updated_at DESC").Count(&total).Find(&systemLogs).Offset((num - 1) * size).Limit(size).Error
		case "operation":
			err = db.DB.Model(model.SystemLog{}).Where("operation LIKE ?", "%"+queryInfo+"%").Order("updated_at DESC").Count(&total).Find(&systemLogs).Offset((num - 1) * size).Limit(size).Error
		}
	} else {
		err = db.DB.Model(model.SystemLog{}).Order("updated_at DESC").Count(&total).Find(&systemLogs).Offset((num - 1) * size).Limit(size).Error
	}
	return
}

func (u systemLogRepository) Save(item *model.SystemLog) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}
