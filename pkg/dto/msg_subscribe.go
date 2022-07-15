package dto

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type MsgSubscribeDTO struct {
	model.MsgSubscribe
	SubConfig MsgSubConfig `json:"subConfig"`
	Users     []model.User `json:"users"`
}

func NewMsgSubscribeDTO(subscribe model.MsgSubscribe) MsgSubscribeDTO {
	var con MsgSubConfig
	var msgDTO MsgSubscribeDTO
	json.Unmarshal([]byte(subscribe.Config), &con)
	msgDTO.MsgSubscribe = subscribe
	msgDTO.SubConfig = con
	return msgDTO
}

type MsgSubConfig struct {
	DingTalk   string `json:"dingTalk"`
	WorkWeiXin string `json:"workWeiXin"`
	Local      string `json:"local"`
	Email      string `json:"email"`
}

type MsgSubscribeUserDTO struct {
	MsgSubscribeID string   `json:"msgSubscribeId"`
	Users          []string `json:"users"`
}
