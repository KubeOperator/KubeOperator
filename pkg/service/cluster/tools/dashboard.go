package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Dashboard struct {
	Cluster             *Cluster
	Tool                *model.ClusterTool
	LocalHostName       string
	LocalRepositoryPort int
}

func NewDashboard(cluster *Cluster, tool *model.ClusterTool) (*Dashboard, error) {
	p := &Dashboard{
		Tool:                tool,
		Cluster:             cluster,
		LocalHostName:       constant.LocalRepositoryDomainName,
		LocalRepositoryPort: cluster.helmRepoPort,
	}
	return p, nil
}

func (d Dashboard) setDefaultValue(toolDetail model.ClusterToolDetail) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)

	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(d.Tool.Vars), &values)
	values["extraArgs[0]"] = "--enable-skip-login"
	values["extraArgs[1]"] = "--enable-insecure-login"
	values["protocolHttp"] = "true"
	values["service.externalPort"] = 9090
	values["metricsScraper.enabled"] = true
	values["metricsScraper.image.repository"] = fmt.Sprintf("%s:%d/%s", d.LocalHostName, d.LocalRepositoryPort, imageMap["metrics_image_name"])
	values["metricsScraper.image.tag"] = imageMap["metrics_image_tag"]
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", d.LocalHostName, d.LocalRepositoryPort, imageMap["dashboard_image_name"])
	values["image.tag"] = imageMap["dashboard_image_tag"]
	str, _ := json.Marshal(&values)
	d.Tool.Vars = string(str)
}

func (d Dashboard) Install(toolDetail model.ClusterToolDetail) error {
	d.setDefaultValue(toolDetail)
	if err := installChart(d.Cluster.HelmClient, d.Tool, constant.DashboardChartName, toolDetail.ChartVersion); err != nil {
		return err
	}

	ingressItem := &Ingress{
		name:    constant.DefaultDashboardIngressName,
		url:     constant.DefaultDashboardIngress,
		service: constant.DefaultDashboardServiceName,
		port:    9090,
		version: d.Cluster.Version,
	}
	if err := createRoute(d.Cluster.Namespace, ingressItem, d.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(d.Cluster.Namespace, constant.DefaultDashboardDeploymentName, 1, d.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (d Dashboard) Upgrade(toolDetail model.ClusterToolDetail) error {
	d.setDefaultValue(toolDetail)
	return upgradeChart(d.Cluster.HelmClient, d.Tool, constant.DashboardChartName, toolDetail.ChartVersion)
}

func (d Dashboard) Uninstall() error {
	return uninstall(d.Cluster.Namespace, d.Tool, constant.DefaultDashboardIngressName, d.Cluster.Version, d.Cluster.HelmClient, d.Cluster.KubeClient)
}
