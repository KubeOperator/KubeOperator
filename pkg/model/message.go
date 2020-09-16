package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Message struct {
	common.BaseModel
	ID      string `json:"-"`
	Title   string `json:"title"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Level   string `json:"level"`
}

func (m *Message) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}

func (m Message) TableName() string {
	return "ko_message"
}
