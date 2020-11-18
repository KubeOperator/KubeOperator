package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type SystemLogController struct {
	Ctx              context.Context
	SystemLogService service.SystemLogService
}

func NewSystemLogController() *SystemLogController {
	return &SystemLogController{
		SystemLogService: service.NewSystemLogService(),
	}
}

// List SystemLog
// @Tags system_logs
// @Summary Show all system_logs
// @Description Show system_logs
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /system_logs/ [get]
func (u SystemLogController) Get() (page.Page, error) {
	p, _ := u.Ctx.Values().GetBool("page")
	if p {
		num, _ := u.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := u.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return u.SystemLogService.Page(num, size)
	} else {
		var page page.Page
		items, err := u.SystemLogService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}
