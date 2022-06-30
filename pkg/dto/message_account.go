package dto

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"reflect"
	"time"
)

type MsgAccountDTO struct {
	ID        string      `json:"-"`
	Name      string      `json:"name" validate:"containsany=EMAILDING_TALKWORK_WEIXIN"`
	Status    string      `json:"status"`
	Config    interface{} `json:"config"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

func CoverToDTO(mo model.MsgAccount) MsgAccountDTO {
	return MsgAccountDTO{
		ID:        mo.ID,
		Name:      mo.Name,
		Status:    mo.Status,
		Config:    CoverToConfig(mo.Name, mo.Config),
		CreatedAt: mo.CreatedAt,
		UpdatedAt: mo.UpdatedAt,
	}
}

func CoverToModel(dto MsgAccountDTO) model.MsgAccount {
	config, _ := json.Marshal(dto.Config)
	return model.MsgAccount{
		Name:   dto.Name,
		ID:     dto.ID,
		Status: dto.Status,
		Config: string(config),
		BaseModel: common.BaseModel{
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
	}
}

func CoverToConfig(name string, config string) interface{} {
	if name == constant.Email {
		var emailConfig EmailConfig
		_ = json.Unmarshal([]byte(config), &emailConfig)
		if reflect.DeepEqual(emailConfig, EmailConfig{}) {
			return MsgConfigs[name]
		} else {
			return emailConfig
		}
	}
	if name == constant.DingTalk {
		var dingTalkConfig DingTalkConfig
		_ = json.Unmarshal([]byte(config), &dingTalkConfig)
		if reflect.DeepEqual(dingTalkConfig, DingTalkConfig{}) {
			return MsgConfigs[name]
		} else {
			return dingTalkConfig
		}
	}
	if name == constant.WorkWeiXin {
		var workWeiXinConfig WorkWeiXinConfig
		_ = json.Unmarshal([]byte(config), &workWeiXinConfig)
		if reflect.DeepEqual(workWeiXinConfig, WorkWeiXinConfig{}) {
			return MsgConfigs[name]
		} else {
			return workWeiXinConfig
		}
	}

	return nil
}

type EmailConfig struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	TestUser string `json:"testUser"`
	Status   string `json:"status"`
}

type DingTalkConfig struct {
	WebHook  string `json:"webHook"`
	TestUser string `json:"testUser"`
	Status   string `json:"status"`
	Secret   string `json:"secret"`
}

type WorkWeiXinConfig struct {
	CorpID     string `json:"corpId"`
	AgentID    string `json:"agentId"`
	CorpSecret string `json:"corpSecret"`
	TestUser   string `json:"testUser"`
	Status     string `json:"status"`
}

var MsgConfigs = map[string]interface{}{
	constant.Email: EmailConfig{
		Status: constant.Disable,
	},
	constant.DingTalk: DingTalkConfig{
		Status: constant.Disable,
	},
	constant.WorkWeiXin: WorkWeiXinConfig{
		Status: constant.Disable,
	},
}
