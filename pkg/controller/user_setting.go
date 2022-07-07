package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type UserSettingController struct {
	Ctx                context.Context
	UserSettingService service.UserSettingService
}

func NewUserSettingController() *UserSettingController {
	return &UserSettingController{
		UserSettingService: service.NewUserSettingService(),
	}
}

func (u *UserSettingController) GetBy(username string) (dto.UserSettingDTO, error) {
	return u.UserSettingService.GetByUsername(username)
}

func (u *UserSettingController) PostUpdate() (dto.UserSettingDTO, error) {
	var updated dto.UserSettingDTO
	if err := u.Ctx.ReadJSON(&updated); err != nil {
		return updated, err
	}
	return u.UserSettingService.Update(updated)
}
