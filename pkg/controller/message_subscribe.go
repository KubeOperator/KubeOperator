package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type MessageSubscribeController struct {
	Ctx                 context.Context
	MsgSubscribeService service.MsgSubscribeService
}

func NewMessageSubscribeController() *MessageSubscribeController {
	return &MessageSubscribeController{
		MsgSubscribeService: service.NewMsgSubscribeService(),
	}
}

func (m MessageSubscribeController) PostSearch() (page.Page, error) {
	//pa, _ := m.Ctx.Values().GetBool("page")
	var p page.Page
	var conditions condition.Conditions
	if m.Ctx.GetContentLength() > 0 {
		if err := m.Ctx.ReadJSON(&conditions); err != nil {
			return p, err
		}
	}
	num, _ := m.Ctx.Values().GetInt(constant.PageNumQueryKey)
	size, _ := m.Ctx.Values().GetInt(constant.PageSizeQueryKey)
	resourceName := m.Ctx.URLParam("resourceName")
	scope := m.Ctx.URLParam("type")
	return m.MsgSubscribeService.Page(scope, resourceName, num, size, conditions)
}

func (m MessageSubscribeController) PostUpdate() error {
	var updated dto.MsgSubscribeDTO
	if err := m.Ctx.ReadJSON(&updated); err != nil {
		return err
	}
	return m.MsgSubscribeService.Update(updated)
}
