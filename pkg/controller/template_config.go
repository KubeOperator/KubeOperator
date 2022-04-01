package controller

import (
	"context"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
)

type TemplateConfigController struct {
	Ctx                   context.Context
	TemplateConfigService service.TemplateConfigService
}

func NewTemplateConfigController() *TemplateConfigController {
	return &TemplateConfigController{
		TemplateConfigService: service.NewTemplateConfigService(),
	}
}

// List TemplateConfigs
// @Tags TemplateConfig
// @Summary Show all TemplateConfigs
// @Description 获取所有的模板配置
// @Accept  json
// @Produce  json
// @Success 200 {object} []dto.TemplateConfig
// @Security ApiKeyAuth
// @Router /templateconfigs [get]
func (t TemplateConfigController) Get() ([]dto.TemplateConfig, error) {
	return t.TemplateConfigService.List()
}
