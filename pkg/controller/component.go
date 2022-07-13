package controller

import (
	"fmt"
	"strings"

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
// @Param cluster path string true "集群名称"
// @Success 200 {object} []dto.Component
// @Security ApiKeyAuth
// @Router /components/ [get]
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
// @Success 200
// @Security ApiKeyAuth
// @Router /components/ [post]
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

// Delete Component
// @Tags components
// @Summary Delete a component
// @Description 删除一个集群组件
// @Accept  json
// @Produce  json
// @Param cluster path string true "集群名称"
// @Param name path string true "组件名称"
// @Success 200
// @Security ApiKeyAuth
// @Router /components/{cluster}/{name} [delete]
func (c ComponentController) DeleteBy(cluster, name string) error {
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_COMPONENT, fmt.Sprintf("%s (%s)", name, cluster))

	return c.ComponentService.Delete(cluster, name)
}

// Sync Component
// @Tags components
// @Summary Sync components
// @Description 同步集群组件
// @Accept  json
// @Produce  json
// @Param request body dto.ComponentSync true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /components/sync [post]
func (c ComponentController) PostSync() error {
	var req dto.ComponentSync
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
	go kolog.Save(operator, constant.SYNC_COMPONENT, fmt.Sprintf("%s (%s)", req.ClusterName, strings.Join(req.Names, ",")))

	return c.ComponentService.Sync(&req)
}
