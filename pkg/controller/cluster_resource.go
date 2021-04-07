package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ClusterResourceController struct {
	Ctx                    context.Context
	ClusterResourceService service.ClusterResourceService
}

func NewClusterResourceController() *ClusterResourceController {
	return &ClusterResourceController{
		ClusterResourceService: service.NewClusterResourceService(),
	}
}

func (c ClusterResourceController) Get() (*page.Page, error) {
	pa, _ := c.Ctx.Values().GetBool("page")
	resourceType := c.Ctx.URLParam("resourceType")
	clusterName := c.Ctx.Params().GetString("cluster")
	if pa {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return c.ClusterResourceService.Page(num, size, clusterName, resourceType)
	} else {
		return nil, nil
	}
}

func (c ClusterResourceController) Post() ([]dto.ClusterResource, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	var req dto.ClusterResourceCreate
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}
	return c.ClusterResourceService.Create(clusterName, req)
}

func (c ClusterResourceController) DeleteBy(name string) error {
	resourceType := c.Ctx.URLParam("resourceType")
	clusterName := c.Ctx.Params().GetString("cluster")
	return c.ClusterResourceService.Delete(name, resourceType, clusterName)
}

func (c ClusterResourceController) GetList() (interface{}, error) {
	resourceType := c.Ctx.URLParam("resourceType")
	projectName := c.Ctx.Values().GetString("project")
	clusterName := c.Ctx.Params().GetString("cluster")
	return c.ClusterResourceService.GetResources(resourceType, projectName, clusterName)
}
