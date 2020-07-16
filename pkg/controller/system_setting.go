package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type SystemSettingController struct {
	Ctx                  context.Context
	SystemSettingService service.SystemSettingService
}

func NewSystemSettingController() *SystemSettingController {
	return &SystemSettingController{
		SystemSettingService: service.NewSystemSettingService(),
	}
}

func (s SystemSettingController) Get() (page.Page, error) {

	var page page.Page
	items, err := s.SystemSettingService.List()
	if err != nil {
		return page, err
	}
	page.Items = items
	page.Total = len(items)
	return page, nil
}

func (s SystemSettingController) Post() ([]dto.SystemSetting, error) {
	var req dto.SystemSettingCreate
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return s.SystemSettingService.Create(req)
}
