package istios

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	IngressImage = "istio/proxyv2:1.8.0"
)

type IngressInterface struct {
	Component *model.ClusterIstio
	HelmInfo  IstioHelmInfo
}

func NewIngressInterface(component *model.ClusterIstio, helmInfo IstioHelmInfo) *IngressInterface {
	return &IngressInterface{
		Component: component,
		HelmInfo:  helmInfo,
	}
}

func (i *IngressInterface) setDefaultValue() map[string]interface{} {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(i.Component.Vars), &values)
	values["global.proxy.image"] = fmt.Sprintf("%s:%d/%s", i.HelmInfo.LocalhostName, i.HelmInfo.LocalhostPort, IngressImage)
	values["global.jwtPolicy"] = "first-party-jwt"
	values["gateways.istio-ingressgateway.resources.requests.cpu"] = fmt.Sprintf("%vm", values["gateways.istio-ingressgateway.resources.requests.cpu"])
	values["gateways.istio-ingressgateway.resources.requests.memory"] = fmt.Sprintf("%vMi", values["gateways.istio-ingressgateway.resources.requests.memory"])
	values["gateways.istio-ingressgateway.resources.limits.cpu"] = fmt.Sprintf("%vm", values["gateways.istio-ingressgateway.resources.limits.cpu"])
	values["gateways.istio-ingressgateway.resources.limits.memory"] = fmt.Sprintf("%vMi", values["gateways.istio-ingressgateway.resources.limits.memory"])
	return values
}

func (i *IngressInterface) Install() error {
	valueMaps := i.setDefaultValue()
	if err := installChart(i.HelmInfo.HelmClient, i.Component, valueMaps, constant.IngressChartName); err != nil {
		return err
	}
	return nil
}

func (i *IngressInterface) Uninstall() error {
	return uninstall(i.Component, i.HelmInfo.HelmClient)
}
