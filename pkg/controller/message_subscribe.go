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
	pa, _ := m.Ctx.Values().GetBool("page")
	var p page.Page
	var conditions condition.Conditions
	if m.Ctx.GetContentLength() > 0 {
		if err := m.Ctx.ReadJSON(&conditions); err != nil {
			return p, err
		}
	}
	resourceName := m.Ctx.URLParam("resourceName")
	scope := m.Ctx.URLParam("type")

	if pa {
		num, _ := m.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := m.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return m.MsgSubscribeService.Page(scope, resourceName, num, size, conditions)
	} else {
		items, err := m.MsgSubscribeService.List(scope, resourceName, conditions)
		if err != nil {
			return p, nil
		}
		p.Items = items
		p.Total = len(items)
		return p, nil
	}
}

func (m MessageSubscribeController) PostUpdate() error {
	var updated dto.MsgSubscribeDTO
	if err := m.Ctx.ReadJSON(&updated); err != nil {
		return err
	}
	return m.MsgSubscribeService.Update(updated)
}

func (m MessageSubscribeController) PostUser() error {
	var add dto.MsgSubscribeUserDTO
	if err := m.Ctx.ReadJSON(&add); err != nil {
		return err
	}
	return m.MsgSubscribeService.AddSubscribeUser(add)
}

func (m MessageSubscribeController) PostDeleteUser() error {
	var del dto.MsgSubscribeUserDTO
	if err := m.Ctx.ReadJSON(&del); err != nil {
		return err
	}
	return m.MsgSubscribeService.DeleteSubscribeUser(del)
}

func (m MessageSubscribeController) GetUsers() (dto.AddSubscribeResponse, error) {
	sessionUser := m.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)
	resourceName := m.Ctx.URLParam("resourceName")
	search := m.Ctx.URLParam("name")

	return m.MsgSubscribeService.GetSubscribeUser(resourceName, search, user)
}
