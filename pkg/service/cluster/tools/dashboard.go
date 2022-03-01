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
		LocalRepositoryPort: constant.LocalDockerRepositoryPort,
	}
	return p, nil
}

func (d Dashboard) setDefaultValue(toolDetail model.ClusterToolDetail) {
	imageMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(toolDetail.Vars), &imageMap); err != nil {
		log.Errorf("json unmarshal falied : %v", toolDetail.Vars)
	}

	values := map[string]interface{}{}
	if err := json.Unmarshal([]byte(d.Tool.Vars), &values); err != nil {
		log.Errorf("json unmarshal falied : %v", d.Tool.Vars)
	}
	values["extraArgs[0]"] = "--enable-skip-login"
	values["extraArgs[1]"] = "--enable-insecure-login"
	values["protocolHttp"] = "true"
	values["service.externalPort"] = 9090
	values["metricsScraper.enabled"] = true
	values["metricsScraper.image.repository"] = fmt.Sprintf("%s:%d/%s", d.LocalHostName, d.LocalRepositoryPort, imageMap["metrics_image_name"])
	values["metricsScraper.image.tag"] = imageMap["metrics_image_tag"]
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", d.LocalHostName, d.LocalRepositoryPort, imageMap["dashboard_image_name"])
	values["image.tag"] = imageMap["dashboard_image_tag"]
	str, err := json.Marshal(&values)
	if err != nil {
		log.Errorf("json marshal falied : %v", values)
	}
	d.Tool.Vars = string(str)
}

func (d Dashboard) Install(toolDetail model.ClusterToolDetail) error {
	d.setDefaultValue(toolDetail)
	if err := installChart(d.Cluster.HelmClient, d.Tool, constant.DashboardChartName, toolDetail.ChartVersion); err != nil {
		return err
	}
	if err := createRoute(d.Cluster.Namespace, constant.DefaultDashboardIngressName, constant.DefaultDashboardIngress, constant.DefaultDashboardServiceName, 9090, d.Cluster.KubeClient); err != nil {
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
	return uninstall(d.Cluster.Namespace, d.Tool, constant.DefaultDashboardIngressName, d.Cluster.HelmClient, d.Cluster.KubeClient)
}
