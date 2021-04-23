package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
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

// Search SystemLog
// @Tags system_logs
// @Summary Search user
// @Description 过滤系统日志
// @Accept  json
// @Produce  json
// @Param conditions body condition.Conditions true "conditions"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /logs/ [post]
func (u SystemLogController) Post() (*page.Page, error) {
	p, _ := u.Ctx.Values().GetBool("page")
	if p {
		var conditions condition.Conditions
		if u.Ctx.GetContentLength() > 0 {
			if err := u.Ctx.ReadJSON(&conditions); err != nil {
				return nil, err
			}
		}
		num, _ := u.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := u.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return u.SystemLogService.Page(num, size, conditions)
	}
	return nil, nil
}
