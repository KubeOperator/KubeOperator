package service

import (
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/typeparse"
)

type SystemLogService interface {
	Create(creation dto.SystemLogCreate) error
	Page(num, size int, queryCondition dto.SystemLogQuery) (page.Page, error)
}

type systemLogService struct{}

func NewSystemLogService() SystemLogService {
	return &systemLogService{}
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

func (u systemLogService) Page(num, size int, queryCondition dto.SystemLogQuery) (page.Page, error) {
	var (
		page      page.Page
		querySQl  string
		logsOfDB  []model.SystemLog
		logsOfDTO []dto.SystemLog
		total     int
		err       error
	)

	if queryCondition.Name.Field != "" {
		nameCondition := typeparse.ParseConditionToSql(queryCondition.Name)
		querySQl += nameCondition + " AND "
	}
	if queryCondition.Operation.Field != "" {
		operationCondition := typeparse.ParseConditionToSql(queryCondition.Operation)
		querySQl += operationCondition + " AND "
	}
	if queryCondition.OperationInfo.Field != "" {
		operationInfoCondition := typeparse.ParseConditionToSql(queryCondition.OperationInfo)
		querySQl += operationInfoCondition + " AND "
	}
	if queryCondition.Quick.Field != "" {
		quickCondition := typeparse.ParseConditionQuickToSql(queryCondition.Quick, "name", "operation", "operation_info")
		querySQl += ("(" + quickCondition + ") AND ")
	}
	if strings.Contains(querySQl, " AND ") {
		querySQl = querySQl[0 : len(querySQl)-5]
	}

	if err = db.DB.Model(&model.SystemLog{}).Where(querySQl).Order("updated_at DESC").
		Count(&total).
		Offset((num - 1) * size).Limit(size).
		Find(&logsOfDB).Error; err != nil {
		return page, err
	}

	for _, mo := range logsOfDB {
		logsOfDTO = append(logsOfDTO, dto.SystemLog{SystemLog: mo})
	}
	page.Total = total
	page.Items = logsOfDTO
	return page, err
}
