package service

import (
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type TaskLogService interface {
	List(clusterName string) ([]dto.TaskLog, error)
	Save(taskLog *model.TaskLog) error
	Start(log *model.TaskLog) error
	End(log *model.TaskLog, success bool, message string) error
	GetRunningLogWithClusterNameAndType(clusterName string, logType string) (model.TaskLog, error)

	StartDetail(detail *model.TaskLogDetail) error
	EndDetail(detail *model.TaskLogDetail, statu string, message string) error
	SaveDetail(detail *model.TaskLogDetail) error
}

type taskLogService struct {
}

func NewTaskLogService() TaskLogService {
	return &taskLogService{}
}

func (c *taskLogService) List(clusterName string) ([]dto.TaskLog, error) {
	var cluster model.Cluster
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return nil, err
	}
	var mos []model.TaskLog
	if err := db.DB.Where("cluster_id = ?", cluster.ID).
		Order("created_at desc").
		Find(&mos).
		Error; err != nil {
		return nil, err
	}

	var items []dto.TaskLog
	for _, mo := range mos {
		items = append(items, dto.TaskLog{TaskLog: mo})
	}
	return items, nil
}

func (c *taskLogService) SaveDetail(detail *model.TaskLogDetail) error {
	if db.DB.NewRecord(detail) {
		return db.DB.Create(detail).Error
	} else {
		return db.DB.Save(detail).Error
	}
}

func (c *taskLogService) StartDetail(detail *model.TaskLogDetail) error {
	detail.StartTime = time.Now()
	detail.Status = constant.StatusWaiting
	return db.DB.Save(detail).Error
}

func (c *taskLogService) EndDetail(detail *model.TaskLogDetail, status string, message string) error {
	detail.EndTime = time.Now()
	detail.Status = status
	detail.Message = message

	return db.DB.Save(detail).Error
}

func (c *taskLogService) GetRunningLogWithClusterNameAndType(clusterName string, logType string) (model.TaskLog, error) {
	var (
		item model.Cluster
		log  model.TaskLog
	)
	var cluster model.Cluster
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return log, err
	}
	now := time.Now()
	h, _ := time.ParseDuration("-12h")
	halfDayAgo := now.Add(h)
	if err := db.DB.Where("cluster_id = ? AND type = ? AND status = ? AND created_at BETWEEN ? AND ?", cluster.ID, logType, constant.ClusterRunning, halfDayAgo, now).
		Find(&item).
		Error; err != nil {
		return log, err
	}
	return log, nil
}

func (c *taskLogService) Save(taskLog *model.TaskLog) error {
	if db.DB.NewRecord(taskLog) {
		return db.DB.Create(taskLog).Error
	} else {
		return db.DB.Save(taskLog).Error
	}
}

func (c *taskLogService) Start(log *model.TaskLog) error {
	log.StartTime = time.Now()
	log.Phase = constant.StatusWaiting
	return db.DB.Save(log).Error
}

func (c *taskLogService) End(log *model.TaskLog, success bool, message string) error {
	log.EndTime = time.Now()
	if success {
		log.Phase = constant.TaskLogStatusSuccess
	} else {
		log.Phase = constant.TaskLogStatusFailed
	}
	log.Message = message

	return db.DB.Save(log).Error
}
