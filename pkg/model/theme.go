package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Theme struct {
	common.BaseModel
	ID           string `json:"-" gorm:"type:varchar(64)"`
	SystemName   string `json:"systemName"`
	LoginImage   string `json:"loginImage"`
	Logo         string `json:"logo"`
	LogoWithText string `json:"logoWithText"`
	Icon         string `json:"icon"`
	LogoAbout    string `json:"logoAbout"`
}

func (c *Theme) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
}
