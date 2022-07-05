package dto

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type MsgSubscribeDTO struct {
	model.MsgSubscribe
	SubConfig interface{} `json:"subConfig"`
}

func (m *MsgSubscribeDTO) CoverToDTO(mo model.MsgSubscribe) {
	var con MsgSubConfig
	json.Unmarshal([]byte(mo.Config), &con)
	m.MsgSubscribe = mo
	m.SubConfig = con
}

type MsgSubConfig struct {
	DingTalk   string `json:"dingTalk"`
	WorkWeiXin string `json:"workWeiXin"`
	Local      string `json:"local"`
	Email      string `json:"email"`
}
