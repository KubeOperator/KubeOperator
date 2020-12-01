package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type License struct {
	common.BaseModel
	ID      string `json:"-"`
	Content string `json:"-" gorm:"type:text(65535)"`
}

func (n *License) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}
