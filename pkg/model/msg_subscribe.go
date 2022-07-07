package model

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
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

func NewMsgSubscribe(name, scope, resourceId string) MsgSubscribe {
	subConfig := MsgConfig{
		DingTalk:   constant.Disable,
		Email:      constant.Disable,
		Local:      constant.Enable,
		WorkWeiXin: constant.Disable,
	}
	configB, _ := json.Marshal(subConfig)
	return MsgSubscribe{
		Name:       name,
		Type:       scope,
		ResourceID: resourceId,
		Config:     string(configB),
	}
}

type MsgConfig struct {
	DingTalk   string
	Email      string
	Local      string
	WorkWeiXin string
}
