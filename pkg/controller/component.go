package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ComponentController struct {
	Ctx              context.Context
	ComponentService service.ComponentService
}

func NewComponentController() *ComponentController {
	return &ComponentController{
		ComponentService: service.NewComponentService(),
	}
}

// List Component
// @Tags clusters
// @Summary Show all components
// @Description Show components
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /component/ [get]
func (c ComponentController) Get() ([]dto.Component, error) {
	clusterName := c.Ctx.URLParam("cluster")
	return c.ComponentService.Get(clusterName)
}

// Create Component
// @Tags components
// @Summary Create a component
// @Description 添加一个集群组件
// @Accept  json
// @Produce  json
// @Param request body dto.ComponentCreate true "request"
// @Success 200 {object} dto.Credential
// @Security ApiKeyAuth
// @Router /credentials/ [post]
func (c ComponentController) Post() error {
	var req dto.ComponentCreate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_COMPONENT, req.Name+req.Version)

	return c.ComponentService.Create(&req)
}
