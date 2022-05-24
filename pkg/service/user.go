package service

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/ldap"
	"github.com/KubeOperator/KubeOperator/pkg/util/message"
	"github.com/jinzhu/gorm"
)

var (
	errOriginalNotMatch  = errors.New("ORIGINAL_NOT_MATCH")
	errUserNotFound      = errors.New("USER_NOT_FOUND")
	errUserIsNotActive   = errors.New("USER_IS_NOT_ACTIVE")
	errUserNameExist     = errors.New("NAME_EXISTS")
	errLdapDisable       = errors.New("LDAP_DISABLE")
	errEmailExist        = errors.New("EMAIL_EXIST")
	errNamePwdFailed     = errors.New("NAME_PASSWORD_SAME_FAILED")
	errEmailDisable      = errors.New("EMAIL_DISABLE")
	errEmailNotMatch     = errors.New("EMAIL_NOT_MATCH")
	errNameOrPasswordErr = errors.New("NAME_PASSWORD_ERROR")
)

type UserService interface {
	Get(name string) (*dto.User, error)
	List(user dto.SessionUser, conditions condition.Conditions) ([]dto.User, error)
	Create(isSuper bool, creation dto.UserCreate) (*dto.User, error)
	Page(num, size int, user dto.SessionUser, conditions condition.Conditions) (*page.Page, error)
	Delete(name string) error
	Update(name string, isSuper bool, update dto.UserUpdate) (*dto.User, error)
	Batch(op dto.UserOp) error
	ChangePassword(isSuper bool, ch dto.UserChangePassword) error
	UserAuth(name string, password string) (user *model.User, err error)
	ResetPassword(fp dto.UserForgotPassword) error
}

type userService struct {
	userRepo                   repository.UserRepository
	userNotificationConfigRepo repository.UserNotificationConfigRepository
	systemService              SystemSettingService
}

func NewUserService() UserService {
	return &userService{
		userRepo:                   repository.NewUserRepository(),
		systemService:              NewSystemSettingService(),
		userNotificationConfigRepo: repository.NewUserNotificationConfigRepository(),
	}
}

func (u *userService) Get(name string) (*dto.User, error) {
	var mo model.User
	if err := db.DB.Where(model.User{Name: name}).
		Preload("CurrentProject").
		First(&mo).Error; err != nil {
		return nil, err
	}
	d := toUserDTO(mo)
	return &d, nil
}

func (u *userService) List(user dto.SessionUser, conditions condition.Conditions) ([]dto.User, error) {
	var userDTOS []dto.User
	var mos []model.User
	d := db.DB.Model(model.User{})
	if err := dbUtil.WithConditions(&d, model.User{}, conditions); err != nil {

		return nil, err
	}
	if user.IsSuper {
		if err := d.Order("name").
			Preload("CurrentProject").
			Find(&mos).Error; err != nil {
			return nil, err
		}
	} else {
		if err := d.Where("is_admin = ? OR name = ?", false, user.Name).
			Order("name").
			Preload("CurrentProject").
			Find(&mos).Error; err != nil {
			return nil, err
		}
	}
	for _, mo := range mos {
		userDTOS = append(userDTOS, toUserDTO(mo))
	}
	return userDTOS, nil
}

func (u *userService) Page(num, size int, user dto.SessionUser, conditions condition.Conditions) (*page.Page, error) {
	var (
		p        page.Page
		userDTOs []dto.User
		mos      []model.User
	)
	d := db.DB.Model(model.User{})
	if err := dbUtil.WithConditions(&d, model.User{}, conditions); err != nil {
		return nil, err
	}
	if user.IsSuper {
		if err := d.
			Count(&p.Total).
			Order("name").
			Offset((num - 1) * size).
			Limit(size).
			Preload("CurrentProject").
			Find(&mos).Error; err != nil {
			return nil, err
		}
	} else {
		if err := d.
			Where("is_admin = ? OR name = ?", false, user.Name).
			Count(&p.Total).
			Order("name").
			Offset((num - 1) * size).
			Limit(size).
			Preload("CurrentProject").
			Find(&mos).Error; err != nil {
			return nil, err
		}
	}

	for _, mo := range mos {
		userDTOs = append(userDTOs, toUserDTO(mo))
	}
	p.Items = userDTOs
	return &p, nil
}

