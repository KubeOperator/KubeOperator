package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type TaskLogController struct {
	Ctx            context.Context
	TaskLogService service.TaskLogService
}

func NewTaskLogController() *TaskLogController {
	return &TaskLogController{
		TaskLogService: service.NewTaskLogService(),
	}
}

// Search TaskLog
// @Tags task_logs
// @Summary Search tasklog
// @Description 过滤任务日志
// @Accept  json
// @Produce  json
// @Param conditions body condition.Conditions true "conditions"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /tasklog/ [post]
func (u TaskLogController) Post() (*page.Page, error) {
	p, _ := u.Ctx.Values().GetBool("page")
	if p {
		num, _ := u.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := u.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		cluster := u.Ctx.URLParam("cluster")
		logtype := u.Ctx.URLParam("logtype")
		p, err := u.TaskLogService.Page(num, size, cluster, logtype)
		return p, err
	}
	return nil, nil
}
