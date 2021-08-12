package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type UserMessage struct {
	common.BaseModel
	ID          string  `json:"id"`
	Receive     string  `json:"receive"`
	UserID      string  `json:"userId"`
	MessageID   string  `json:"messageId"`
	SendType    string  `json:"sendType"`
	SendStatus  string  `json:"sendStatus"`
	ReadStatus  string  `json:"readStatus"`
	Message     Message `json:"message"`
	ClusterName string  `json:"clusterName"`
}

func (u *UserMessage) BeforeCreate() error {
	u.ID = uuid.NewV4().String()
	return nil
}
