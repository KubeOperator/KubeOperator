package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
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
// @Router /template/search [post]
func (t TemplateConfigController) PostSearch() (*page.Page, error) {

	p, _ := t.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	if t.Ctx.GetContentLength() > 0 {
		if err := t.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	if p {
		num, _ := t.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := t.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return t.TemplateConfigService.Page(num, size, conditions)
	} else {
		var p page.Page
		items, err := t.TemplateConfigService.List()
		if err != nil {
			return nil, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

func (t TemplateConfigController) Delete() error {
	return nil
}

func (t TemplateConfigController) Update() (dto.TemplateConfig, error) {
	return dto.TemplateConfig{}, nil
}

func (t TemplateConfigController) Create() (dto.TemplateConfig, error) {
	return dto.TemplateConfig{}, nil
}
