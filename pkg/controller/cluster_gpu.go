package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type ClusterGpuController struct {
	Ctx               context.Context
	ClusterGpuService service.ClusterGpuService
}

func NewClusterGpuController() *ClusterGpuController {
	return &ClusterGpuController{
		ClusterGpuService: service.NewClusterGpuService(),
	}
}

func (c ClusterGpuController) PostBy(clusterName string, operation string) (*dto.ClusterGpu, error) {
	cts, err := c.ClusterGpuService.HandleGPU(clusterName, operation)
	if err != nil {
		return nil, err
	}

	return cts, nil
}

func (c ClusterGpuController) GetBy(clusterName string) (*dto.ClusterGpu, error) {
	cts, err := c.ClusterGpuService.Get(clusterName)
	if err != nil {
		return nil, err
	}

	return cts, nil
}
