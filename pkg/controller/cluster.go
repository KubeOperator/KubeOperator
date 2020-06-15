package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/kataras/iris"
)

type clusterController struct {
	Ctx            iris.Context
	clusterService service.ClusterService
}

func (c clusterController) Get() ([]dto.Cluster, error) {
	return c.clusterService.List()
}

func (c clusterController) GetBy(name string) (dto.Cluster, error) {
	return c.clusterService.Get(name)
}

func (c clusterController) GetStatus(name string) (dto.ClusterStatus, error) {
	return c.clusterService.GetStatus(name)
}

func (c clusterController) Post() error {
	var req dto.ClusterCreate
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	return c.clusterService.Create(req)
}

func (c clusterController) Delete(name string) error {
	return c.clusterService.Delete(name)
}
