package model

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterStatus struct {
	commonModel.BaseModel
	ID         string
	Message    string                   `json:"message" gorm:"type:text(65535)"`
	Phase      string                   `json:"phase"`
	Conditions []ClusterStatusCondition `json:"conditions" gorm:"save_associations:false"`
}

func (s *ClusterStatus) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}
func (s ClusterStatus) TableName() string {
	return "ko_cluster_status"
}
