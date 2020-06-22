package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/google/uuid"
)

type ClusterMonitor struct {
	common.BaseModel
	ID      string
	Enable  bool
	Domain  string
	Status  string
	Message string
}

func (c *ClusterMonitor) BeforeCreate() (err error) {
	c.ID = uuid.New().String()
	return nil
}
func (c ClusterMonitor) TableName() string {
	return "ko_cluster_monitor"
}
