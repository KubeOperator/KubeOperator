package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/log"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ForgotPasswordController struct {
	Ctx         context.Context
	UserService service.UserService
}

func NewForgotPasswordController() *ForgotPasswordController {
	return &ForgotPasswordController{
		UserService: service.NewUserService(),
	}
}

func (u ForgotPasswordController) PostForgotPassword() error {
	var req dto.UserForgotPassword
	err := u.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}

	go log.Save("N/A", constant.FORGOT_USER_PASSWORD, req.Username)

	return u.UserService.ResetPassword(req)
}
