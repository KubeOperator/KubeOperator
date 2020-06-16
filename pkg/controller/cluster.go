package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/kataras/iris/v12/context"
)

type ClusterController struct {
	Ctx            context.Context
	ClusterService service.ClusterService
}

func NewClusterController() *ClusterController {
	return &ClusterController{
		ClusterService: service.NewClusterService(),
	}
}

func (c ClusterController) Get() (dto.ClusterPage, error) {
	page, _ := c.Ctx.Values().GetBool("page")
	if page {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return c.ClusterService.Page(num, size)
	} else {
		var page dto.ClusterPage
		items, err := c.ClusterService.List()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}
}

func (c ClusterController) GetBy(name string) (dto.Cluster, error) {
	return c.ClusterService.Get(name)
}

func (c ClusterController) GetStatus(name string) (dto.ClusterStatus, error) {
	return c.ClusterService.GetStatus(name)
}

func (c ClusterController) Post() error {
	var req dto.ClusterCreate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	return c.ClusterService.Create(req)
}

func (c ClusterController) Delete(name string) error {
	return c.ClusterService.Delete(name)
}
