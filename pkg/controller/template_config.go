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
// @Tags templateConfigs
// @Summary Show all TemplateConfigs
// @Description 获取所有的模板配置
// @Accept  json
// @Produce  json
// @Success 200 {object} []dto.TemplateConfig
// @Security ApiKeyAuth
// @Router /templates/search [post]
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

func (t TemplateConfigController) Get() (*page.Page, error) {
	p, _ := t.Ctx.Values().GetBool("page")
	if p {
		num, _ := t.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := t.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return t.TemplateConfigService.Page(num, size, condition.TODO())
	} else {
		var page page.Page
		items, err := t.TemplateConfigService.List()
		if err != nil {
			return nil, err
		}
		page.Items = items
		page.Total = len(items)
		return &page, nil
	}

}

// Update TemplateConfig
// @Tags templateConfigs
// @Summary Update a TemplateConfig
// @Description 更新模板配置
// @Accept  json
// @Produce  json
// @Param request body dto.TemplateConfig true "request"
// @Param name path string true "模板配置名称"
// @Success 200 {object} dto.TemplateConfig
// @Security ApiKeyAuth
// @Router /templates/{name} [patch]
func (t TemplateConfigController) PatchBy(name string) (*dto.TemplateConfig, error) {
	var req dto.TemplateConfig
	err := t.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return t.TemplateConfigService.Update(name, req)
}

// Create TemplateConfig
// @Tags templateConfigs
// @Summary Create a TemplateConfig
// @Description 创建模版
// @Accept  json
// @Produce  json
// @Param request body dto.TemplateConfigCreate true "request"
// @Success 200 {object} dto.TemplateConfig
// @Security ApiKeyAuth
// @Router /templates/create [post]

func (t TemplateConfigController) PostCreate() (*dto.TemplateConfig, error) {

	var req dto.TemplateConfigCreate
	if err := t.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}
	operator := t.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_TEMPLATE, req.Name)

	return t.TemplateConfigService.Create(req)
}

// Get TemplateConfig
// @Tags templateConfigs
// @Summary Show a TemplateConfig
// @Description 获取单个模版配置
// @Accept  json
// @Produce  json
// @Param name path string true "模版名称"
// @Success 200 {object} dto.TemplateConfig
// @Security ApiKeyAuth
// @Router /templates/{name} [get]
func (t TemplateConfigController) GetBy(name string) (*dto.TemplateConfig, error) {
	return t.TemplateConfigService.Get(name)
}

// Delete TemplateConfig
// @Tags templateConfigs
// @Summary Delete a templateConfig
// @Description 删除模版配置
// @Accept  json
// @Produce  json
// @Param name path string true "模版名称"
// @Security ApiKeyAuth
// @Router /templates/{name} [delete]
func (t TemplateConfigController) DeleteBy(name string) error {
	return t.TemplateConfigService.Delete(name)
}
