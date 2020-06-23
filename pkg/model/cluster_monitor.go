package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/google/uuid"
)

type ClusterMonitor struct {
	common.BaseModel
	ID           string `json:"id"`
	Enable       bool   `json:"enable"`
	Domain       string `json:"domain"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	DashboardUrl string `json:"dashboardUrl"`
}

func (c *ClusterMonitor) BeforeCreate() (err error) {
	c.ID = uuid.New().String()
	return nil
}
func (c ClusterMonitor) TableName() string {
	return "ko_cluster_monitor"
}
