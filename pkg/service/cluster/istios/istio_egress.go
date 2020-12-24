package istios

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	EgressImage = "istio/proxyv2:1.8.0"
)

type EgressInterface struct {
	Component *model.ClusterIstio
	HelmInfo  IstioHelmInfo
}

func NewEgressInterface(component *model.ClusterIstio, helmInfo IstioHelmInfo) *EgressInterface {
	return &EgressInterface{
		Component: component,
		HelmInfo:  helmInfo,
	}
}

func (e *EgressInterface) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(e.Component.Vars), &values)
	values["global.proxy.image"] = fmt.Sprintf("%s:%d/%s", e.HelmInfo.LocalhostName, constant.LocalDockerRepositoryPort, EgressImage)
	values["global.jwtPolicy"] = "first-party-jwt"
	values["gateways.istio-egressgateway.resources.requests.cpu"] = fmt.Sprintf("%vm", values["gateways.istio-egressgateway.resources.requests.cpu"])
	values["gateways.istio-egressgateway.resources.requests.memory"] = fmt.Sprintf("%vMi", values["gateways.istio-egressgateway.resources.requests.memory"])
	values["gateways.istio-egressgateway.resources.limits.cpu"] = fmt.Sprintf("%vm", values["gateways.istio-egressgateway.resources.limits.cpu"])
	values["gateways.istio-egressgateway.resources.limits.memory"] = fmt.Sprintf("%vMi", values["gateways.istio-egressgateway.resources.limits.memory"])

	str, _ := json.Marshal(&values)
	e.Component.Vars = string(str)
}

func (e *EgressInterface) Install() error {
	e.setDefaultValue()
	if err := installChart(e.HelmInfo.HelmClient, e.Component, constant.EgressChartName); err != nil {
		return err
	}
	return nil
}

func (e *EgressInterface) Uninstall() error {
	return uninstall(e.Component, e.HelmInfo.HelmClient)
}
