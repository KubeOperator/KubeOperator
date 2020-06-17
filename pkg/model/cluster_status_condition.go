package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type ClusterStatusCondition struct {
	common.BaseModel
	ID              string `json:"_"`
	Name            string `json:"name"`
	ClusterStatusID string `json:"_"`
	Status          string `json:"status"`
	Message         string `json:"message"gorm:"type:text"`
	OrderNum        int    `json:"_"`
	LastProbeTime   time.Time
}

func (c *ClusterStatusCondition) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}

func (c ClusterStatusCondition) TableName() string {
	return "ko_cluster_status_condition"
}
