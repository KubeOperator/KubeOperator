package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Volume struct {
	common.BaseModel
	ID     string `json:"id" gorm:"type:varchar(64)"`
	HostID string `json:"hostId" gorm:"type:varchar(64)"`
	Size   string `json:"size" gorm:"type:varchar(64)"`
	Name   string `json:"name" gorm:"type:varchar(256)"`
}

func (v *Volume) BeforeCreate() (err error) {
	v.ID = uuid.NewV4().String()
	return err
}
