package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/ldap"
	"github.com/jinzhu/gorm"
)

type SystemSettingService interface {
	Get(name string) (dto.SystemSetting, error)
	GetLocalHostName() string
	List() (dto.SystemSettingResult, error)
	Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error)
	LdapCreate(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error)
	LdapSync(creation dto.SystemSettingCreate) error
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

func (s systemSettingService) Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error) {

	var result []dto.SystemSetting
	for k, v := range creation.Vars {
		systemSetting, err := s.systemSettingRepo.Get(k)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				systemSetting.Key = k
				systemSetting.Value = v
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
			err := s.systemSettingRepo.Save(&systemSetting)
			if err != nil {
				return result, err
			}
			result = append(result, dto.SystemSetting{SystemSetting: systemSetting})
		}
	}
	return result, nil
}

func (s systemSettingService) LdapCreate(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error) {

	var result []dto.SystemSetting
	err := s.ldapValidCheck(creation)
	if err != nil {
		return nil, err
	}
	for k, v := range creation.Vars {
		systemSetting, err := s.systemSettingRepo.Get(k)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				systemSetting.Key = k
				systemSetting.Value = v
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
			err := s.systemSettingRepo.Save(&systemSetting)
			if err != nil {
				return result, err
			}
			result = append(result, dto.SystemSetting{SystemSetting: systemSetting})
		}
	}
	return result, nil
}

func (s systemSettingService) LdapSync(creation dto.SystemSettingCreate) error {
	err := s.ldapValidCheck(creation)
	if err != nil {
		return err
	}
	go s.ldapSync(creation)
	return nil
}

func (s systemSettingService) GetLocalHostName() string {
	mo, err := s.systemSettingRepo.Get("ip")
	if err != nil || mo.Value == "" {
		return ""
	}
	return mo.Value
}

func (s systemSettingService) ldapValidCheck(creation dto.SystemSettingCreate) error {

	ldapClient := ldap.NewLdap(creation.Vars)
	err := ldapClient.Connect()
	if err != nil {
		return err
	}
	return nil
}

func (s systemSettingService) ldapSync(creation dto.SystemSettingCreate) {
	ldapClient := ldap.NewLdap(creation.Vars)
	err := ldapClient.Connect()
	if err != nil {
		log.Error(err)
	}
	entries, err := ldapClient.Search()
	if err != nil {
		log.Error(err)
	}
	for _, entry := range entries {
		user := new(model.User)
		for _, at := range entry.Attributes {
			if at.Name == "cn" {
				user.Name = at.Values[0]
			}
			if at.Name == "mail" {
				user.Email = at.Values[0]
			}
		}
		if user.Email == "" {
			continue
		}
		user.Type = constant.Ldap
		user.Language = "zh-CN"
		_, err := s.userRepo.Get(user.Name)
		if gorm.IsRecordNotFoundError(err) {
			err = s.userRepo.Save(user)
			if err != nil {
				log.Errorf("user "+user.Name+"add failed,Error:", err)
			}
		}
	}
}
