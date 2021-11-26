package model

import (
	"errors"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)



type CisTask struct {
	common.BaseModel
	ID        string    `json:"id"`
	ClusterID string    `json:"clusterId"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Policy    string    `json:"policy"`
	Message   string    `json:"message" gorm:"type:text(65535)"`
	//Results   []CisTaskResult `json:"results"`
	Status    string `json:"status"`
	TotalPass int    `json:"totalPass"`
	TotalFail int    `json:"totalFail"`
	TotalWarn int    `json:"totalWarn"`
	TotalInfo int    `json:"totalInfo"`
}

func (c *CisTask) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}

func (c *CisTask) BeforeDelete() error {
	if c.Status == constant.ClusterRunning {
		return errors.New("task is running")
	}
	return nil
}

type CisTaskWithResult struct {
	CisTask
	Result string `json:"_"`
}

func (CisTaskWithResult) TableName() string {
	return "ko_cis_task"
}
