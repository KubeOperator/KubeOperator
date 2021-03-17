package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
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
// @Router /logs/ [get]
func (u SystemLogController) Post() (page.Page, error) {
	num, _ := u.Ctx.Values().GetInt(constant.PageNumQueryKey)
	size, _ := u.Ctx.Values().GetInt(constant.PageSizeQueryKey)

	var page page.Page
	var queryCondition dto.SystemLogQuery
	err := u.Ctx.ReadJSON(&queryCondition)
	if err != nil {
		return page, err
	}

	return u.SystemLogService.Page(num, size, queryCondition)
}
