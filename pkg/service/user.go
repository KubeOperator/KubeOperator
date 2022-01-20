package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/ldap"
	"github.com/KubeOperator/KubeOperator/pkg/util/message"
	"github.com/jinzhu/gorm"
)

var (
	OriginalNotMatch  = errors.New("ORIGINAL_NOT_MATCH")
	UserNotFound      = errors.New("USER_NOT_FOUND")
	UserIsNotActive   = errors.New("USER_IS_NOT_ACTIVE")
	UserNameExist     = errors.New("NAME_EXISTS")
	LdapDisable       = errors.New("LDAP_DISABLE")
	EmailExist        = errors.New("EMAIL_EXIST")
	NamePwdFailed     = errors.New("NAME_PASSWORD_SAME_FAILED")
	EmailDisable      = errors.New("EMAIL_DISABLE")
	EmailNotMatch     = errors.New("EMAIL_NOT_MATCH")
	NameOrPasswordErr = errors.New("NAME_PASSWORD_ERROR")
)

type UserService interface {
	Get(name string) (dto.User, error)
	List() ([]dto.User, error)
	Create(creation dto.UserCreate) (*dto.User, error)
	Page(num, size int) (page.Page, error)
	Delete(name string) error
	Batch(op dto.UserOp) error
	ChangePassword(ch dto.UserChangePassword) error
	UserAuth(name string, password string) (user *model.User, err error)
	ResetPassword(fp dto.UserForgotPassword) error
}

type userService struct {
	userRepo      repository.UserRepository
	systemService SystemSettingService
}

func NewUserService() UserService {
	return &userService{
		userRepo:      repository.NewUserRepository(),
		systemService: NewSystemSettingService(),
	}
}

func (u userService) Get(name string) (dto.User, error) {
	var userDTO dto.User
	mo, err := u.userRepo.Get(name)
	if err != nil {
		return userDTO, err
	}
	userDTO = dto.User{
		ID:        mo.ID,
		Name:      mo.Name,
		Email:     mo.Email,
		IsActive:  mo.IsActive,
		Language:  mo.Language,
		IsAdmin:   mo.IsAdmin,
		Type:      mo.Type,
		CreatedAt: mo.CreatedAt,
	}
	return userDTO, err
}

func (u userService) List() ([]dto.User, error) {
	var userDTOS []dto.User
	mos, err := u.userRepo.List()
	if err != nil {
		return userDTOS, err
	}
	for _, mo := range mos {
		user := dto.User{
			ID:        mo.ID,
			Name:      mo.Name,
			Email:     mo.Email,
			IsActive:  mo.IsActive,
			Language:  mo.Language,
			IsAdmin:   mo.IsAdmin,
			Type:      mo.Type,
			CreatedAt: mo.CreatedAt,
		}
		userDTOS = append(userDTOS, user)
	}
	return userDTOS, err
}

func (u userService) Create(creation dto.UserCreate) (*dto.User, error) {

	if creation.Name == creation.Password {
		return nil, NamePwdFailed
	}

	old, _ := u.Get(creation.Name)
	if old.ID != "" {
		return nil, UserNameExist
	}

	var userEmail model.User
	db.DB.Where("email = ?", creation.Email).First(&userEmail)
	if userEmail.ID != "" {
		return nil, EmailExist
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
		IsAdmin:  creation.IsAdmin,
		Type:     constant.Local,
	}
	err = u.userRepo.Save(&user)
	if err != nil {
		return nil, err
	}

	userDTO := dto.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		Language:  user.Language,
		IsAdmin:   user.IsAdmin,
		Type:      user.Type,
		CreatedAt: user.CreatedAt,
	}

	return &userDTO, err
}

func (u userService) Page(num, size int) (page.Page, error) {
	var page page.Page
	var userDTOs []dto.User
	total, mos, err := u.userRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		user := dto.User{
			ID:        mo.ID,
			Name:      mo.Name,
			Email:     mo.Email,
			IsActive:  mo.IsActive,
			Language:  mo.Language,
			IsAdmin:   mo.IsAdmin,
			Type:      mo.Type,
			CreatedAt: mo.CreatedAt,
		}
		userDTOs = append(userDTOs, user)
	}
	page.Total = total
	page.Items = userDTOs
	return page, err
}

func (u userService) Delete(name string) error {
	return u.userRepo.Delete(name)
}

func (u userService) Batch(op dto.UserOp) error {
	var deleteItems []model.User
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.User{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	return u.userRepo.Batch(op.Operation, deleteItems)
}

func (u userService) ChangePassword(ch dto.UserChangePassword) error {
	user, err := u.userRepo.Get(ch.Name)
	if err != nil {
		return err
	}
	success, err := user.ValidateOldPassword(ch.Original)
	if err != nil {
		return err
	}
	if !success {
		return OriginalNotMatch
	}
	if ch.Password == user.Name {
		return NamePwdFailed
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

func (u userService) UserAuth(name string, password string) (user *model.User, err error) {
	var dbUser model.User
	if db.DB.Where("name = ?", name).First(&dbUser).RecordNotFound() {
		if db.DB.Where("email = ?", name).First(&dbUser).RecordNotFound() {
			return nil, NameOrPasswordErr
		}
	}
	if !dbUser.IsActive {
		return nil, UserIsNotActive
	}

	if dbUser.Type == constant.Ldap {
		enable, err := NewSystemSettingService().Get("ldap_status")
		if err != nil {
			return nil, err
		}
		if enable.Value == "DISABLE" {
			return nil, LdapDisable
		}
		result, err := NewSystemSettingService().List()
		if err != nil {
			return nil, err
		}
		ldapClient := ldap.NewLdap(result.Vars)
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
			return nil, NameOrPasswordErr
		}
	}
	return &dbUser, nil
}

func (u userService) ResetPassword(fp dto.UserForgotPassword) error {
	user, err := u.userRepo.Get(fp.Username)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return UserNotFound
		}
		return err
	}
	if user.Email != fp.Email {
		return EmailNotMatch
	}

	password := GetPasswd()
	user.Password, err = encrypt.StringEncrypt(password)
	if err != nil {
		return err
	}
	systemSetting, err := NewSystemSettingService().ListByTab("EMAIL")
	if err != nil {
		return err
	}
	if systemSetting.Vars == nil || systemSetting.Vars["EMAIL_STATUS"] != "ENABLE" {
		return EmailDisable
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
	vars["CONTENT"] = user.Name + "您的密码被重置为" + password
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

func GetPasswd() string {
	rand.Seed(time.Now().UnixNano())
	lenNum := rand.Intn(5)

	var numbers = []rune("0123456789")
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 8)
	for j := 0; j < lenNum; j++ {
		b[j] = numbers[rand.Intn(len(numbers))]
	}
	for k := lenNum; k < 8; k++ {
		b[k] = letters[rand.Intn(len(letters))]
	}
	return (string(b))
}
