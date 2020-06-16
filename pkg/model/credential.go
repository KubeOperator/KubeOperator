package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Credential struct {
	common.BaseModel
	ID         string `json:"id"`
	Name       string `json:"name"gorm:"not null;unique"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	Type       string `json:"type"`
}

func (c *Credential) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return err
}

func (c Credential) TableName() string {
	return "ko_credential"
}
