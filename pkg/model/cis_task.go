package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type CisTask struct {
	common.BaseModel
	ID        string      `json:"-"`
	ClusterID string      `json:"clusterId"`
	StartTime time.Time   `json:"startTime"`
	EndTime   time.Time   `json:"endTime"`
	Message   string      `json:"message" gorm:"type:text(65535)"`
	Results   []CisResult `json:"results"`
	Status    string      `json:"status"`
}

func (c *CisTask) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}

func (c CisTask) TableName() string {
	return "ko_cis_task"
}

