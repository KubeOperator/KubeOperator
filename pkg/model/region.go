package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Region struct {
	common.BaseModel
	ID         string `json:"id" gorm:"type:varchar(64)"`
	Name       string `json:"name" gorm:"type:varchar(256);not null;unique"`
	Datacenter string `json:"datacenter" gorm:"type:varchar(64)"`
	Provider   string `json:"provider" gorm:"type:varchar(64)"`
	Vars       string `json:"vars" gorm:"type text(65535)"`
}

func (c *Region) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return err
}
