package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type ClusterIstioController struct {
	Ctx                 context.Context
	ClusterIstioService service.ClusterIstioService
}

func NewClusterIstioController() *ClusterIstioController {
	return &ClusterIstioController{
		ClusterIstioService: service.NewClusterIstioService(),
	}
}

func (c ClusterIstioController) GetBy(clusterName string) ([]dto.ClusterIstio, error) {
	cts, err := c.ClusterIstioService.List(clusterName)
	if err != nil {
		return nil, err
	}
	return cts, nil
}

func (c ClusterIstioController) PostEnableBy(clusterName string) (*[]dto.ClusterIstio, error) {
	var req []dto.ClusterIstio
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	cts, err := c.ClusterIstioService.Enable(clusterName, req)
	if err != nil {
		return nil, err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.ENABLE_CLUSTER_ISTIO, clusterName)

	return &cts, nil
}

func (c ClusterIstioController) PostDisableBy(clusterName string) (*[]dto.ClusterIstio, error) {
	var req []dto.ClusterIstio
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	cts, err := c.ClusterIstioService.Disable(clusterName, req)
	if err != nil {
		return nil, err
	}

	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DISABLE_CLUSTER_ISTIO, clusterName)

	return &cts, nil
}
