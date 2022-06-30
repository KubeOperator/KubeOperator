package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type MessageAccountController struct {
	Ctx        context.Context
	MsgService service.MsgAccountService
}

func NewMessageAccountController() *MessageAccountController {
	return &MessageAccountController{
		MsgService: service.NewMsgAccountService(),
	}
}

func (m MessageAccountController) GetBy(name string) (dto.MsgAccountDTO, error) {
	return m.MsgService.GetByName(name)
}

func (m MessageAccountController) Post() (dto.MsgAccountDTO, error) {
	var req dto.MsgAccountDTO
	err := m.Ctx.ReadJSON(&req)
	if err != nil {
		return req, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return req, err
	}
	return m.MsgService.CreateOrUpdate(req)
}

func (m MessageAccountController) PostVerify() error {
	var req dto.MsgAccountDTO
	err := m.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	return m.MsgService.Verify(req)
}
