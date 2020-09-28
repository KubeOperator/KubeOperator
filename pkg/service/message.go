package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/message"
	"github.com/KubeOperator/KubeOperator/pkg/util/message/client"
	"github.com/jinzhu/gorm"
	"time"
)

type MessageService interface {
	SendMessage(mType string, result bool, content string, clusterName string, title string) error
}

type messageService struct {
	messageRepo                repository.MessageRepository
	clusterService             ClusterService
	projectMemberRepo          repository.ProjectMemberRepository
	projectResourceRepo        repository.ProjectResourceRepository
	systemSettingService       SystemSettingService
	userNotificationConfigRepo repository.UserNotificationConfigRepository
	userReceiverRepo           repository.UserReceiverRepository
	userRepo                   repository.UserRepository
	userMessageRepo            repository.UserMessageRepository
	projectRepo                repository.ProjectRepository
}

func NewMessageService() MessageService {
	return &messageService{
		messageRepo:                repository.NewMessageRepository(),
		clusterService:             NewClusterService(),
		projectMemberRepo:          repository.NewProjectMemberRepository(),
		projectResourceRepo:        repository.NewProjectResourceRepository(),
		systemSettingService:       NewSystemSettingService(),
		userReceiverRepo:           repository.NewUserReceiverRepository(),
		userRepo:                   repository.NewUserRepository(),
		userMessageRepo:            repository.NewUserMessageRepository(),
		userNotificationConfigRepo: repository.NewUserNotificationConfigRepository(),
		projectRepo:                repository.NewProjectRepository(),
	}
}

func (m messageService) SendMessage(mType string, result bool, content string, clusterName string, title string) error {
	var msg model.Message
	msg.Type = mType
	if result {
		msg.Level = constant.MsgInfo
	} else {
		msg.Level = constant.MsgWarning
	}
	msg.Content = content
	msg.Title = title
	cluster, err := m.clusterService.Get(clusterName)
	if err != nil {
		return err
	}
	msg.ClusterID = cluster.ID
	err = m.messageRepo.Save(&msg)
	if err != nil {
		return err
	}
	userMessages, err := m.GetUserMessages(msg)
	go m.SendUserMessage(userMessages, clusterName)
	return nil
}
func (m messageService) GetContentByTitleAndType(content, title, sendType, clusterName string) string {
	date := time.Now().Format("2006-01-02 15:04:05")
	var result string
	detail := make(map[string]string)
	json.Unmarshal([]byte(content), &detail)
	cluster, err := m.clusterService.Get(clusterName)
	if err != nil {
		return ""
	}
	var project model.Project
	proResources, err := m.projectResourceRepo.ListByResourceIdAndType(cluster.ID, constant.ResourceCluster)
	if err != nil {
		return ""
	}
	if len(proResources) != 0 {
		projectId := proResources[0].ProjectID
		if err := db.DB.Where(model.Project{ID: projectId}).First(&project).Error; err != nil {
			return ""
		}
	}
	if sendType == constant.Email {
		if title == constant.ClusterEventWarning {
			result = "<html>" +
				"<head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\"></head>" +
				"<body><style> table { font-size: 14px; table-layout:fixed;border:5px solid #F2F2F2;}td { font-family: Arial; WORD-WRAP: break-word }</style>" +
				"<div align=\"center\"> <table border=\"0\" cellspacing=\"2\" cellpadding=\"2\" width=\"900\"> <tr bgcolor=\"#D1D1D1\"> " +
				"<th align=\"left\" style=\"font-size:23px;\">" + Tr(title) + "</th></tr><tr><td align=\"left\">" +
				"项目:" + project.Name + "</td></tr>" +
				"<tr><td align=\"left\">集群:" + clusterName + "</td>" +
				"</tr><tr><td align=\"left\">" + detail["name"] + "</td></tr>" +
				"<tr><td align=\"left\">类别:" + detail["type"] + "</td></tr>" +
				"<tr><td align=\"left\">原因:" + detail["reason"] + "</td></tr>" +
				"<tr><td align=\"left\">组件:" + detail["component"] + "</td></tr>" +
				"<tr><td align=\"left\">NameSpace: " + detail["namespace"] + "</td></tr>" +
				"<tr><td align=\"left\">主机:" + detail["host"] + "</td> </tr> " +
				"<tr><td align=\"left\">告警时间: " + date + " </td></tr>" +
				"<tr><td align=\"left\">详情: " + detail["message"] + "</td>/tr></table>" +
				"<p>此邮件为KubeOperator平台自动发送，请勿回复!</p></div></body></html>"
		} else {
			result = "<html><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\"></head><body>" +
				"<style>table {font-size: 14px;able-layout:fixed;border:5px solid #F2F2F2;}td {font-family: Arial; WORD-WRAP: break-word }</style>" +
				"<div align=\"center\"><table border=\"0\" cellspacing=\"2\" cellpadding=\"2\" width=\"900\">" +
				"<tr bgcolor=\"#D1D1D1\"><th align=\"left\" style=\"font-size:23px;\">" + Tr(title) + "</th></tr>" +
				"<tr><td align=\"left\">项目: " + project.Name + "</td></tr>" +
				"<tr><td align=\"left\">集群: " + clusterName + "</td></tr> " +
				"<tr><td align=\"left\">详情: " + detail["message"] + "</td></tr>" +
				"<tr><td align=\"left\">时间: " + date + "</td></tr></table> " +
				"<p>此邮件为KubeOperator平台自动发送，请勿回复!</p></div></body></html>"
		}
		return result
	}
	if sendType == constant.DingTalk || sendType == constant.WorkWeiXin {
		if title == constant.ClusterEventWarning {
			result = "### " + Tr(title) + "\n\n" +
				"> **项目**:" + project.Name + "\n\n" +
				"> **集群**:" + clusterName + "\n\n" +
				"> **名称**:" + detail["name"] + " \n\n " +
				"> **类别**:" + detail["type"] + " \n\n " +
				"> **原因**:" + detail["reason"] + " \n\n " +
				"> **组件**:" + detail["component"] + " \n\n " +
				"> **类型**:" + detail["kind"] + " \n\n " +
				"> **NameSpace**:" + detail["namespace"] + " \n\n " +
				"> **详情**:" + detail["message"] + "\n\n" +
				"> **时间**:" + date + "\n\n" +
				"<font color=\"info\">本消息由KubeOperator自动发送</font>"
		} else {
			result = "### " + Tr(title) + "\n\n" +
				"> **项目**:" + project.Name + "\n\n" +
				"> **集群**:" + clusterName + "\n\n" +
				"> **详情**:" + detail["message"] + "\n\n" +
				"> **时间**:" + date + "\n\n" +
				"<font color=\"info\">本消息由KubeOperator自动发送</font>"
		}
	}

	return result
}

