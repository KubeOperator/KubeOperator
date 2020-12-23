package istios

import (
	"encoding/json"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type RemoteInterface struct {
	Component *model.ClusterIstio
	HelmInfo  IstioHelmInfo
}

func NewRemoteInterface(component *model.ClusterIstio, helmInfo IstioHelmInfo) *RemoteInterface {
	return &RemoteInterface{
		Component: component,
		HelmInfo:  helmInfo,
	}
}

func (r *RemoteInterface) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(r.Component.Vars), &values)
	values["global.proxy.autoInject"] = "enabled"

	str, _ := json.Marshal(&values)
	r.Component.Vars = string(str)
}

func (r *RemoteInterface) Install() error {
	r.setDefaultValue()
	if err := installChart(r.HelmInfo.HelmClient, r.Component, constant.RemoteChartName); err != nil {
		return err
	}
	return nil
}

func (r *RemoteInterface) Uninstall() error {
	return uninstall(r.Component, r.HelmInfo.HelmClient)
}
