package service

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	uuid "github.com/satori/go.uuid"
)

type TaskLogService interface {
	Page(num, size int, clusterName string, logtype string) (*page.Page, error)
	GetByID(id string) (model.TaskLog, error)
	GetTaskDetailByID(id string) (*dto.TaskLog, error)
	GetTaskLogByID(clusterId, logId string) (*dto.Logs, error)
	GetTaskLogByName(clusterName, logId string) (*dto.Logs, error)
	Save(taskLog *model.TaskLog) error
	Start(log *model.TaskLog) error
	End(log *model.TaskLog, success bool, message string) error
	NewTerminalTask(clusterID string, logtype string) (*model.TaskLog, error)

	IsTaskOn(clusterName string) bool

	StartDetail(detail *model.TaskLogDetail) error
	EndDetail(detail *model.TaskLogDetail, statu string, message string) error
	SaveDetail(detail *model.TaskLogDetail) error

	RestartTask(cluster *model.Cluster, operation string) error
}

type taskLogService struct {
}

func NewTaskLogService() TaskLogService {
	return &taskLogService{}
}

func (c *taskLogService) Page(num, size int, clusterName string, logtype string) (*page.Page, error) {
	var (
		datas []dto.TaskLog
		p     page.Page
	)
	if logtype == "cluster" {
		var tasklogs []model.TaskLog
		d := db.DB.Model(model.TaskLog{})
		if len(clusterName) != 0 {
			var cluster model.Cluster
			if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
				return nil, err
			}
			d = d.Where("cluster_id = ?", cluster.ID)
		}
		if err := d.Count(&p.Total).Preload("Details").Order("created_at desc").Offset((num - 1) * size).Limit(size).Find(&tasklogs).Error; err != nil {
			return nil, err
		}
		for t := 0; t < len(tasklogs); t++ {
			sort.Slice(tasklogs[t].Details, func(i, j int) bool {
				return tasklogs[t].Details[i].StartTime > tasklogs[t].Details[j].StartTime
			})
			datas = append(datas, dto.TaskLog{TaskLog: tasklogs[t]})
		}
		p.Items = datas
		return &p, nil
	}

	var tasklogs []model.TaskLogDetail
	d := db.DB.Model(model.TaskLogDetail{}).Where("cluster_id != ?", "")
	if len(clusterName) != 0 {
		var cluster model.Cluster
		if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
			return nil, err
		}
		d = d.Where("cluster_id = ?", cluster.ID)
	}
	if err := d.Count(&p.Total).Order("created_at desc").Offset((num - 1) * size).Limit(size).Find(&tasklogs).Error; err != nil {
		return nil, err
	}
	for t := 0; t < len(tasklogs); t++ {
		datas = append(datas, dto.TaskLog{
			TaskLog: model.TaskLog{
				ID:        tasklogs[t].ID,
				ClusterID: tasklogs[t].ClusterID,
				Phase:     tasklogs[t].Status,
				Message:   tasklogs[t].Message,
				Type:      tasklogs[t].Task,
				StartTime: tasklogs[t].StartTime,
				EndTime:   tasklogs[t].EndTime,
			},
		})
	}
	p.Items = datas
	return &p, nil
}

func (c *taskLogService) GetByID(id string) (model.TaskLog, error) {
	var tasklog model.TaskLog
	if err := db.DB.Where("id = ?", id).Preload("Details").First(&tasklog).Error; err != nil {
		return tasklog, err
	}
	return tasklog, nil
}

func (c *taskLogService) GetTaskDetailByID(id string) (*dto.TaskLog, error) {
	var (
		tasklog   model.TaskLog
		retrylogs []model.TaskRetryLog
	)
	if err := db.DB.Where("id = ?", id).Preload("Details").First(&tasklog).Error; err != nil {
		return &dto.TaskLog{TaskLog: tasklog}, err
	}
	if err := db.DB.Where("task_log_id = ?", id).Find(&retrylogs).Error; err != nil {
		return &dto.TaskLog{TaskLog: tasklog}, err
	}
	for _, re := range retrylogs {
		item := model.TaskLogDetail{
			Task:      constant.TaskLogStatusFailed,
			TaskLogID: re.TaskLogID,
			ClusterID: re.ClusterID,
			StartTime: re.LastFailedTime,
			EndTime:   re.LastFailedTime,
			Status:    constant.TaskLogStatusFailed,
			Message:   re.Message,
		}
		tasklog.Details = append(tasklog.Details, item)
		item2 := model.TaskLogDetail{
			Task:      constant.TaskLogStatusRedo,
			TaskLogID: re.TaskLogID,
			ClusterID: re.ClusterID,
			StartTime: re.RestartTime,
			EndTime:   re.RestartTime,
			Status:    constant.TaskLogStatusSuccess,
			Message:   re.Message,
		}
		tasklog.Details = append(tasklog.Details, item2)
	}
	sort.Slice(tasklog.Details, func(i, j int) bool {
		return tasklog.Details[i].StartTime < tasklog.Details[j].StartTime
	})
	return &dto.TaskLog{TaskLog: tasklog}, nil
}

func (c *taskLogService) GetTaskLogByID(clusterId, logId string) (*dto.Logs, error) {
	var cluster model.Cluster
	if err := db.DB.Where("id = ?", clusterId).First(&cluster).Error; err != nil {
		return nil, err
	}
	r, err := ansible.GetAnsibleLogReader(cluster.Name, logId)
	if err != nil {
		return nil, err
	}
	var chunk []byte
	for {
		buffer := make([]byte, 1024)
		n, err := r.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		chunk = append(chunk, buffer[:n]...)
	}
	return &dto.Logs{Msg: string(chunk)}, nil
}

