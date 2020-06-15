package model

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterStatus struct {
	commonModel.BaseModel
	ID      string
	Message string `gorm:"type:text(65535)"`
	Phase   string
}

func (s *ClusterStatus) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}
func (s ClusterStatus) TableName() string {
	return "ko_cluster_status"
}
