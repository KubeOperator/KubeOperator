package tool

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Status struct {
	commonModel.BaseModel
	ID         string
	Message    string `gorm:"type:text(65535)"`
	Phase      string
	Conditions []Condition `gorm:"save_associations:false"`
}

func (s *Status) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s Status) TableName() string {
	return "ko_tool_status"
}
