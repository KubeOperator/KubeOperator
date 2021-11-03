package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type NtpServer struct {
	common.BaseModel
	ID      string `json:"id" gorm:"type:varchar(64)"`
	Name    string `json:"name" gorm:"type:varchar(256);not null;unique"`
	Address string `json:"address"`
	Status  string `json:"status"`
}

func (c *NtpServer) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return err
}
