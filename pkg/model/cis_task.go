package model

import (
	"errors"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type CisTask struct {
	common.BaseModel
	ID        string          `json:"id"`
	ClusterID string          `json:"clusterId"`
	StartTime time.Time       `json:"startTime"`
	EndTime   time.Time       `json:"endTime"`
	Message   string          `json:"message" gorm:"type:text(65535)"`
	Results   []CisTaskResult `json:"results"`
	Status    string          `json:"status"`
}

func (c *CisTask) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}

func (c *CisTask) BeforeDelete() error {
	if c.Status == constant.ClusterRunning {
		return errors.New("task is running")
	}
	if err := db.DB.Where(CisTaskResult{CisTaskId: c.ID}).Delete(&CisTaskResult{}).Error; err != nil {
		return err
	}
	return nil
}
