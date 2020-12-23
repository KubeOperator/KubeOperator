package istios

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	OperatorImageHub = "istio"
	OperatorTag      = "1.8.0"
)

type OperatorInterface struct {
	Component *model.ClusterIstio
	HelmInfo  IstioHelmInfo
}

func NewOperatorInterface(component *model.ClusterIstio, helmInfo IstioHelmInfo) *OperatorInterface {
	return &OperatorInterface{
		Component: component,
		HelmInfo:  helmInfo,
	}
}

func (o OperatorInterface) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(o.Component.Vars), &values)
	values["hub"] = fmt.Sprintf("%s:%o/%s", o.HelmInfo.LocalhostName, constant.LocalDockerRepositoryPort, OperatorImageHub)
	values["tag"] = OperatorTag
	values["operator.resources.requests.cpu"] = fmt.Sprintf("%vm", values["operator.resources.requests.cpu"])
	values["operator.resources.requests.memory"] = fmt.Sprintf("%vMi", values["operator.resources.requests.memory"])
	values["operator.resources.limits.cpu"] = fmt.Sprintf("%vm", values["operator.resources.limits.cpu"])
	values["operator.resources.limits.memory"] = fmt.Sprintf("%vMi", values["operator.resources.limits.memory"])

	str, _ := json.Marshal(&values)
	o.Component.Vars = string(str)
}

func (o OperatorInterface) Install() error {
	o.setDefaultValue()
	if err := installChart(o.HelmInfo.HelmClient, o.Component, constant.OperatorChartName); err != nil {
		return err
	}
	return nil
}

func (o OperatorInterface) Uninstall() error {
	return uninstall(o.Component, o.HelmInfo.HelmClient)
}