func (u *userService) Create(isSuper bool, creation dto.UserCreate) (*dto.User, error) {

	if creation.Name == creation.Password {
		return nil, errNamePwdFailed
	}

	old, _ := u.Get(creation.Name)
	if old != nil {
		return nil, errUserNameExist
	}

	if creation.Email == "" {
		return nil, errEmailNotMatch
	}
	var userEmail model.User
	db.DB.Where("email = ?", creation.Email).First(&userEmail)
	if userEmail.ID != "" {
		return nil, errEmailExist
	}
	password, err := encrypt.StringEncrypt(creation.Password)
	if err != nil {
		return nil, err
	}
	user := model.User{
		Name:     creation.Name,
		Email:    creation.Email,
		Password: password,
		IsActive: true,
		Language: model.ZH,
		IsAdmin:  strings.ToLower(creation.Role) == constant.SystemRoleAdmin,
		IsSuper:  false,
		Type:     constant.Local,
	}
	if !isSuper {
		user.IsAdmin = false
	}
	err = u.userRepo.Save(&user)
	if err != nil {
		return nil, err
	}

	userNotificationConfig1 := new(model.UserNotificationConfig)
	userNotificationConfig1.UserID = user.ID
	userNotificationConfig1.Type = constant.System
	vars := make(map[string]string)
	vars[constant.LocalMail] = constant.Enable
	vars[constant.Email] = constant.Disable
	vars[constant.DingTalk] = constant.Disable
	vars[constant.WorkWeiXin] = constant.Disable
	data, _ := json.Marshal(vars)
	userNotificationConfig1.Vars = string(data)
	err = u.userNotificationConfigRepo.Save(userNotificationConfig1)
	if err != nil {
		return nil, err
	}
	userNotificationConfig2 := new(model.UserNotificationConfig)
	userNotificationConfig2.UserID = user.ID
	userNotificationConfig2.Vars = string(data)
	userNotificationConfig2.Type = constant.Cluster
	err = u.userNotificationConfigRepo.Save(userNotificationConfig2)
	if err != nil {
		return nil, err
	}

	d := toUserDTO(user)
	return &d, err
}

func (u *userService) Update(name string, isSuper bool, update dto.UserUpdate) (*dto.User, error) {
	var mo model.User
	if err := db.DB.Where(model.User{Name: name}).First(&mo).Error; err != nil {
		return nil, err
	}
	if update.Email != "" {
		mo.Email = update.Email
	}
	if update.Language != "" {
		mo.Language = update.Language
	}

	if isSuper {
		if update.Role != "" {
			mo.IsAdmin = (strings.ToLower(update.Role) == constant.SystemRoleAdmin || strings.ToLower(update.Role) == constant.SystemRoleSuperAdmin)
		}
	}

	if update.Status != "" {
		mo.IsActive = strings.ToLower(update.Status) == constant.UserStatusActive
	}

	if update.CurrentProject != "" {
		var p model.Project
		if err := db.DB.Where(model.Project{Name: update.CurrentProject}).First(&p).Error; err != nil {
			return nil, err
		}
		mo.CurrentProjectID = p.ID
		mo.CurrentProject = p
	}

	if err := db.DB.Save(&mo).Error; err != nil {
		return nil, err
	}
	d := toUserDTO(mo)
	return &d, nil
}

func (u *userService) Delete(name string) error {
	return u.userRepo.Delete(name)
}

