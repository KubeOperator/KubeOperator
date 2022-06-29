package service

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/ldap"
	"github.com/jinzhu/gorm"
	"reflect"
	"strconv"
	"strings"
)

type LdapService interface {
	Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error)
	TestConnect(creation dto.SystemSettingCreate) (int, error)
	TestLogin(username, password string) error
	LdapSync() ([]dto.LdapUser, error)
	ImpostUsers(users []dto.LdapUser) error
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

func (l ldapService) GetLdapConfig() (map[string]string, error) {
	var vars map[string]string
	settings, err := l.systemSettingRepo.ListByTab("LDAP")
	if err != nil {
		return nil, err
	}
	vars = make(map[string]string, len(settings))
	for _, s := range settings {
		vars[s.Key] = s.Value
	}
	return vars, nil
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

func (l ldapService) ToLdap(vars map[string]string) (dto.LdapSetting, error) {
	config := dto.LdapSetting{}
	con, err := json.Marshal(vars)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(con, &config)
	if err != nil {
		return config, err
	}
	config.SizeLimit, err = strconv.Atoi(vars["size_limit"])
	if err != nil {
		return config, err
	}
	config.TimeLimit, err = strconv.Atoi(vars["time_limit"])
	if err != nil {
		return config, err
	}
	return config, nil
}

func (l ldapService) TestConnect(creation dto.SystemSettingCreate) (int, error) {
	users := 0
	setting, err := l.ToLdap(creation.Vars)
	if err != nil {
		return users, err
	}
	if setting.Status != "ENABLE" {
		return users, errors.New("LDAP_DISABLE_STATUS")
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

func (l ldapService) TestLogin(username, password string) error {
	vars, err := l.GetLdapConfig()
	if err != nil {
		return err
	}
	setting, err := l.ToLdap(vars)
	if err != nil {
		return err
	}
	if setting.Status != "ENABLE" {
		return errors.New("LDAP_DISABLE_STATUS")
	}
	ldapClient, err := ldap.NewLdap(vars)
	if err != nil {
		return err
	}
	err = ldapClient.Connect()
	if err != nil {
		return err
	}
	return ldapClient.Login(username, password)
}

func (l ldapService) LdapSync() ([]dto.LdapUser, error) {
	var users []dto.LdapUser
	vars, err := l.GetLdapConfig()
	if err != nil {
		return users, err
	}
	setting, err := l.ToLdap(vars)
	if err != nil {
		return users, err
	}
	if setting.Status != "ENABLE" {
		return users, errors.New("LDAP_DISABLE_STATUS")
	}
	ldapClient, err := ldap.NewLdap(vars)
	if err != nil {
		return users, err
	}
	err = ldapClient.Connect()
	if err != nil {
		return users, err
	}
	attributes, err := ldapClient.Config.GetAttributes()
	if err != nil {
		return users, err
	}
	mappings, err := ldapClient.Config.GetMappings()
	if err != nil {
		return users, err
	}
	entries, err := ldapClient.Search(setting.UserDn, setting.Filter, setting.SizeLimit, setting.TimeLimit, attributes)
	if err != nil {
		return users, err
	}
	names, err := l.userRepo.ListNames()
	if err != nil {
		return users, err
	}
	for _, entry := range entries {
		user := &dto.LdapUser{}
		rv := reflect.ValueOf(&user).Elem().Elem()
		for _, at := range entry.Attributes {
			for k, v := range mappings {
				if v == at.Name && len(at.Values) > 0 {
					fv := rv.FieldByName(k)
					if fv.IsValid() {
						fv.Set(reflect.ValueOf(strings.Trim(at.Values[0], " ")))
					}
				}
			}
		}
		if user.Name == "" {
			continue
		}
		if names[user.Name] {
			user.Available = false
		} else {
			user.Available = true
		}
		users = append(users, *user)
	}

	return users, err
}

func (l ldapService) ImpostUsers(users []dto.LdapUser) error {
	tx := db.DB.Begin()
	for _, imp := range users {
		us := model.User{
			Name:  imp.Name,
			Email: imp.Email,
		}
		if us.Email == "" {
			us.Email = us.Name + "@example.com"
		}
		us.Type = constant.Ldap
		us.Language = "zh-CN"
		us.IsActive = true
		_ = l.userRepo.Save(&us)
	}
	tx.Commit()
	return nil
}
