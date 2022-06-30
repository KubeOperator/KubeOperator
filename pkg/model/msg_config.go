package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type MsgConfig struct {
	common.BaseModel
	ID     string `json:"-"`
	Name   string `json:"name"`
	Config string `json:"config"`
	Type   string `json:"type"`
}

func (m *MsgConfig) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}
