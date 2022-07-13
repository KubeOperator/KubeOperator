package service

import (
	"bytes"
	"github.com/KubeOperator/KubeOperator/bindata"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"html/template"
	"io"
)

type MsgService interface {
	GetMsgContent(msgType, sendType string, content map[string]interface{}) (string, error)
}

type msgService struct {
	UserSettingService UserSettingService
	UserService        UserService
}

func NewMsgService() MsgService {
	return &msgService{
		UserSettingService: NewUserSettingService(),
		UserService:        NewUserService(),
	}
}

func (m msgService) SendMsg(name string, resource interface{}, success bool, content map[string]string) error {

	var (
		userSettings []model.UserSetting
		msg          model.Msg
		//resourceId   string
	)
	msg.Name = name

	switch resource.(type) {
	case model.Cluster:
		re := resource.(model.Cluster)
		content["resourceName"] = re.Name
		content["createdAt"] = re.Name
		//resourceId = re.ID
		msg.Type = constant.Cluster
	case map[string]string:
		re := resource.(map[string]string)
		content["resourceName"] = re["name"]
		content["createdAt"] = re["createdAt"]
		msg.Type = constant.System
	}
	title := constant.MsgTitle[name]
	if success {
		msg.Level = constant.MsgInfo
		content["title"] = title + "成功"
	} else {
		msg.Level = constant.MsgWarning
		content["title"] = title + "失败"
	}

	//accounts, err := getMsgAccounts(name, msg.Type, resourceId)
	//if err != nil {
	//	return err
	//}
	if err := db.DB.Model(model.Msg{}).Create(&msg).Error; err != nil {
		return err
	}

	if err := db.DB.Model(model.UserSetting{}).Find(&userSettings).Error; err != nil {
		return err
	}

	return nil
}

func getMsgAccounts(name, msgType, resourceId string) ([]model.MsgAccount, error) {
	var (
		accounts    []model.MsgAccount
		msgAccounts []model.MsgAccount
		subscribe   model.MsgSubscribe
	)

	operate := name
	if name != constant.ClusterInstall && name != constant.LicenseExpires {
		operate = constant.ClusterOperator
	}
	if err := db.DB.Model(model.MsgSubscribe{}).Where("name = ? AND type = ? AND resource_id = ?", operate, msgType, resourceId).First(&subscribe).Error; err != nil {
		return nil, err
	}
	msgSubDTO := dto.NewMsgSubscribeDTO(subscribe)
	if err := db.DB.Model(model.MsgAccount{}).Where("status = ?", constant.Enable).Find(&accounts).Error; err != nil {
		return nil, err
	}
	for _, account := range accounts {
		if account.Status == constant.Enable {
			if account.Name == constant.Email && msgSubDTO.SubConfig.Email == constant.Enable {
				msgAccounts = append(msgAccounts, account)
			}
			if account.Name == constant.WorkWeiXin && msgSubDTO.SubConfig.WorkWeiXin == constant.Enable {
				msgAccounts = append(msgAccounts, account)
			}
			if account.Name == constant.DingTalk && msgSubDTO.SubConfig.DingTalk == constant.Enable {
				msgAccounts = append(msgAccounts, account)
			}
		}
	}
	return msgAccounts, nil
}

func (m msgService) GetMsgContent(msgType, sendType string, content map[string]interface{}) (string, error) {
	tempUrl := constant.Templates[msgType][sendType]
	data, err := bindata.Asset(tempUrl)
	if err != nil {
		return "", err
	}
	newTm := template.New(sendType)
	tm, err := newTm.Parse(string(data))
	if err != nil {
		return "", err
	}
	reader, outStream := io.Pipe()
	go func() {
		err = tm.Execute(outStream, content)
		if err != nil {
			panic(err)
		}
		outStream.Close()
	}()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
