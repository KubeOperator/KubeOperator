package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type SystemLogService interface {
	Create(creation dto.SystemLogCreate) error
	Page(num, size int) (page.Page, error)
	List() ([]dto.SystemLog, error)
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
		OperationUnit: creation.OperationUnit,
		Operation:     creation.Operation,
		RequestPath:   creation.RequestPath,
	}
	if err := s.systemLogRepo.Save(&log); err != nil {
		return err
	}
	return nil
}

func (u systemLogService) Page(num, size int) (page.Page, error) {
	var page page.Page
	var systemLogResults []dto.SystemLog
	total, mos, err := u.systemLogRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		systemLogResults = append(systemLogResults, dto.SystemLog{SystemLog: mo})
	}
	page.Total = total
	page.Items = systemLogResults
	return page, err
}

func (u systemLogService) List() ([]dto.SystemLog, error) {
	var logDTOS []dto.SystemLog
	mos, err := u.systemLogRepo.List()
	if err != nil {
		return logDTOS, err
	}
	for _, mo := range mos {
		logDTOS = append(logDTOS, dto.SystemLog{SystemLog: mo})
	}
	return logDTOS, err
}
