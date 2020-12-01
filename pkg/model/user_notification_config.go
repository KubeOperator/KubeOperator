package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type UserNotificationConfig struct {
	common.BaseModel
	ID     string `json:"-"`
	Vars   string `json:"vars"`
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

func (u *UserNotificationConfig) BeforeCreate() error {
	u.ID = uuid.NewV4().String()
	return nil
}
