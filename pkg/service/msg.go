package service

import (
	"bytes"
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/bindata"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	msgClient "github.com/KubeOperator/KubeOperator/pkg/util/msg"
	"github.com/jinzhu/gorm"
	"html/template"
	"io"
	"reflect"
	"time"
)

type MsgService interface {
	SendMsg(name, scope string, resource interface{}, success bool, content map[string]string) error
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

func (m msgService) SendMsg(name, scope string, resource interface{}, success bool, content map[string]string) error {
	var (
		msg        model.Msg
		resourceId string
	)
	msg.Name = name
	msg.Type = scope
	switch resource.(type) {
	case model.Cluster:
		re := resource.(model.Cluster)
		content["resourceName"] = re.Name
		if scope == constant.Cluster {
			resourceId = re.ID
		}
		var project model.Project
		db.DB.Where("id = ?", re.ProjectID).First(&project)
		content["projectName"] = project.Name
	case map[string]string:
		re := resource.(map[string]string)
		content["resourceName"] = re["name"]
	}
	date := time.Now().Add(time.Hour * 8).Format("2006-01-02 15:04:05")
	content["createdAt"] = date

	title := constant.MsgTitle[name]
	content["operator"] = title
	if success {
		msg.Level = constant.MsgInfo
		content["title"] = title + "成功"
	} else {
		msg.Level = constant.MsgWarning
		content["title"] = title + "失败"
	}

	var (
		subscribe      model.MsgSubscribe
		userSubscribes []model.MsgSubscribeUser
		accounts       []model.MsgAccount
		msgAccounts    map[string]model.MsgAccount
		userIds        []string
		userSettings   []model.UserSetting
		userAccounts   map[string][]model.UserSetting
	)

	operate := name
	if name != constant.ClusterInstall && name != constant.LicenseExpires {
		operate = constant.ClusterOperator
	}
	if err := db.DB.Model(model.MsgSubscribe{}).Where("name = ? AND type = ? AND resource_id = ?", operate, msg.Type, resourceId).First(&subscribe).Error; err != nil {
		return err
	}
	if reflect.DeepEqual(subscribe, model.MsgSubscribe{}) {
		return nil
	}
	if err := db.DB.Model(model.MsgSubscribeUser{}).Where("subscribe_id = ?", subscribe.ID).First(&userSubscribes).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if len(userSubscribes) == 0 {
		return nil
	}
	for _, us := range userSubscribes {
		userIds = append(userIds, us.UserID)
	}
	if err := db.DB.Model(model.UserSetting{}).Where("user_id in (?)", userIds).Find(&userSettings).Error; err != nil {
		return err
	}
	if len(userSettings) == 0 {
		return nil
	}
	msgSubDTO := dto.NewMsgSubscribeDTO(subscribe)
	if err := db.DB.Model(model.MsgAccount{}).Where("status = ?", constant.Enable).Find(&accounts).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if len(accounts) == 0 {
		return nil
	}
	msgAccounts = make(map[string]model.MsgAccount, len(accounts))
	userAccounts = make(map[string][]model.UserSetting, len(accounts)+1)
	for _, account := range accounts {
		if account.Status == constant.Enable {
			if account.Name == constant.Email && msgSubDTO.SubConfig.Email == constant.Enable {
				msgAccounts[constant.Email] = account
				for _, us := range userSettings {
					if us.GetMsgSetting().Email.Account != "" && us.GetMsgSetting().Email.Receive == constant.Enable {
						userAccounts[constant.Email] = append(userAccounts[constant.Email], us)
					}
				}
			}
			if account.Name == constant.WorkWeiXin && msgSubDTO.SubConfig.WorkWeiXin == constant.Enable {
				msgAccounts[constant.WorkWeiXin] = account
				for _, us := range userSettings {
					if us.GetMsgSetting().WorkWeiXin.Account != "" && us.GetMsgSetting().WorkWeiXin.Receive == constant.Enable {
						userAccounts[constant.WorkWeiXin] = append(userAccounts[constant.WorkWeiXin], us)
					}
				}
			}
			if account.Name == constant.DingTalk && msgSubDTO.SubConfig.DingTalk == constant.Enable {
				msgAccounts[constant.DingTalk] = account
				for _, us := range userSettings {
					if us.GetMsgSetting().DingTalk.Account != "" && us.GetMsgSetting().DingTalk.Receive == constant.Enable {
						userAccounts[constant.DingTalk] = append(userAccounts[constant.DingTalk], us)
					}
				}
			}
		}
	}
	if msgSubDTO.SubConfig.Local == constant.Enable {
		for _, us := range userSettings {
			userAccounts[constant.Local] = append(userAccounts[constant.Local], us)
		}
	}

	go sendUserMegs(msgAccounts, userAccounts, msg, content)

	return nil
}

func sendUserMegs(msgAccounts map[string]model.MsgAccount, userAccounts map[string][]model.UserSetting, msg model.Msg, content map[string]string) {
	by, err := json.Marshal(content)
	if err != nil {
		logger.Log.Errorf("send message failed,create msg error: %v\n", err.Error())
	}
	msg.Content = string(by)
	if err := db.DB.Create(&msg).Error; err != nil {
		logger.Log.Errorf("send message failed,create msg error: %v\n", err.Error())
	}
	for _, l := range userAccounts[constant.Local] {
		userMsg := &model.UserMsg{
			MsgID:      msg.ID,
			UserID:     l.UserID,
			SendStatus: constant.SendSuccess,
			ReadStatus: constant.UnRead,
			SendType:   constant.Local,
		}
		db.DB.Create(&userMsg)
	}

	for k, v := range msgAccounts {
		client, err := msgClient.NewMsgClient(k, dto.CoverToConfig(k, v.Config))
		if err != nil {
			logger.Log.Errorf("send message failed,create msg client error: %v\n", err.Error())
			continue
		}
		uAccounts := userAccounts[k]
		receivers := []string{}
		userMsgs := []model.UserMsg{}
		for _, ua := range uAccounts {
			var receive string
			if k == constant.Email {
				receive = ua.GetMsgSetting().Email.Account
			}
			if k == constant.WorkWeiXin {
				receive = ua.GetMsgSetting().WorkWeiXin.Account
			}
			if k == constant.DingTalk {
				receive = ua.GetMsgSetting().DingTalk.Account
			}
			if receive != "" {
				receivers = append(receivers, receive)
			}
			userMsgs = append(userMsgs, model.UserMsg{
				SendType:   k,
				UserID:     ua.UserID,
				MsgID:      msg.ID,
				Receive:    receive,
				ReadStatus: constant.UnRead,
				SendStatus: constant.SendSuccess,
			})
		}
		if len(receivers) == 0 {
			logger.Log.Errorf("send message failed,get msg receivers error: receivers is null")
			continue
		}
		detail, err := GetMsgContent(msg.Name, k, content)
		if err != nil {
			logger.Log.Errorf("send message failed,get msg content error: %v\n", err.Error())
			continue
		}
		if err := client.Send(receivers, content["title"], []byte(detail)); err != nil {
			for i := range userMsgs {
				userMsgs[i].SendStatus = constant.SendFailed
			}
		}
		for _, us := range userMsgs {
			db.DB.Create(&us)
		}
	}
}

func GetMsgContent(msgType, sendType string, content map[string]string) (string, error) {
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
