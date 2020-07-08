package tools

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Dashboard struct {
	Cluster *Cluster
	Tool    *model.ClusterTool
}

func NewDashboard(cluster *Cluster, tool *model.ClusterTool) (*Dashboard, error) {
	p := &Dashboard{
		Tool:    tool,
		Cluster: cluster,
	}
	return p, nil
}

func (c Dashboard) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &values)
	values["extraArgs[0]"] = "--enable-skip-login"
	values["extraArgs[1]"] = "--enable-insecure-login"
	values["protocolHttp"] = "true"
	values["service.externalPort"] = 9090
	values["metricsScraper.enabled"] = true
	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
}

func (c Dashboard) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.DashboardChartName); err != nil {
		return err
	}
	if err := createRoute(constant.DefaultDashboardIngressName, constant.DefaultDashboardIngress, constant.DefaultDashboardServiceName, 9090, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(constant.DefaultDashboardDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c Dashboard) Uninstall() error {
	return uninstall(c.Tool, constant.DefaultDashboardIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
