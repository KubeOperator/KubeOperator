package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type VmConfigController struct {
	Ctx             context.Context
	VmConfigService service.VmConfigService
}

func NewVmConfigController() *VmConfigController {
	return &VmConfigController{
		VmConfigService: service.NewVmConfigService(),
	}
}

// List VmConfigs
// @Tags vmConfigs
// @Summary Show all vmConfigs
// @Description Show vmConfigs
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /vm/configs/ [get]
func (v VmConfigController) Get() (page.Page, error) {
	pa, _ := v.Ctx.Values().GetBool("page")
	if pa {
		num, _ := v.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := v.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return v.VmConfigService.Page(num, size)
	} else {
		var page page.Page
		items, err := v.VmConfigService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

// Create VmConfig
// @Tags vmConfigs
// @Summary Create a VmConfig
// @Description create a VmConfig
// @Accept  json
// @Produce  json
// @Param request body dto.VmConfigCreate true "request"
// @Success 200 {object} dto.VmConfig
// @Security ApiKeyAuth
// @Router /vm/config/ [post]
func (v VmConfigController) Post() (*dto.VmConfig, error) {
	var req dto.VmConfigCreate
	err := v.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return v.VmConfigService.Create(req)
}

func (v VmConfigController) PostBatch() error {
	var req dto.VmConfigOp
	err := v.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = v.VmConfigService.Batch(req)
	if err != nil {
		return err
	}
	return err
}
