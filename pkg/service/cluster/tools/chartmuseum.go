package tools

import (
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
)

const (
	ChartmuseumImageName = "chartmuseum/chartmuseum"
	ChartmuseumTag       = "v0.12.0"
)

type Chartmuseum struct {
	Cluster *Cluster
	Tool    *model.ClusterTool
}

func NewChartmuseum(cluster *Cluster, tool *model.ClusterTool) (*Chartmuseum, error) {
	p := &Chartmuseum{
		Tool:    tool,
		Cluster: cluster,
	}
	return p, nil
}

func (c Chartmuseum) setDefaultValue() {
	systemService := service.NewSystemSettingService()
	locahostName := systemService.GetLocalHostName()
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &values)
	values["env.open.DISABLE_API"] = false
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", locahostName, constant.LocalDockerRepositoryPort, ChartmuseumImageName)
	values["image.tag"] = ChartmuseumTag
	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
}

func (c Chartmuseum) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.ChartmuseumChartName); err != nil {
		return err
	}
	if err := createRoute(constant.DefaultChartmuseumIngressName, constant.DefaultChartmuseumIngress, constant.DefaultChartmuseumServiceName, 8080, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(constant.DefaultChartmuseumDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c Chartmuseum) Uninstall() error {
	return uninstall(c.Tool, constant.DefaultChartmuseumIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
