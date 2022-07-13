package model

import uuid "github.com/satori/go.uuid"

type MsgSubscribeUser struct {
	ID          string `json:"id"`
	SubscribeID string `json:"subscribeId"`
	UserID      string `json:"userId"`
}

func (m *MsgSubscribeUser) BeforeCreate() (err error) {
	m.ID = uuid.NewV4().String()
	return nil
}
