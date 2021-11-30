package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
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
// @Summary Show users
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /users [get]
func (u *UserController) Get() (*page.Page, error) {

	p, _ := u.Ctx.Values().GetBool("page")
	sessionUser := u.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)
	if p {
		num, _ := u.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := u.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return u.UserService.Page(num, size, user, condition.TODO())
	} else {
		var p page.Page
		items, err := u.UserService.List(user, condition.TODO())
		if err != nil {
			return &p, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Search User
// @Tags users
// @Summary Search user
// @Description 过滤用户
// @Accept  json
// @Produce  json
// @Param conditions body condition.Conditions true "conditions"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /users/search [post]
func (u *UserController) PostSearch() (*page.Page, error) {
	var conditions condition.Conditions
	if u.Ctx.GetContentLength() > 0 {
		if err := u.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	p, _ := u.Ctx.Values().GetBool("page")
	sessionUser := u.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)

	if p {
		num, _ := u.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := u.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return u.UserService.Page(num, size, user, conditions)
	} else {
		var p page.Page
		items, err := u.UserService.List(user, conditions)
		if err != nil {
			return &p, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Get User
// @Tags users
// @Summary Show a user
// @Description 获取单个用户
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.User
// @Security ApiKeyAuth
// @Router /users/{name} [get]
func (u *UserController) GetBy(name string) (*dto.User, error) {
	return u.UserService.Get(name)
}

// Create User
// @Tags users
// @Summary Create a user
// @Description 创建用户
// @Accept  json
// @Produce  json
// @Param request body dto.UserCreate true "request"
// @Success 200 {object} dto.User
// @Security ApiKeyAuth
// @Router /users [post]
func (u *UserController) Post() (*dto.User, error) {
	var req dto.UserCreate
	err := u.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}

	sessionUser := u.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)

	operator := u.Ctx.Values().GetString("operator")
	kolog.Save(operator, constant.CREATE_USER, req.Name)

	return u.UserService.Create(user.IsSuper, req)
}

// Delete User
// @Tags users
// @Summary Delete a user
// @Description 删除用户
// @Accept  json
// @Produce  json
// @Param name path string true "用户名"
// @Security ApiKeyAuth
// @Router /users/{name} [delete]
func (u *UserController) DeleteBy(name string) error {
	operator := u.Ctx.Values().GetString("operator")
	kolog.Save(operator, constant.DELETE_USER, name)

	return u.UserService.Delete(name)
}

// Update User
// @Tags users
// @Summary Update a user
// @Description 更新用户
// @Accept  json
// @Produce  json
// @Param request body dto.UserUpdate true "request"
// @Param name path string true "用户名"
// @Success 200 {object} dto.User
// @Security ApiKeyAuth
// @Router /users/{name} [patch]
func (u *UserController) PatchBy(name string) (*dto.User, error) {
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

	sessionUser := u.Ctx.Values().Get("user")
	sessions, _ := sessionUser.(dto.SessionUser)

	user, err := u.UserService.Update(name, sessions.IsSuper, req)
	if err != nil {
		return nil, err
	}

	operator := u.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_USER, name)

	return user, err
}

func (u *UserController) PostBatch() error {
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
		delUser += userItem.Name + ","
	}
	go kolog.Save(operator, constant.DELETE_USER, delUser)

	return err
}

// Change User Password
// @Tags users
// @Summary Change user password
// @Description 更新用户密码
// @Accept  json
// @Produce  json
// @Param request body dto.UserChangePassword true "request"
// @Success 200 {object} dto.User
// @Security ApiKeyAuth
// @Router /users/change/password [post]
func (u *UserController) PostChangePassword() error {
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
	sessionUser := u.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)

	err = u.UserService.ChangePassword(user.IsSuper, req)
	if err != nil {
		return err
	}

	operator := u.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_USER_PASSWORD, req.Name)

	return err
}
