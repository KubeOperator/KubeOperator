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

func (l LdapController) PostSync() error {
	var req dto.SystemSettingCreate
	err := l.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = l.LdapService.LdapSync(req)
	if err != nil {
		return err
	}
	return nil
}
