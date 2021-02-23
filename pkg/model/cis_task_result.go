package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type CisTaskResult struct {
	common.BaseModel
	ID          string `json:"-"`
	ClusterID   string `json:"clusterId"`
	CisTaskId   string `json:"cisTaskId"`
	Number      string `json:"number"`
	Desc        string `json:"desc"`
	Remediation string `json:"remediation" gorm:"type:text(65535)"`
	Status      string `json:"status"`
	Scored      bool   `json:"scored"`
}

func (c *CisTaskResult) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}
