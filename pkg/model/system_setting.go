package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type SystemSetting struct {
	common.BaseModel
	ID    string `json:"id" gorm:"type:varchar(64)"`
	Key   string `json:"key" gorm:"type:varchar(256);not null;unique"`
	Value string `json:"value" gorm:"type:varchar(256);not null;"`
	Tab   string `json:"tab"`
}

func (s *SystemSetting) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return err
}
