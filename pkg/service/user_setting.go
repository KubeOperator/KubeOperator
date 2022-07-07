package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type UserSettingService interface {
	Update(settingDTO dto.UserSettingDTO) (dto.UserSettingDTO, error)
	GetByUsername(username string) (dto.UserSettingDTO, error)
}

type userSettingService struct {
	UserService UserService
}

func NewUserSettingService() UserSettingService {
	return &userSettingService{
		UserService: NewUserService(),
	}
}

func (u userSettingService) GetByUsername(username string) (dto.UserSettingDTO, error) {
	var (
		setting    model.UserSetting
		settingDto dto.UserSettingDTO
	)
	user, err := u.UserService.Get(username)
	if err != nil {
		return settingDto, err
	}

	if err := db.DB.Model(&model.MsgSetting{}).Where("user_id = ?", user.ID).First(&setting).Error; err != nil {
		return settingDto, err
	}

	settingDto.UserSetting = setting
	msgConfig := setting.GetMsgSetting()
	if msgConfig.Email.Account == "" {
		msgConfig.Email.Account = user.Email
	}
	settingDto.MsgConfig = msgConfig

	return settingDto, nil
}

func (u userSettingService) Update(settingDTO dto.UserSettingDTO) (dto.UserSettingDTO, error) {
	setting := settingDTO.UserSetting
	msgSetting, err := settingDTO.GetMsgConfig()
	if err != nil {
		return settingDTO, err
	}
	setting.Msg = msgSetting

	return settingDTO, db.DB.Save(&setting).Error
}
