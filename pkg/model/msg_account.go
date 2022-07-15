package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type MsgAccount struct {
	common.BaseModel
	ID     string `json:"-"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Config string `json:"config"`
}

func (m *MsgAccount) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}
