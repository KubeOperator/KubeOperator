package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type SystemLog struct {
	common.BaseModel
	ID            string `json:"-" gorm:"type:varchar(64)"`
	Name          string `json:"name" gorm:"type:varchar(256);not null;"`
	Operation     string `json:"operation" gorm:"type:varchar(256);not null;"`
	OperationInfo string `json:"operationInfo" gorm:"type:varchar(256);"`
}

func (s *SystemLog) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return err
}
