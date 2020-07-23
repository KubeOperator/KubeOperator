package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Credential struct {
	common.BaseModel
	ID         string `json:"id" gorm:"type:varchar(64)"`
	Name       string `json:"name" gorm:"type:varchar(256);not null;unique"`
	Username   string `json:"username" gorm:"type:varchar(64)"`
	Password   string `json:"password" gorm:"type:varchar(256)"`
	PrivateKey string `json:"privateKey" gorm:"type: text(0)"`
	Type       string `json:"type" gorm:"type:varchar(64)"`
}

func (c *Credential) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return err
}

func (c Credential) TableName() string {
	return "ko_credential"
}
