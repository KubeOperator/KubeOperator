package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type MsgSubscribe struct {
	common.BaseModel
	ID         string `json:"-"`
	Name       string `json:"name"`
	Config     string `json:"-"`
	Type       string `json:"type"`
	ResourceID string `json:"resourceId"`
}

func (m *MsgSubscribe) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}
