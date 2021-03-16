package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
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

// List User
// @Tags users
// @Summary Show all users
// @Description Show users
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /users/ [get]
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

// Get User
// @Tags users
// @Summary Show a user
// @Description show a user by name
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.User
// @Security ApiKeyAuth
// @Router /users/{name}/ [get]
func (u UserController) GetBy(name string) (dto.User, error) {
	return u.UserService.Get(name)
}

// Create User
// @Tags users
// @Summary Create a user
// @Description create a user
// @Accept  json
// @Produce  json
// @Param request body dto.UserCreate true "request"
// @Success 200 {object} dto.Host
// @Security ApiKeyAuth
// @Router /users/ [post]
func (u UserController) Post() (*dto.User, error) {
	var req dto.UserCreate
	err := u.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}

	operator := u.Ctx.Values().GetString("operator")
	kolog.Save(operator, constant.CREATE_USER, req.Name)

	return u.UserService.Create(req)
}

// Delete User
// @Tags users
// @Summary Delete a user
// @Description delete a user by name
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /users/{name}/ [delete]
func (u UserController) DeleteBy(name string) error {
	operator := u.Ctx.Values().GetString("operator")
	kolog.Save(operator, constant.DELETE_USER, name)

	return u.UserService.Delete(name)
}

// Update User
// @Tags users
// @Summary Update a user
// @Description Update a user
// @Accept  json
// @Produce  json
// @Param request body dto.UserUpdate true "request"
// @Success 200 {object} dto.User
// @Security ApiKeyAuth
// @Router /users/{name}/ [patch]
func (u UserController) PatchBy(name string) (*dto.User, error) {
	var req dto.UserUpdate
	err := u.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	user, err := u.UserService.Update(name, req)
	if err != nil {
		return nil, err
	}

	operator := u.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_USER, name)

	return user, err
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
		return err
	}

	operator := u.Ctx.Values().GetString("operator")
	delUser := ""
	for _, userItem := range req.Items {
		delUser += (userItem.Name + ",")
	}
	go kolog.Save(operator, constant.DELETE_USER, delUser)

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
		return err
	}

	operator := u.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_USER_PASSWORD, req.Name)

	return err
}
