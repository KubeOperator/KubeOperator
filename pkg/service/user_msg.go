package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
)

type UserMsgService interface {
	UpdateLocalMsg(msgID string, user dto.SessionUser) error
	PageLocalMsg(num, size int, user dto.SessionUser, conditions condition.Conditions) (page.Page, error)
}

type userMsgService struct {
}

func NewUserMsgService() UserMsgService {
	return userMsgService{}
}

func (u userMsgService) PageLocalMsg(num, size int, user dto.SessionUser, conditions condition.Conditions) (page.Page, error) {
	var (
		p           page.Page
		msgs        []model.UserMsg
		userMsgDTOS []dto.UserMsgDTO
	)
	d := db.DB.Model(model.UserMsg{})
	if err := dbUtil.WithConditions(&d, model.UserMsg{}, conditions); err != nil {
		return page.Page{}, err
	}
	if err := d.Where("user_id = ?", user.UserId).Count(&p.Total).Order("created_at desc").Offset((num - 1) * size).Limit(size).Preload("Msg").Find(&msgs).Error; err != nil {
		return page.Page{}, err
	}

	for _, msg := range msgs {
		userMsgDTOS = append(userMsgDTOS, dto.UserMsgDTO{
			UserMsg: msg,
		})
	}
	p.Items = userMsgDTOS

	return p, nil
}

func (u userMsgService) UpdateLocalMsg(msgID string, user dto.SessionUser) error {
	var old model.UserMsg
	if err := db.DB.Where("id = ? AND user_id = ?", msgID, user.UserId).Find(&old).Error; err != nil {
		return err
	}
	old.ReadStatus = constant.Read
	return nil
}
