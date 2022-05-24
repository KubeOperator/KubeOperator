package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/ldap"
	"github.com/jinzhu/gorm"
	"strings"
)

type LdapService interface {
	Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error)
	LdapSync(creation dto.SystemSettingCreate) error
}

type ldapService struct {
	systemSettingRepo repository.SystemSettingRepository
	userRepo          repository.UserRepository
}

func NewLdapService() LdapService {
	return &ldapService{
		systemSettingRepo: repository.NewSystemSettingRepository(),
		userRepo:          repository.NewUserRepository(),
	}
}

func (l ldapService) Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error) {
	var result []dto.SystemSetting
	err := l.ldapValidCheck(creation)
	if err != nil {
		return nil, err
	}
	for k, v := range creation.Vars {
		systemSetting, err := l.systemSettingRepo.Get(k)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				systemSetting.Key = k
				systemSetting.Value = v
				systemSetting.Tab = creation.Tab
				err := l.systemSettingRepo.Save(&systemSetting)
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
			err := l.systemSettingRepo.Save(&systemSetting)
			if err != nil {
				return result, err
			}
			result = append(result, dto.SystemSetting{SystemSetting: systemSetting})
		}
	}
	return result, nil
}

func (l ldapService) LdapSync(creation dto.SystemSettingCreate) error {
	err := l.ldapValidCheck(creation)
	if err != nil {
		return err
	}
	go l.ldapSync(creation)
	return nil
}

func (l ldapService) ldapValidCheck(creation dto.SystemSettingCreate) error {

	ldapClient, err := ldap.NewLdap(creation.Vars)
	if err != nil {
		logger.Log.Error(err)
	}
	err = ldapClient.Connect()
	if err != nil {
		return err
	}
	return nil
}

func (l ldapService) ldapSync(creation dto.SystemSettingCreate) {
	ldapClient, err := ldap.NewLdap(creation.Vars)
	if err != nil {
		logger.Log.Error(err)
	}
	err = ldapClient.Connect()
	if err != nil {
		logger.Log.Error(err)
	}
	entries, err := ldapClient.Search()
	if err != nil {
		logger.Log.Error(err)
	}
	for _, entry := range entries {
		user := new(model.User)
		for _, at := range entry.Attributes {
			if at.Name == "cn" {
				user.Name = strings.Trim(at.Values[0], " ")
			}
			if at.Name == "mail" {
				user.Email = strings.Trim(at.Values[0], " ")
			}
		}
		if user.Email == "" {
			continue
		}
		user.Type = constant.Ldap
		user.Language = "zh-CN"
		_, err := l.userRepo.Get(user.Name)
		if gorm.IsRecordNotFoundError(err) {
			err = l.userRepo.Save(user)
			if err != nil {
				logger.Log.Errorf("user "+user.Name+"add failed,Error:", err)
			}
		}
	}
}
