package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/ldap"
)

var (
	PasswordNotMatch = errors.New("PASSWORD_NOT_MATCH")
	OriginalNotMatch = errors.New("ORIGINAL_NOT_MATCH")
	UserNotFound     = errors.New("USER_NOT_FOUND")
	UserIsNotActive  = errors.New("USER_IS_NOT_ACTIVE")
	UserNameExist    = errors.New("NAME_EXISTS")
	LdapDisable      = errors.New("LDAP_DISABLE")
	NamePwdFailed    = errors.New("NAME_PASSWORD_SAME_FAILED")
)

type UserService interface {
	Get(name string) (dto.User, error)
	List() ([]dto.User, error)
	Create(creation dto.UserCreate) (*dto.User, error)
	Page(num, size int) (page.Page, error)
	Delete(name string) error
	Update(update dto.UserUpdate) (*dto.User, error)
	Batch(op dto.UserOp) error
	ChangePassword(ch dto.UserChangePassword) error
	UserAuth(name string, password string) (user *model.User, err error)
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
	userDTO.User = mo
	return userDTO, err
}

func (u userService) List() ([]dto.User, error) {

	var userDTOS []dto.User
	mos, err := u.userRepo.List()
	if err != nil {
		return userDTOS, err
	}
	for _, mo := range mos {
		userDTOS = append(userDTOS, dto.User{User: mo})
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
	return &dto.User{User: user}, err
}

func (u userService) Update(update dto.UserUpdate) (*dto.User, error) {

	old, err := u.Get(update.Name)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:       old.ID,
		Name:     update.Name,
		Email:    update.Email,
		IsActive: update.IsActive,
		IsAdmin:  update.IsAdmin,
		Type:     constant.Local,
		Password: old.Password,
		Language: model.ZH,
	}
	err = u.userRepo.Save(&user)
	if err != nil {
		return nil, err
	}
	return &dto.User{User: user}, err
}

func (u userService) Page(num, size int) (page.Page, error) {

	var page page.Page
	var userDTOs []dto.User
	total, mos, err := u.userRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		userDTOs = append(userDTOs, dto.User{User: mo})
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
			return nil, UserNotFound
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
		password, err = encrypt.StringEncrypt(password)
		if err != nil {
			return nil, err
		}
		if dbUser.Password != password {
			return nil, PasswordNotMatch
		}
	}
	return &dbUser, nil
}
