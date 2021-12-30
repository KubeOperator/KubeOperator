package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/util/prometheus"
	"github.com/kataras/iris/v12/context"
)

type MonitorController struct {
	Ctx                context.Context
	ClusterToolService service.ClusterToolService
}

func NewMonitorController() *MonitorController {
	return &MonitorController{
		ClusterToolService: service.NewClusterToolService(),
	}
}

func (m MonitorController) PostSearchBy(clusterName string) ([]dto.Metric, error) {
	var req dto.QueryOptions
	if err := m.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	endPoint, err := m.ClusterToolService.GetNodePort(clusterName, "prometheus")
	if err != nil {
		return nil, err
	}
	prometheusClient, err := prometheus.NewPrometheusService(endPoint)
	if err != nil {
		return nil, err
	}

	return prometheusClient.GetNamedMetricsOverTime(&req), nil
}
