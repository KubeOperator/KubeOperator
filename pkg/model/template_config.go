package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type TemplateConfig struct {
	common.BaseModel
	ID     string `json:"-" gorm:"type:varchar(64)"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
}

func (t *TemplateConfig) BeforeCreate() (err error) {
	t.ID = uuid.NewV4().String()
	return nil
}
