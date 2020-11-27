package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	MetricsScraperImageName = "kubernetesui/metrics-scraper"
	MetricsScraperImageTag  = "v1.0.4"
	DashboardImageName      = "kubernetesui/dashboard"
	DashboardImageTag       = "v2.0.3"
)

type Dashboard struct {
	Cluster       *Cluster
	Tool          *model.ClusterTool
	LocalHostName string
}

func NewDashboard(cluster *Cluster, localhostName string, tool *model.ClusterTool) (*Dashboard, error) {
	p := &Dashboard{
		Tool:          tool,
		Cluster:       cluster,
		LocalHostName: localhostName,
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
	values["metricsScraper.image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, constant.LocalDockerRepositoryPort, MetricsScraperImageName)
	values["metricsScraper.image.tag"] = MetricsScraperImageTag
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, constant.LocalDockerRepositoryPort, DashboardImageName)
	values["image.tag"] = DashboardImageTag
	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
}

func (c Dashboard) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.DashboardChartName); err != nil {
		return err
	}
	if err := createRoute(c.Cluster.Namespace, constant.DefaultDashboardIngressName, constant.DefaultDashboardIngress, constant.DefaultDashboardServiceName, 9090, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(c.Cluster.Namespace, constant.DefaultDashboardDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c Dashboard) Uninstall() error {
	return uninstall(c.Cluster.Namespace, c.Tool, constant.DefaultDashboardIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
