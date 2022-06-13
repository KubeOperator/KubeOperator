package model

import (
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

// log + details if apply single
type TaskLog struct {
	common.BaseModel
	ID        string    `json:"id"`
	ClusterID string    `json:"clusterID"`
	Type      string    `json:"type"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`

	Phase    string `json:"phase"`
	PrePhase string `json:"prePhase"`
	Message  string `json:"message" gorm:"type:text(65535)"`

	Details []TaskLogDetail `json:"details"`
}

// detail if apply multi
type TaskLogDetail struct {
	common.BaseModel
	ID            string    `json:"id"`
	Task          string    `json:"task"`
	TaskID        string    `json:"taskID"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	LastProbeTime time.Time `json:"lastProbeTime"`
	Message       string    `json:"message" gorm:"type:text(65535)"`
	Status        string    `json:"status"`
}

func (n *TaskLog) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}

func (n *TaskLogDetail) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}
