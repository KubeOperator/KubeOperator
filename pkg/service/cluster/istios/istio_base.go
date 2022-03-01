package istios

import (
	"encoding/json"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

var log = logger.Default

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

func (b *BaseInterface) setDefaultValue() map[string]interface{} {
	values := map[string]interface{}{}
	if err := json.Unmarshal([]byte(b.Component.Vars), &values); err != nil {
		log.Errorf("json unmarshal falied : %v", b.Component.Vars)
	}

	return values
}

func (b *BaseInterface) Install() error {
	valueMaps := b.setDefaultValue()
	if err := installChart(b.HelmInfo.HelmClient, b.Component, valueMaps, constant.BaseChartName); err != nil {
		return err
	}
	return nil
}

func (b *BaseInterface) Uninstall() error {
	return uninstall(b.Component, b.HelmInfo.HelmClient)
}
