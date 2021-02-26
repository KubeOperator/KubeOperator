package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Grafana struct {
	Cluster             *Cluster
	Tool                *model.ClusterTool
	LocalHostName       string
	LocalRepositoryPort int
	prometheusNs        string
	lokiNs              string
}

func NewGrafana(cluster *Cluster, tool *model.ClusterTool, prometheusNs, lokiNs string) (*Grafana, error) {
	p := &Grafana{
		Tool:                tool,
		Cluster:             cluster,
		LocalHostName:       constant.LocalRepositoryDomainName,
		LocalRepositoryPort: constant.LocalDockerRepositoryPort,
		prometheusNs:        prometheusNs,
		lokiNs:              lokiNs,
	}
	return p, nil
}

func (g Grafana) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)

	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(g.Tool.Vars), &values)
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", g.LocalHostName, g.LocalRepositoryPort, imageMap["grafana_image_name"])
	values["image.tag"] = imageMap["grafana_image_tag"]
	values["initChownData.enabled"] = true
	values["initChownData.image.repository"] = fmt.Sprintf("%s:%d/%s", g.LocalHostName, g.LocalRepositoryPort, imageMap["busybox_image_name"])
	values["initChownData.image.tag"] = imageMap["busybox_image_tag"]
	values["downloadDashboardsImage.repository"] = imageMap["curl_image_name"]
	values["downloadDashboardsImage.tag"] = imageMap["curl_image_tag"]

	if isInstall {
		values["grafana\\.ini.server.root_url"] = "%(protocol)s://%(domain)s:%(http_port)s/proxy/grafana/" + g.Cluster.Name + "/"
		values["grafana\\.ini.server.serve_from_sub_path"] = true

		values["datasources.'datasources\\.yaml'.apiVersion"] = 1

		if len(g.prometheusNs) != 0 {
			values["datasources.'datasources\\.yaml'.datasources[0].name"] = "MYDS_Prometheus"
			values["datasources.'datasources\\.yaml'.datasources[0].type"] = "prometheus"
			values["datasources.'datasources\\.yaml'.datasources[0].url"] = "http://prometheus-server." + g.prometheusNs
			values["datasources.'datasources\\.yaml'.datasources[0].access"] = "proxy"
			values["datasources.'datasources\\.yaml'.datasources[0].isDefault"] = true
		}
		if len(g.lokiNs) != 0 {
			values["datasources.'datasources\\.yaml'.datasources[1].name"] = "Loki"
			values["datasources.'datasources\\.yaml'.datasources[1].type"] = "loki"
			values["datasources.'datasources\\.yaml'.datasources[1].url"] = "http://loki." + g.lokiNs + ":3100"
			values["datasources.'datasources\\.yaml'.datasources[1].access"] = "proxy"
		}

		values["dashboardProviders.'dashboardproviders\\.yaml'.apiVersion"] = 1
		values["dashboardProviders.'dashboardproviders\\.yaml'.providers[0].name"] = "default"
		values["dashboardProviders.'dashboardproviders\\.yaml'.providers[0].orgId"] = 1
		values["dashboardProviders.'dashboardproviders\\.yaml'.providers[0].folder"] = ""
		values["dashboardProviders.'dashboardproviders\\.yaml'.providers[0].type"] = "file"
		values["dashboardProviders.'dashboardproviders\\.yaml'.providers[0].disableDeletion"] = false
		values["dashboardProviders.'dashboardproviders\\.yaml'.providers[0].editable"] = true
		values["dashboardProviders.'dashboardproviders\\.yaml'.providers[0].options.path"] = "/var/lib/grafana/dashboards/default"
		values["dashboards.default.custom-dashboard.file"] = "dashboards/custom-dashboard.json"

		if _, ok := values["persistence.size"]; ok {
			values["persistence.size"] = fmt.Sprintf("%vGi", values["persistence.size"])
		}
		if va, ok := values["persistence.enabled"]; ok {
			if hasPers, _ := va.(bool); !hasPers {
				delete(values, "nodeSelector.kubernetes\\.io/hostname")
			}
		}
	}

	str, _ := json.Marshal(&values)
	g.Tool.Vars = string(str)
}

func (g Grafana) Install(toolDetail model.ClusterToolDetail) error {
	g.setDefaultValue(toolDetail, true)
	if err := installChart(g.Cluster.HelmClient, g.Tool, constant.GrafanaChartName, toolDetail.ChartVersion); err != nil {
		return err
	}
	if err := createRoute(g.Cluster.Namespace, constant.DefaultGrafanaIngressName, constant.DefaultGrafanaIngress, constant.DefaultGrafanaServiceName, 80, g.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(g.Cluster.Namespace, constant.DefaultGrafanaDeploymentName, 1, g.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (g Grafana) Upgrade(toolDetail model.ClusterToolDetail) error {
	g.setDefaultValue(toolDetail, false)
	return upgradeChart(g.Cluster.HelmClient, g.Tool, constant.GrafanaChartName, toolDetail.ChartVersion)
}

func (g Grafana) Uninstall() error {
	return uninstall(g.Cluster.Namespace, g.Tool, constant.DefaultGrafanaIngressName, g.Cluster.HelmClient, g.Cluster.KubeClient)
}
