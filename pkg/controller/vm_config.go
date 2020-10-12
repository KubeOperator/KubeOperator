package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/service"
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
// @Tags vmConfig
// @Summary Show all vmConfigs
// @Description Show vmConfigs
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /vmConfigs/ [get]
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
