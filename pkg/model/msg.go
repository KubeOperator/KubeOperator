package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Msg struct {
	common.BaseModel
	ID      string `json:"-"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Level   string `json:"level"`
}

func (m *Msg) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}
