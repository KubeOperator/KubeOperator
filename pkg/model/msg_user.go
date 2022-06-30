package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type MsgUser struct {
	common.BaseModel
	ID         string `json:"-"`
	Receive    string `json:"receive"`
	UserID     string `json:"userId"`
	MsgID      string `json:"msgId"`
	SendStatus string `json:"sendStatus"`
	ReadStatus string `json:"readStatus"`
	SendType   string `json:"sendType"`
}

func (m *MsgUser) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}
