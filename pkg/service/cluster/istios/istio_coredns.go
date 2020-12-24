package istios

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type CorednsInterface struct {
	Component *model.ClusterIstio
	HelmInfo  IstioHelmInfo
}

func NewCorednsInterface(component *model.ClusterIstio, helmInfo IstioHelmInfo) *CorednsInterface {
	return &CorednsInterface{
		Component: component,
		HelmInfo:  helmInfo,
	}
}

func (c *CorednsInterface) Install() error {
	if err := installChart(c.HelmInfo.HelmClient, c.Component, constant.CorednsChartName); err != nil {
		return err
	}
	return nil
}

func (c *CorednsInterface) Uninstall() error {
	return uninstall(c.Component, c.HelmInfo.HelmClient)
}
