package istios

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type BaseInterface struct {
	Component *model.ClusterIstio
	HelmInfo  IstioHelmInfo
}

func NewBaseInterface(component *model.ClusterIstio, helmInfo IstioHelmInfo) *BaseInterface {
	return &BaseInterface{
		Component: component,
		HelmInfo:  helmInfo,
	}
}

func (b BaseInterface) Install() error {
	if err := installChart(b.HelmInfo.HelmClient, b.Component, constant.BaseChartName); err != nil {
		return err
	}
	return nil
}

func (b BaseInterface) Uninstall() error {
	return uninstall(b.Component, b.HelmInfo.HelmClient)
}
