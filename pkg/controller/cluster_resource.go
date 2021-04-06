package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/service"
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
