package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
)

type UserMsgService interface {
	UpdateLocalMsg(msgID string, user dto.SessionUser) error
	PageLocalMsg(num, size int, user dto.SessionUser, conditions condition.Conditions) (dto.UserMsgResponse, error)
}

type userMsgService struct {
}

func NewUserMsgService() UserMsgService {
	return userMsgService{}
}

func (u userMsgService) PageLocalMsg(num, size int, user dto.SessionUser, conditions condition.Conditions) (dto.UserMsgResponse, error) {
	var (
		res     dto.UserMsgResponse
		msgs    []model.UserMsg
		msgDTOs []dto.UserMsgDTO
		unread  int
	)
	d := db.DB.Model(model.UserMsg{})
	if err := dbUtil.WithConditions(&d, model.UserMsg{}, conditions); err != nil {
		return res, err
	}
	if err := d.Where("user_id = ? AND send_type = ?", user.UserId, constant.Local).Count(&res.Total).Order("created_at desc").Offset((num - 1) * size).Limit(size).Preload("Msg").Find(&msgs).Error; err != nil {
		return res, err
	}
	if err := db.DB.Model(model.UserMsg{}).Where("user_id = ? AND send_type = ? AND read_status = ?", user.UserId, constant.Local, constant.UnRead).Count(&unread).Error; err != nil {
		return res, err
	}

	for _, m := range msgs {
		var con map[string]string
		err := json.Unmarshal([]byte(m.Msg.Content), &con)
		if err != nil {
			return res, err
		}
		msgDTOs = append(msgDTOs, dto.UserMsgDTO{
			UserMsg: m,
			Content: con,
			Type:    m.Msg.Type,
		})
	}

	res.Items = msgDTOs
	res.Unread = unread

	return res, nil
}

func (u userMsgService) UpdateLocalMsg(msgID string, user dto.SessionUser) error {
	var old model.UserMsg
	if err := db.DB.Where("id = ? AND user_id = ?", msgID, user.UserId).Find(&old).Error; err != nil {
		return err
	}
	old.ReadStatus = constant.Read

	return db.DB.Save(&old).Error
}
