package service

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/ldap"
	"github.com/jinzhu/gorm"
	"strconv"
)

type LdapService interface {
	Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error)
	LdapSync(creation dto.SystemSettingCreate) error
	TestConnect(creation dto.SystemSettingCreate) (int, error)
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

func (l ldapService) ToLdap(creation dto.SystemSettingCreate) (dto.LdapSetting, error) {
	config := dto.LdapSetting{}
	con, err := json.Marshal(creation.Vars)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(con, &config)
	if err != nil {
		return config, err
	}
	config.SizeLimit, err = strconv.Atoi(creation.Vars["size_limit"])
	if err != nil {
		return config, err
	}
	config.TimeLimit, err = strconv.Atoi(creation.Vars["time_limit"])
	if err != nil {
		return config, err
	}
	return config, nil
}

func (l ldapService) TestConnect(creation dto.SystemSettingCreate) (int, error) {
	users := 0
	setting, err := l.ToLdap(creation)
	if err != nil {
		return users, err
	}
	if setting.Status != "ENABLE" {
		return users, errors.New("请先启用LDAP")
	}
	ldapClient, err := ldap.NewLdap(creation.Vars)
	if err != nil {
		return users, nil
	}
	if err := ldapClient.Connect(); err != nil {
		return users, err
	}
	attributes, err := ldapClient.Config.GetAttributes()
	if err != nil {
		return users, err
	}
	entries, err := ldapClient.Search(setting.UserDn, setting.Filter, setting.SizeLimit, setting.TimeLimit, attributes)
	if err != nil {
		return users, err
	}
	if len(entries) == 0 {
		return users, nil
	}

	return len(entries), nil
}

//func (l ldapService) TestLogin(username,password string) error  {
//
//}

func (l ldapService) ldapSync(creation dto.SystemSettingCreate) {
	//ldapClient, err := ldap.NewLdap(creation.Vars)
	//if err != nil {
	//	logger.Log.Error(err)
	//}
	//err = ldapClient.Connect()
	//if err != nil {
	//	logger.Log.Error(err)
	//}
	//
	//attributes, err := ldapClient.Config.GetAttributes()
	//if err != nil {
	//	logger.Log.Error(err)
	//}
	//mappings, err := ldapClient.Config.GetMappings()
	//if err != nil {
	//	logger.Log.Error(err)
	//}
	//entries, err := ldapClient.Search(attributes)
	//if err != nil {
	//	logger.Log.Error(err)
	//}
	//
	//for _, entry := range entries {
	//	user := new(model.User)
	//	rv := reflect.ValueOf(&user).Elem().Elem()
	//	for _, at := range entry.Attributes {
	//		for k, v := range mappings {
	//			if v == at.Name && len(at.Values) > 0 {
	//				fv := rv.FieldByName(k)
	//				if fv.IsValid() {
	//					fv.Set(reflect.ValueOf(strings.Trim(at.Values[0], " ")))
	//				}
	//			}
	//		}
	//	}
	//	if user.Email == "" || user.Name == "" {
	//		continue
	//	}
	//	user.Type = constant.Ldap
	//	user.Language = "zh-CN"
	//	_, err := l.userRepo.Get(user.Name)
	//	if gorm.IsRecordNotFoundError(err) {
	//		err = l.userRepo.Save(user)
	//		if err != nil {
	//			logger.Log.Errorf("user "+user.Name+"add failed,Error:", err)
	//		}
	//	}
	//}
}
