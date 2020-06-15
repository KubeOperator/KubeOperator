package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type ClusterStatusCondition struct {
	common.BaseModel
	ID            string
	Name          string
	StatusID      string
	Status        string
	Message       string `gorm:"type:text"`
	OrderNum      int
	LastProbeTime time.Time
}

func (c *ClusterStatusCondition) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}

func (c ClusterStatusCondition) TableName() string {
	return "ko_cluster_status_condition"
}