func (c *taskLogService) GetTaskLogByName(clusterName, logId string) (*dto.Logs, error) {
	r, err := ansible.GetAnsibleLogReader(clusterName, logId)
	if err != nil {
		return nil, err
	}
	var chunk []byte
	for {
		buffer := make([]byte, 1024)
		n, err := r.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		chunk = append(chunk, buffer[:n]...)
	}
	return &dto.Logs{Msg: string(chunk)}, nil
}

func (c *taskLogService) NewTerminalTask(clusterID string, logtype string) (*model.TaskLog, error) {
	task := model.TaskLog{
		ClusterID: clusterID,
		Type:      logtype,
		Phase:     constant.TaskLogStatusRunning,
		StartTime: time.Now().Unix(),
		Details: []model.TaskLogDetail{
			{
				ID:            uuid.NewV4().String(),
				Task:          logtype,
				Status:        constant.TaskLogStatusRunning,
				LastProbeTime: time.Now().Unix(),
				StartTime:     time.Now().Unix(),
			},
		},
	}
	return &task, db.DB.Create(&task).Error
}

func (c *taskLogService) SaveDetail(detail *model.TaskLogDetail) error {
	if db.DB.NewRecord(detail) {
		detail.StartTime = time.Now().Unix()
		return db.DB.Create(detail).Error
	} else {
		return db.DB.Save(detail).Error
	}
}

func (c *taskLogService) StartDetail(detail *model.TaskLogDetail) error {
	detail.StartTime = time.Now().Unix()
	return db.DB.Create(detail).Error
}

func (c *taskLogService) EndDetail(detail *model.TaskLogDetail, status string, message string) error {
	detail.Status = status
	detail.Message = message
	detail.EndTime = time.Now().Unix()

	return db.DB.Save(detail).Error
}

func (c *taskLogService) IsTaskOn(clusterName string) bool {
	var (
		cluster model.Cluster
		log     model.TaskLog
	)
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return true
	}
	if cluster.CurrentTaskID == "" {
		return false
	}
	if err := db.DB.Where("id = ?", cluster.CurrentTaskID).First(&log).Error; err != nil {
		return false
	}
	return !(log.Phase == constant.TaskLogStatusFailed || log.Phase == constant.TaskLogStatusSuccess)
}

func (c *taskLogService) Save(taskLog *model.TaskLog) error {
	for i := 0; i < len(taskLog.Details); i++ {
		if taskLog.Details[i].ID == "" {
			taskLog.Details[i].ID = uuid.NewV4().String()
		}
	}
	if db.DB.NewRecord(taskLog) {
		taskLog.StartTime = time.Now().Unix()
		return db.DB.Create(&taskLog).Error
	} else {
		return db.DB.Save(&taskLog).Error
	}
}

func (c *taskLogService) Start(log *model.TaskLog) error {
	log.Phase = constant.TaskLogStatusWaiting
	log.StartTime = time.Now().Unix()
	return db.DB.Create(log).Error
}

func (c *taskLogService) End(log *model.TaskLog, success bool, message string) error {
	status := constant.TaskLogStatusFailed
	if success {
		status = constant.TaskLogStatusSuccess
		log.Phase = constant.TaskLogStatusSuccess
	} else {
		log.Phase = constant.TaskLogStatusFailed
	}
	log.EndTime = time.Now().Unix()
	for i := 0; i < len(log.Details); i++ {
		if log.Details[i].Status == constant.TaskLogStatusRunning {
			log.Details[i].Status = status
			log.Details[i].Message = message
			log.Details[i].EndTime = time.Now().Unix()
		}
	}
	log.Message = message

	return db.DB.Save(log).Error
}

func (c taskLogService) RestartTask(cluster *model.Cluster, operation string) error {
	isON := c.IsTaskOn(cluster.Name)
	if isON {
		return errors.New("TASK_IN_EXECUTION")
	}
	if operation != cluster.TaskLog.Type {
		return fmt.Errorf("restart failed, task type do not match %s - %s", operation, cluster.TaskLog.Type)
	}
	retrylog := &model.TaskRetryLog{
		ClusterID:   cluster.ID,
		TaskLogID:   cluster.TaskLog.ID,
		Message:     cluster.TaskLog.Message,
		RestartTime: time.Now().Unix(),
	}
	if len(cluster.TaskLog.Details) > 0 {
		for i := range cluster.TaskLog.Details {
			if cluster.TaskLog.Details[i].Status == constant.TaskLogStatusFailed {
				cluster.TaskLog.Details[i].Status = constant.TaskLogStatusRunning
				cluster.TaskLog.Details[i].Message = ""
				cluster.TaskLog.Details[i].StartTime = time.Now().Unix()
				retrylog.LastFailedTime = cluster.TaskLog.Details[i].EndTime
			}
		}
	}
	if err := db.DB.Create(retrylog).Error; err != nil {
		return err
	}
	cluster.TaskLog.Phase = constant.TaskLogStatusWaiting
	if err := c.Save(&cluster.TaskLog); err != nil {
		return fmt.Errorf("reset contidion err %s", err.Error())
	}
	return nil
}
