package service

import (
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/message"
	"github.com/KubeOperator/KubeOperator/pkg/util/message/client"
	"github.com/jinzhu/gorm"
)

type SystemSettingService interface {
	Get(name string) (dto.SystemSetting, error)
	GetLocalIP() (string, error)
	List() (dto.SystemSettingResult, error)
	Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error)
	ListByTab(tabName string) (dto.SystemSettingResult, error)
	CheckSettingByType(tabName string, creation dto.SystemSettingCreate) error
}

type systemSettingService struct {
	systemSettingRepo repository.SystemSettingRepository
	userRepo          repository.UserRepository
}

func NewSystemSettingService() SystemSettingService {
	return &systemSettingService{
		systemSettingRepo: repository.NewSystemSettingRepository(),
		userRepo:          repository.NewUserRepository(),
	}
}

func (s systemSettingService) Get(key string) (dto.SystemSetting, error) {
	var systemSettingDTO dto.SystemSetting
	mo, err := s.systemSettingRepo.Get(key)
	if err != nil {
		return systemSettingDTO, err
	}
	systemSettingDTO.SystemSetting = mo
	return systemSettingDTO, err
}

func (s systemSettingService) List() (dto.SystemSettingResult, error) {
	var systemSettingResult dto.SystemSettingResult
	vars := make(map[string]string)
	mos, err := s.systemSettingRepo.List()
	if err != nil {
		return systemSettingResult, err
	}
	for _, mo := range mos {
		vars[mo.Key] = mo.Value
	}
	systemSettingResult.Vars = vars
	return systemSettingResult, err
}
func (s systemSettingService) ListByTab(tabName string) (dto.SystemSettingResult, error) {
	var systemSettingResult dto.SystemSettingResult
	vars := make(map[string]string)
	mos, err := s.systemSettingRepo.ListByTab(tabName)
	if err != nil {
		return systemSettingResult, err
	}
	for _, mo := range mos {
		vars[mo.Key] = mo.Value
	}
	if len(mos) > 0 {
		systemSettingResult.Tab = tabName
	}
	systemSettingResult.Vars = vars
	return systemSettingResult, err
}

func (s systemSettingService) Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error) {

	var result []dto.SystemSetting
	for k, v := range creation.Vars {
		systemSetting, err := s.systemSettingRepo.Get(k)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				systemSetting.Key = k
				systemSetting.Value = v
				systemSetting.Tab = creation.Tab
				err := s.systemSettingRepo.Save(&systemSetting)
				if err != nil {
					return result, err
				}
				result = append(result, dto.SystemSetting{SystemSetting: systemSetting})
			} else {
				return result, err
			}
		} else if systemSetting.ID != "" {
			systemSetting.Value = v
			if systemSetting.Tab == "" {
				systemSetting.Tab = creation.Tab
			}
			err := s.systemSettingRepo.Save(&systemSetting)
			if err != nil {
				return result, err
			}
			result = append(result, dto.SystemSetting{SystemSetting: systemSetting})
		}
	}
	return result, nil
}

func (s systemSettingService) GetLocalIP() (string, error) {
	var arch_type model.SystemSetting
	if err := db.DB.Model(&model.SystemSetting{}).Where("key = ?", "arch_type").First(&arch_type).Error; err != nil {
		return "", fmt.Errorf("can't found arch_type from system setting, err %s", err.Error())
	}

	if arch_type.Value == "single" {
		var sysSetting model.SystemSetting
		if err := db.DB.Model(&model.SystemSetting{}).Where("key = ?", "ip").First(&sysSetting).Error; err != nil {
			return "", fmt.Errorf("can't found ip from system setting, err %s", err.Error())
		}
		return sysSetting.Value, nil
	}
	var sysRegistry model.SystemRegistry
	if err := db.DB.Model(&model.SystemRegistry{}).Where("architecture = ?", "amd64").First(&sysRegistry).Error; err != nil {
		return "", fmt.Errorf("can't found registry from system registry, err %s", err.Error())
	}
	return sysRegistry.RegistryHostname, nil
}

func (s systemSettingService) CheckSettingByType(tabName string, creation dto.SystemSettingCreate) error {

	vars := make(map[string]interface{})
	for k, value := range creation.Vars {
		vars[k] = value
	}
	if tabName == constant.Email {
		vars["type"] = constant.Email
		vars["RECEIVERS"] = vars["SMTP_TEST_USER"]
		vars["TITLE"] = "KubeOperator测试邮件"
		vars["CONTENT"] = "此邮件由 KubeOperator 发送，用于测试邮件发送，请勿回复"
	} else if tabName == constant.DingTalk {
		vars["type"] = constant.DingTalk
		vars["RECEIVERS"] = vars["DING_TALK_TEST_USER"]
		vars["TITLE"] = "KubeOperator测试消息"
		vars["CONTENT"] = "此邮件由 KubeOperator 发送，用于测试消息发送"
	} else if tabName == constant.WorkWeiXin {
		vars["type"] = constant.WorkWeiXin
		vars["CONTENT"] = "此邮件由 KubeOperator 发送，用于测试消息发送"
		vars["RECEIVERS"] = vars["WORK_WEIXIN_TEST_USER"]
	}
	c, err := message.NewMessageClient(vars)
	if err != nil {
		return err
	}
	if tabName == constant.WorkWeiXin {
		token, err := client.GetToken(vars)
		if err != nil {
			return err
		}
		vars["TOKEN"] = token
	}
	err = c.SendMessage(vars)
	if err != nil {
		return err
	}
	return nil
}