func Tr(title string) string {
	var result string
	switch title {
	case constant.ClusterInstall:
		result = "集群安装"
		break
	case constant.ClusterUnInstall:
		result = "集群卸载"
		break
	case constant.ClusterUpgrade:
		result = "集群升级"
		break
	case constant.ClusterScale:
		result = "集群伸缩"
		break
	case constant.ClusterAddWorker:
		result = "集群扩容"
		break
	case constant.ClusterRemoveWorker:
		result = "集群缩容"
		break
	case constant.ClusterRestore:
		result = "集群恢复"
		break
	case constant.ClusterBackup:
		result = "集群备份"
		break
	case constant.ClusterEventWarning:
		result = "集群事件告警"
		break
	}
	return result
}

func (m messageService) SendUserMessage(messages []model.UserMessage, clusterName string) {
	userMsgRepo := repository.NewUserMessageRepository()
	for _, msg := range messages {
		systemSetting, _ := NewSystemSettingService().ListByTab(msg.SendType)
		if systemSetting.Vars != nil && systemSetting.Vars[msg.SendType+"_STATUS"] == "ENABLE" {
			vars := make(map[string]interface{})
			vars["type"] = msg.SendType
			for k, value := range systemSetting.Vars {
				vars[k] = value
			}
			mClient, err := message.NewMessageClient(vars)
			if err != nil {
				msg.SendStatus = constant.SendFailed
				_ = userMsgRepo.Save(&msg)
				log.Errorf("send message failed,create client error:", err.Error())
				continue
			}
			if msg.SendType == constant.WorkWeiXin {
				token, err := client.GetToken(vars)
				if err != nil {
					msg.SendStatus = constant.SendFailed
					_ = userMsgRepo.Save(&msg)
					log.Errorf("send message failed, get token error:", err.Error())
					continue
				}
				vars["TOKEN"] = token
			}
			vars["type"] = msg.SendType
			vars["TITLE"] = Tr(msg.Message.Title)
			vars["CONTENT"] = m.GetContentByTitleAndType(msg.Message.Content, msg.Message.Title, msg.SendType, clusterName)
			vars["RECEIVERS"] = msg.Receive
			err = mClient.SendMessage(vars)
			if err != nil {
				msg.SendStatus = constant.SendFailed
				_ = userMsgRepo.Save(&msg)
				log.Errorf("send message failed,send message error:", err.Error())
				continue
			}
			_ = userMsgRepo.Save(&msg)
		}
	}
}

