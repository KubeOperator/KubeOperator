package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type SystemLogService interface {
	Create(creation dto.SystemLogCreate) error
	Page(num, size int, queryOption, queryInfo string) (page.Page, error)
}

type systemLogService struct {
	systemLogRepo repository.SystemLogRepository
}

func NewSystemLogService() SystemLogService {
	return &systemLogService{
		systemLogRepo: repository.NewSystemLogRepository(),
	}
}

func (s systemLogService) Create(creation dto.SystemLogCreate) error {
	log := model.SystemLog{
		Name:          creation.Name,
		Operation:     creation.Operation,
		OperationInfo: creation.OperationInfo,
	}

	if db.DB.NewRecord(log) {
		return db.DB.Create(&log).Error
	} else {
		return db.DB.Save(&log).Error
	}
}

func (u systemLogService) Page(num, size int, queryOption, queryInfo string) (page.Page, error) {
	var (
		page      page.Page
		logsOfDB  []model.SystemLog
		logsOfDTO []dto.SystemLog
		total     int
		err       error
	)

	if len(queryInfo) != 0 {
		switch queryOption {
		case "name":
			err = db.DB.Model(model.SystemLog{}).Where("name LIKE ?", "%"+queryInfo+"%").Order("updated_at DESC").Count(&total).Offset((num - 1) * size).Limit(size).Find(&logsOfDB).Error
		case "operation":
			err = db.DB.Model(model.SystemLog{}).Where("operation LIKE ?", "%"+queryInfo+"%").Order("updated_at DESC").Count(&total).Offset((num - 1) * size).Limit(size).Find(&logsOfDB).Error
		case "operationInfo":
			err = db.DB.Model(model.SystemLog{}).Where("operation_info LIKE ?", "%"+queryInfo+"%").Order("updated_at DESC").Count(&total).Offset((num - 1) * size).Limit(size).Find(&logsOfDB).Error
		}
	} else {
		err = db.DB.Model(model.SystemLog{}).Order("updated_at DESC").Count(&total).Offset((num - 1) * size).Limit(size).Find(&logsOfDB).Error
	}

	if err != nil {
		return page, err
	}
	for _, mo := range logsOfDB {
		logsOfDTO = append(logsOfDTO, dto.SystemLog{SystemLog: mo})
	}
	page.Total = total
	page.Items = logsOfDTO
	return page, err
}
