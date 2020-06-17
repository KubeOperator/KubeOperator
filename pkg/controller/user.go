package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type UserController struct {
	Ctx         context.Context
	UserService service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		UserService: service.NewUserService(),
	}
}

func (u UserController) Get() (dto.UserPage, error) {

	page, _ := u.Ctx.Values().GetBool("page")
	if page {
		num, _ := u.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := u.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return u.UserService.Page(num, size)
	} else {
		var page dto.UserPage
		items, err := u.UserService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

func (u UserController) GetBy(name string) (dto.User, error) {
	return u.UserService.Get(name)
}

func (u UserController) Post() (dto.User, error) {
	var req dto.UserCreate
	err := u.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.User{}, err
	}
	return u.UserService.Create(req)
}

func (u UserController) Delete(name string) error {
	return u.UserService.Delete(name)
}

func (u UserController) PatchBy(name string) (dto.User, error) {
	var req dto.UserUpdate
	err := u.Ctx.ReadJSON(&req)
	if err != nil {
		return dto.User{}, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return dto.User{}, err
	}
	return u.UserService.Update(req)
}

//func (u userController) Batch(operation string, items []dto.User) error {
//	_, err := u.userService.Batch(operation, items)
//	return err
//}
