package istios

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	PilotImageName = "istio/pilot:1.8.0"
)

type PilotInterface struct {
	Component *model.ClusterIstio
	HelmInfo  IstioHelmInfo
}

func NewPilotInterface(component *model.ClusterIstio, helmInfo IstioHelmInfo) *PilotInterface {
	return &PilotInterface{
		Component: component,
		HelmInfo:  helmInfo,
	}
}

func (d *PilotInterface) setDefaultValue() map[string]interface{} {
	values := map[string]interface{}{}
	if err := json.Unmarshal([]byte(d.Component.Vars), &values); err != nil {
		log.Errorf("json unmarshal falied : %v", d.Component.Vars)
	}
	values["pilot.image"] = fmt.Sprintf("%s:%d/%s", d.HelmInfo.LocalhostName, constant.LocalDockerRepositoryPort, PilotImageName)
	values["global.jwtPolicy"] = "first-party-jwt"
	values["pilot.resources.requests.cpu"] = fmt.Sprintf("%vm", values["pilot.resources.requests.cpu"])
	values["pilot.resources.requests.memory"] = fmt.Sprintf("%vMi", values["pilot.resources.requests.memory"])
	values["pilot.resources.limits.cpu"] = fmt.Sprintf("%vm", values["pilot.resources.limits.cpu"])
	values["pilot.resources.limits.memory"] = fmt.Sprintf("%vMi", values["pilot.resources.limits.memory"])

	return values
}

func (d *PilotInterface) Install() error {
	valueMaps := d.setDefaultValue()
	if err := installChart(d.HelmInfo.HelmClient, d.Component, valueMaps, constant.PilotChartName); err != nil {
		return err
	}
	return nil
}

func (d *PilotInterface) Uninstall() error {
	return uninstall(d.Component, d.HelmInfo.HelmClient)
}
