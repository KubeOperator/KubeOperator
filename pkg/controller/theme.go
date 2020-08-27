package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type ThemeController struct {
	Ctx          context.Context
	ThemeService service.ThemeService
}

func NewThemeController() *ThemeController {
	return &ThemeController{
		ThemeService: service.NewThemeService(),
	}
}

func (l *ThemeController) Get() (*dto.Theme, error) {
	return l.ThemeService.GetConsumerTheme()
}