func (m messageService) GetUserMessages(message model.Message) ([]model.UserMessage, error) {
	var projectId string
	var userMessages []model.UserMessage
	var userIds []string
	msgReceivers := make(map[string][]string)
	proResources, err := m.projectResourceRepo.ListByResourceIdAndType(message.ClusterID, constant.ResourceCluster)
	if err != nil {
		return nil, err
	}
	if len(proResources) != 0 {
		projectId = proResources[0].ProjectID
		projectMembers, err := m.projectMemberRepo.ListByProjectId(projectId)
		if err != nil {
			return nil, err
		}
		for _, member := range projectMembers {
			userIds = append(userIds, member.UserID)
		}
	}

	adminUsers, _ := m.userRepo.ListIsAdmin()
	for _, admin := range adminUsers {
		userIds = append(userIds, admin.ID)
	}
	for _, userId := range userIds {
		sendTypes := m.getUserSendTypes(userId, message.Type)
		if len(sendTypes) == 0 {
			continue
		}
		for _, sendType := range sendTypes {
			receiver, _ := m.GetUserReceiver(userId)
			if receiver.ID == "" || receiver.Vars[sendType] == "" {
				continue
			}
			if msgReceivers[sendType] != nil {
				msgReceivers[sendType] = append(msgReceivers[sendType], receiver.Vars[sendType])
			} else {
				var res []string
				msgReceivers[sendType] = append(res, receiver.Vars[sendType])
			}
		}
	}

	for k, v := range msgReceivers {
		receivers := ""
		for _, receiver := range v {
			if k == constant.Email || k == constant.DingTalk {
				if len(receivers) == 0 {
					receivers = receiver
				} else {
					receivers = receiver + ","
				}
			} else {
				if len(receivers) == 0 {
					receivers = receiver
				} else {
					receivers = receiver + "|"
				}
			}
		}
		userMessage := model.UserMessage{
			UserID:     "",
			MessageID:  message.ID,
			SendStatus: constant.SendSuccess,
			ReadStatus: constant.UnRead,
			SendType:   k,
			Receive:    receivers,
			Message:    message,
		}
		userMessages = append(userMessages, userMessage)
	}
	m.AddLocalUserMessage(message, userIds)
	return userMessages, nil
}

func (m messageService) AddLocalUserMessage(message model.Message, userIds []string) {
	for _, userId := range userIds {
		userConfig, err := m.GetUserNotificationConfig(userId, message.Type)
		if err != nil {
			continue
		}
		if userConfig.Vars[constant.LocalMail] == "ENABLE" {
			userMessage := model.UserMessage{
				UserID:     userId,
				MessageID:  message.ID,
				SendStatus: constant.SendSuccess,
				ReadStatus: constant.UnRead,
				SendType:   constant.LocalMail,
			}
			_ = m.userMessageRepo.Save(&userMessage)
		}
	}
}

func (m messageService) getUserSendTypes(userId string, mType string) []string {
	var sendTypes []string
	userConfig, err := m.GetUserNotificationConfig(userId, mType)
	if err != nil {
		return sendTypes
	}
	smtp, _ := m.systemSettingService.Get("EMAIL_STATUS")
	if smtp.ID != "" && smtp.Value == "ENABLE" && userConfig.Vars[constant.Email] == "ENABLE" {
		sendTypes = append(sendTypes, constant.Email)
	}
	dingTalk, _ := m.systemSettingService.Get("DINGTALK_STATUS")
	if dingTalk.ID != "" && dingTalk.Value == "ENABLE" && userConfig.Vars[constant.DingTalk] == "ENABLE" {
		sendTypes = append(sendTypes, constant.DingTalk)
	}
	workWeixin, _ := m.systemSettingService.Get("WORK_WEIXIN_STATUS")
	if workWeixin.ID != "" && workWeixin.Value == "ENABLE" && userConfig.Vars[constant.WorkWeiXin] == "ENABLE" {
		sendTypes = append(sendTypes, constant.WorkWeiXin)
	}
	return sendTypes
}

func (m messageService) GetUserNotificationConfig(userId string, mType string) (*dto.UserNotificationConfigDTO, error) {
	var result dto.UserNotificationConfigDTO
	config, err := m.userNotificationConfigRepo.GetByType(userId, mType)
	if err != nil {
		return nil, err
	}
	v := make(map[string]string)
	json.Unmarshal([]byte(config.Vars), &v)
	result = dto.UserNotificationConfigDTO{
		ID:     config.ID,
		UserID: config.UserID,
		Type:   config.Type,
		Vars:   v,
	}
	return &result, nil
}

func (m messageService) GetUserReceiver(userId string) (*dto.UserReceiverDTO, error) {
	var result dto.UserReceiverDTO
	userReceiver, err := m.userReceiverRepo.Get(userId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	result.ID = userReceiver.ID
	result.UserID = userId
	v := make(map[string]string)
	json.Unmarshal([]byte(userReceiver.Vars), &v)
	result.Vars = v
	return &result, err
}
