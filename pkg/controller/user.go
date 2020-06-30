package controller

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/controller/warp"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
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

func (u UserController) Get() (page.Page, error) {

	p, _ := u.Ctx.Values().GetBool("page")
	if p {
		num, _ := u.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := u.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return u.UserService.Page(num, size)
	} else {
		var page page.Page
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

func (u UserController) PostBatch() error {
	var req dto.UserOp
	err := u.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = u.UserService.Batch(req)
	if err != nil {
		return warp.NewControllerError(errors.New(u.Ctx.Tr(err.Error())))
	}
	return err
}

func (u UserController) PostChangePassword() error {
	var req dto.UserChangePassword
	err := u.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = u.UserService.ChangePassword(req)
	if err != nil {
		return warp.NewControllerError(errors.New(u.Ctx.Tr(err.Error())))
	}
	return err
}
