package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

// log + details if apply single
type TaskLog struct {
	common.BaseModel
	ID        string `json:"id"`
	ClusterID string `json:"clusterID"`
	Type      string `json:"type"`

	Phase     string `json:"phase"`
	Message   string `json:"message" gorm:"type:text(65535)"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`

	Details []TaskLogDetail `json:"details"`
}

// detail if apply multi
type TaskLogDetail struct {
	common.BaseModel
	ID            string `json:"id"`
	Task          string `json:"task"`
	TaskLogID     string `json:"taskLogID"`
	ClusterID     string `json:"clusterID"`
	LastProbeTime int64  `json:"lastProbeTime"`
	StartTime     int64  `json:"startTime"`
	EndTime       int64  `json:"endTime"`
	Message       string `json:"message" gorm:"type:text(65535)"`
	Status        string `json:"status"`
}

type TaskRetryLog struct {
	common.BaseModel
	ID        string `json:"id"`
	TaskLogID string `json:"taskLogID"`
	ClusterID string `json:"clusterID"`
	Message   string `json:"message" gorm:"type:text(65535)"`
}

func (n *TaskLog) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}

func (n *TaskLogDetail) BeforeCreate() (err error) {
	if len(n.ID) == 0 {
		n.ID = uuid.NewV4().String()
	}
	return nil
}