func (u *userService) Batch(op dto.UserOp) error {
	var deleteItems []model.User
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.User{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	return u.userRepo.Batch(op.Operation, deleteItems)
}

func (u *userService) ChangePassword(isSuper bool, ch dto.UserChangePassword) error {
	user, err := u.userRepo.Get(ch.Name)
	if err != nil {
		return err
	}
	if !isSuper {
		success, err := user.ValidateOldPassword(ch.Original)
		if err != nil {
			return err
		}
		if !success {
			return errOriginalNotMatch
		}
		if ch.Password == user.Name {
			return errNamePwdFailed
		}
	}

	user.Password, err = encrypt.StringEncrypt(ch.Password)
	if err != nil {
		return err
	}
	err = u.userRepo.Save(&user)
	if err != nil {
		return err
	}
	return err
}

func (u *userService) UserAuth(name string, password string) (user *model.User, err error) {
	var dbUser model.User
	if db.DB.Where("name = ?", name).Preload("CurrentProject").First(&dbUser).RecordNotFound() {
		if db.DB.Where("email = ?", name).Preload("CurrentProject").First(&dbUser).RecordNotFound() {
			return nil, errNameOrPasswordErr
		}
	}
	if !dbUser.IsActive {
		return nil, errUserIsNotActive
	}

	if dbUser.Type == constant.Ldap {
		enable, err := NewSystemSettingService().Get("ldap_status")
		if err != nil {
			return nil, err
		}
		if enable.Value == "DISABLE" {
			return nil, errLdapDisable
		}
		result, err := NewSystemSettingService().List()
		if err != nil {
			return nil, err
		}
		ldapClient, err := ldap.NewLdap(result.Vars)
		if err != nil {
			return nil, err
		}
		err = ldapClient.Connect()
		if err != nil {
			return nil, err
		}
		err = ldapClient.Login(name, password)
		if err != nil {
			return nil, err
		}
	} else {
		uPassword, err := encrypt.StringDecrypt(dbUser.Password)
		if err != nil {
			return nil, err
		}
		if uPassword != password {
			return nil, errNameOrPasswordErr
		}
	}
	return &dbUser, nil
}

func (u *userService) ResetPassword(fp dto.UserForgotPassword) error {
	user, err := u.userRepo.Get(fp.Username)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errUserNotFound
		}
		return err
	}
	if user.Email != fp.Email {
		return errEmailNotMatch
	}
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	password := string(b)
	user.Password, err = encrypt.StringEncrypt(password)
	if err != nil {
		return err
	}
	systemSetting, err := NewSystemSettingService().ListByTab("EMAIL")
	if err != nil {
		return err
	}
	if systemSetting.Vars == nil || systemSetting.Vars["EMAIL_STATUS"] != "ENABLE" {
		return errEmailDisable
	}
	vars := make(map[string]interface{})
	vars["type"] = "EMAIL"
	for k, value := range systemSetting.Vars {
		vars[k] = value
	}
	mClient, err := message.NewMessageClient(vars)
	if err != nil {
		return err
	}
	vars["TITLE"] = "重置密码"
	vars["CONTENT"] = "<html>您好：" + user.Name + "</br>您的密码被重置为" + password + "</html>"
	vars["RECEIVERS"] = fp.Email
	err = mClient.SendMessage(vars)
	if err != nil {
		return err
	}
	err = u.userRepo.Save(&user)
	if err != nil {
		return err
	}
	return nil
}

func toUserDTO(user model.User) dto.User {
	u := dto.User{User: user}
	u.Role = func() string {
		if u.IsAdmin {
			if u.IsSuper {
				return constant.SystemRoleSuperAdmin
			} else {
				return constant.SystemRoleAdmin
			}
		}
		return constant.SystemRoleUser
	}()
	u.Status = func() string {
		if u.IsActive {
			return constant.UserStatusActive
		}
		return constant.UserStatusPassive
	}()
	u.CurrentProject = user.CurrentProject.Name
	return u
}
