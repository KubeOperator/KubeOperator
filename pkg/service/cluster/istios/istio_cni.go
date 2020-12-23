package istios

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	CniImageName = "istio/install-cni:1.8.0"
)

type CniInterface struct {
	Component *model.ClusterIstio
	HelmInfo  IstioHelmInfo
}

func NewCniInterface(component *model.ClusterIstio, helmInfo IstioHelmInfo) *CniInterface {
	return &CniInterface{
		Component: component,
		HelmInfo:  helmInfo,
	}
}

func (c *CniInterface) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Component.Vars), &values)
	values["cni.image"] = fmt.Sprintf("%s:%d/%s", c.HelmInfo.LocalhostName, constant.LocalDockerRepositoryPort, CniImageName)

	str, _ := json.Marshal(&values)
	c.Component.Vars = string(str)
}

func (c *CniInterface) Install() error {
	c.setDefaultValue()
	if err := installChart(c.HelmInfo.HelmClient, c.Component, constant.CniChartName); err != nil {
		return err
	}
	return nil
}

func (c *CniInterface) Uninstall() error {
	return uninstall(c.Component, c.HelmInfo.HelmClient)
}
