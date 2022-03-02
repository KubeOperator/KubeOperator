package service

import (
	"errors"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/ldap"
	"github.com/jinzhu/gorm"
)

var (
	OriginalNotMatch  = errors.New("ORIGINAL_NOT_MATCH")
	UserNotFound      = errors.New("USER_NOT_FOUND")
	UserIsNotActive   = errors.New("USER_IS_NOT_ACTIVE")
	UserNameExist     = errors.New("NAME_EXISTS")
	LdapDisable       = errors.New("LDAP_DISABLE")
	NamePwdFailed     = errors.New("NAME_PASSWORD_SAME_FAILED")
	EmailDisable      = errors.New("EMAIL_DISABLE")
	EmailNotMatch     = errors.New("EMAIL_NOT_MATCH")
	NameOrPasswordErr = errors.New("NAME_PASSWORD_ERROR")
)

type UserService interface {
	Get(name string) (dto.User, error)
	List() ([]dto.User, error)
	Create(creation dto.UserCreate) (*dto.User, error)
	Update(update dto.UserUpdate) (*dto.User, error)
	Page(num, size int) (page.Page, error)
	Delete(name string) error
	Batch(op dto.UserOp) error
	ChangePassword(ch dto.UserChangePassword) (bool, error)
	UserAuth(name string, password string, isSystem bool) (user *model.User, err error)
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

	old, err := u.Get(creation.Name)
	if !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if old.ID != "" {
		return nil, UserNameExist
	}

	salt, password := encrypt.Encode(creation.Password, nil)
	user := model.User{
		Name:     creation.Name,
		Salt:     salt,
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

func (u userService) Update(update dto.UserUpdate) (*dto.User, error) {
	old, err := u.Get(update.Name)
	if err != nil {
		return nil, err
	}

	if err := db.DB.Model(&model.User{}).Where("id = ?", old.ID).Updates(map[string]interface{}{"IsActive": update.IsActive, "IsAdmin": update.IsAdmin}).Error; err != nil {
		return nil, err
	}

	userDTO := dto.User{
		ID:        old.ID,
		Name:      old.Name,
		IsActive:  old.IsActive,
		Language:  old.Language,
		IsAdmin:   old.IsAdmin,
		Type:      old.Type,
		CreatedAt: old.CreatedAt,
	}

	return &userDTO, err
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

func (u userService) ChangePassword(ch dto.UserChangePassword) (bool, error) {
	isFirstLogin := false
	user, err := u.userRepo.Get(ch.Name)
	if err != nil {
		return isFirstLogin, err
	}
	if user.IsFirst {
		isFirstLogin = true
		user.IsFirst = false
	}

	success, err := validateOldPassword(user, ch.Original)
	if !success || err != nil {
		return isFirstLogin, err
	}

	user.Salt, user.Password = encrypt.Encode(ch.Password, nil)
	if err = u.userRepo.Save(&user); err != nil {
		return isFirstLogin, err
	}
	return isFirstLogin, err
}

func (u userService) UserAuth(name string, password string, isSystem bool) (user *model.User, err error) {
	var dbUser model.User
	if err := db.DB.Where("name = ? AND is_system = ?", name, isSystem).First(&dbUser).Error; err != nil {
		return nil, err
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
		success, err := validateOldPassword(dbUser, password)
		if !success || err != nil {
			return nil, err
		}
	}
	return &dbUser, nil
}

func validateOldPassword(user model.User, password string) (bool, error) {
	if !user.UpdatedAt.Before(time.Now().Add(-1*time.Minute)) && user.ErrCount > 4 {
		return false, errors.New("TOO_MANY_FAILURES")
	}

	if !encrypt.Verify(password, user.Salt, user.Password, nil) {
		if user.UpdatedAt.Before(time.Now().Add(-1 * time.Minute)) {
			_ = db.DB.Model(&model.User{}).Where("id = ?", user.ID).Update("err_count", 1)
		} else {
			_ = db.DB.Model(&model.User{}).Where("id = ?", user.ID).Update("err_count", gorm.Expr("err_count + 1"))
		}
		return false, NameOrPasswordErr
	}
	return true, nil
}
