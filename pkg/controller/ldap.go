package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type LdapController struct {
	Ctx         context.Context
	LdapService service.LdapService
}

func NewLdapController() *LdapController {
	return &LdapController{
		LdapService: service.NewLdapService(),
	}
}

func (l LdapController) Post() ([]dto.SystemSetting, error) {
	var req dto.SystemSettingCreate
	err := l.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	result, err := l.LdapService.Create(req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l *LdapController) PostTestConnect() (*dto.LdapResult, error) {
	ctx := l.Ctx
	var req dto.SystemSettingCreate
	if err := ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	users, err := l.LdapService.TestConnect(req)
	if err != nil {
		return nil, err
	}
	return &dto.LdapResult{Data: users}, nil
}

func (l *LdapController) PostTestLogin() error {
	ctx := l.Ctx
	var req dto.LdapLogin
	if err := ctx.ReadJSON(&req); err != nil {
		return err
	}
	return l.LdapService.TestLogin(req.Username, req.Password)
}

func (l *LdapController) GetSync() ([]dto.LdapUser, error) {
	users, err := l.LdapService.LdapSync()
	return users, err
}

func (l *LdapController) PostImportUsers() error {
	ctx := l.Ctx
	var req dto.ImportRequest
	if err := ctx.ReadJSON(&req); err != nil {
		return err
	}
	return l.LdapService.ImpostUsers(req.Users)
}
