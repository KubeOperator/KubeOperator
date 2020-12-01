package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type UserReceiver struct {
	common.BaseModel
	ID     string `json:"-"`
	UserID string `json:"userId"`
	Vars   string `json:"vars"`
}

func (u *UserReceiver) BeforeCreate() error {
	u.ID = uuid.NewV4().String()
	return nil
}
