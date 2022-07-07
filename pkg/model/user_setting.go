package model

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type UserSetting struct {
	common.BaseModel
	ID     string `json:"id"`
	Msg    string `json:"-"`
	UserID string `json:"userId"`
}

func (u *UserSetting) BeforeCreate() error {
	u.ID = uuid.NewV4().String()
	return nil
}

func (u *UserSetting) GetMsgSetting() MsgSetting {
	var msg MsgSetting
	json.Unmarshal([]byte(u.Msg), &msg)
	return msg
}

func NewUserSetting(userId string) UserSetting {
	msgConfig := &MsgSetting{
		DingTalk: ReceiveSetting{
			Account: "",
			Receive: constant.Disable,
		},
		Email: ReceiveSetting{
			Account: "",
			Receive: constant.Disable,
		},
		WorkWeiXin: ReceiveSetting{
			Account: "",
			Receive: constant.Disable,
		},
		Local: ReceiveSetting{
			Account: "",
			Receive: constant.Enable,
		},
	}
	msg, _ := json.Marshal(msgConfig)
	return UserSetting{
		UserID: userId,
		Msg:    string(msg),
	}
}

type MsgSetting struct {
	DingTalk   ReceiveSetting `json:"dingTalk"`
	Email      ReceiveSetting `json:"email"`
	WorkWeiXin ReceiveSetting `json:"workWeiXin"`
	Local      ReceiveSetting `json:"local"`
}

type ReceiveSetting struct {
	Account string `json:"account"`
	Receive string `json:"receive"`
}
