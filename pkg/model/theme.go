package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Theme struct {
	common.BaseModel
	ID         string `json:"-" gorm:"type:varchar(64)"`
	SystemName string `json:"systemName"`
	Logo       string `json:"logo" gorm:"type:text"`
}

func (c *Theme) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}
