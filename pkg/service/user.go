package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/i18n"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
)

var (
	UserNotFound     = errors.New(i18n.Tr("user_not_found", nil))
	PasswordNotMatch = errors.New(i18n.Tr("password_not_match", nil))
	UserIsNotActive  = errors.New(i18n.Tr("user_is_not_active", nil))
)

type UserService interface {
	Get(name string) (dto.User, error)
	List() ([]dto.User, error)
	Create(creation dto.UserCreate) (dto.User, error)
	Page(num, size int) (dto.UserPage, error)
	Delete(name string) error
	Update(update dto.UserUpdate) (dto.User, error)
	//Batch(operation string, items []model.User) ([]model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService() UserService {
	return &userService{
		userRepo: repository.NewUserRepository(),
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

func (u userService) Create(creation dto.UserCreate) (dto.User, error) {

	user := model.User{
		Name:     creation.Name,
		Email:    creation.Email,
		Password: creation.Password,
		IsActive: true,
		Language: model.ZH,
	}
	err := u.userRepo.Save(&user)
	if err != nil {
		return dto.User{}, err
	}
	return dto.User{User: user}, err
}

func (u userService) Update(update dto.UserUpdate) (dto.User, error) {
	user := model.User{
		ID:       update.ID,
		Name:     update.Name,
		Email:    update.Email,
		IsActive: update.IsActive,
		Language: update.Language,
	}
	err := u.userRepo.Save(&user)
	if err != nil {
		return dto.User{}, err
	}
	return dto.User{User: user}, err
}

func (u userService) Page(num, size int) (dto.UserPage, error) {

	var page dto.UserPage
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

//func (u userService) Batch(operation string, items []dto.User) ([]dto.User, error) {
//	var deleteItems []model.User
//	var notOpItems []model.User
//	switch operation {
//	case constant.BatchOperationDelete:
//		tx := db.DB.Begin()
//		for _, item := range items {
//			err := db.DB.Model(model.User{}).First(&item).Delete(&item).Error
//			if err != nil {
//				tx.Rollback()
//				return nil, err
//			}
//			deleteItems = append(deleteItems, item)
//			tx.Commit()
//		}
//	default:
//		return nil, constant.NotSupportedBatchOperation
//	}
//	return deleteItems, nil
//}

func UserAuth(name string, password string) (sessionUser *auth.SessionUser, err error) {
	var dbUser model.User
	if db.DB.Where("name = ?", name).First(&dbUser).RecordNotFound() {
		return nil, UserNotFound
	}
	if dbUser.IsActive == false {
		return dbUser.ToSessionUser(), UserIsNotActive
	}
	password, err = encrypt.StringEncrypt(password)
	if err != nil {
		return nil, err
	}
	if dbUser.Password != password {
		return nil, PasswordNotMatch
	}
	return dbUser.ToSessionUser(), nil
}
