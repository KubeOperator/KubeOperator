package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type UserMsgController struct {
	Ctx            context.Context
	UserMsgService service.UserMsgService
}

func NewUserMsgController() *UserMsgController {
	return &UserMsgController{
		UserMsgService: service.NewUserMsgService(),
	}
}

func (u *UserMsgController) Get() (dto.UserMsgResponse, error) {
	p, _ := u.Ctx.Values().GetBool("page")
	sessionUser := u.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)
	if p {
		num, _ := u.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := u.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return u.UserMsgService.PageLocalMsg(num, size, user, condition.TODO())
	}
	return dto.UserMsgResponse{}, nil
}

func (u *UserMsgController) PostReadBy(msgID string) error {
	sessionUser := u.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)

	return u.UserMsgService.UpdateLocalMsg(msgID, user)
}
