package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type UserSetting struct {
	common.BaseModel
	ID      string `json:"-"`
	Msg     string `json:"msg"`
	Receive string `json:"receive"`
	UserID  string `json:"userId"`
}

func (u *UserSetting) BeforeCreate() error {
	u.ID = uuid.NewV4().String()
	return nil
}
