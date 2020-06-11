package tools

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	clusterService "github.com/KubeOperator/KubeOperator/pkg/service/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
)

type EFKManager struct {
	clusterName string
}

func (e *EFKManager) Install(values map[string]interface{}) error {
	helmClient, err := clusterService.GetHelmClient(e.clusterName)
	if err != nil {
		return err
	}
	chart, err := helm.LoadCharts(constant.EFKChartPath)
	if err != nil {
		return err
	}
	_, err = helmClient.Install(constant.EFKReleaseName, chart, values)
	if err != nil {
		return err
	}
	return nil
}
