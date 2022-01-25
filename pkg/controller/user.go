package controller

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/koregexp"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/session"
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
// @Param  pageNum  query  int  true "page number"
// @Param  pageSize  query  int  true "page size"
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
	validate := validator.New()
	if err := validate.RegisterValidation("kopassword", koregexp.CheckPasswordPattern); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, err
	}
	if req.Password == reverseString(req.Name) || req.Password == req.Name {
		return nil, errors.New("NAME_PASSWORD_SAME_FAILED")
	}

	operator := u.Ctx.Values().GetString("operator")
	kolog.Save(operator, constant.CREATE_USER, req.Name)

	return u.UserService.Create(req)
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
	user, err := u.UserService.Update(req)
	if err != nil {
		return nil, err
	}

	var onlineSessionIDList = session.GloablSessionMgr.GetSessionIDList()
	for _, onlineSessionID := range onlineSessionIDList {
		if userInfo, ok := session.GloablSessionMgr.GetSessionVal(onlineSessionID, constant.SessionUserKey); ok {
			if value, ok := userInfo.(*dto.Profile); ok {
				if value.User.Name == req.Name {
					session.GloablSessionMgr.EndSessionBy(onlineSessionID)
				}
			}
		}
	}

	operator := u.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_USER, name)

	return user, err
}

func (u UserController) DeleteBy(name string) error {
	operator := u.Ctx.Values().GetString("operator")

	var onlineSessionIDList = session.GloablSessionMgr.GetSessionIDList()
	for _, onlineSessionID := range onlineSessionIDList {
		if userInfo, ok := session.GloablSessionMgr.GetSessionVal(onlineSessionID, constant.SessionUserKey); ok {
			if value, ok := userInfo.(*dto.Profile); ok {
				if value.User.Name == name {
					session.GloablSessionMgr.EndSessionBy(onlineSessionID)
				}
			}
		}
	}

	kolog.Save(operator, constant.DELETE_USER, name)

	return u.UserService.Delete(name)
}

// Delete Users
// @Tags users
// @Summary Delete user list
// @Description delete user list
// @Accept  json
// @Produce  json
// @Param request body dto.UserOp true "request"
// @Security ApiKeyAuth
// @Router /users/batch [post]
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

	var onlineSessionIDList = session.GloablSessionMgr.GetSessionIDList()
	for _, onlineSessionID := range onlineSessionIDList {
		if userInfo, ok := session.GloablSessionMgr.GetSessionVal(onlineSessionID, constant.SessionUserKey); ok {
			if value, ok := userInfo.(*dto.Profile); ok {
				for _, userItem := range req.Items {
					if value.User.Name == userItem.Name {
						session.GloablSessionMgr.EndSessionBy(onlineSessionID)
					}
				}
			}
		}
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
	if err := validate.RegisterValidation("kopassword", koregexp.CheckPasswordPattern); err != nil {
		return err
	}
	if err := validate.Struct(req); err != nil {
		return err
	}

	if req.Password == reverseString(req.Name) || req.Password == req.Name || req.Password == req.Original {
		return errors.New("NAME_PASSWORD_SAME_FAILED")
	}

	err = u.UserService.ChangePassword(req)
	if err != nil {
		return err
	}

	var onlineSessionIDList = session.GloablSessionMgr.GetSessionIDList()
	for _, onlineSessionID := range onlineSessionIDList {
		if userInfo, ok := session.GloablSessionMgr.GetSessionVal(onlineSessionID, constant.SessionUserKey); ok {
			if value, ok := userInfo.(*dto.Profile); ok {
				if value.User.Name == req.Name {
					session.GloablSessionMgr.EndSessionBy(onlineSessionID)
				}
			}
		}
	}

	operator := u.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_USER_PASSWORD, req.Name)

	return err
}

func reverseString(s string) string {
	runes := []rune(s)

	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}

	return string(runes)
}
