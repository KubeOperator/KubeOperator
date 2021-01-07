package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	GrafanaImageName    = "kubeoperator/grafana"
	GrafanaTagAmd64Name = "7.3.3-amd64"
	GrafanaTagArm64Name = "7.3.3-arm64"

	initChownDataImageName    = "kubeoperator/busybox"
	initChownDataTagAmd64Name = "1.28-amd64"
	initChownDataTagArm64Name = "1.28-arm64"
)

type Grafana struct {
	Cluster       *Cluster
	Tool          *model.ClusterTool
	LocalHostName string
	prometheusNs  string
	lokiNs        string
}

func (c Grafana) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &values)
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, constant.LocalDockerRepositoryPort, GrafanaImageName)

	if c.Cluster.Spec.Architectures == "amd64" {
		values["image.tag"] = GrafanaTagAmd64Name
		values["initChownData.image.tag"] = initChownDataTagAmd64Name
	} else {
		values["image.tag"] = GrafanaTagArm64Name
		values["initChownData.image.tag"] = initChownDataTagArm64Name
	}

	values["grafana\\.ini.server.root_url"] = "%(protocol)s://%(domain)s:%(http_port)s/proxy/grafana/" + c.Cluster.Name + "/"
	values["grafana\\.ini.server.serve_from_sub_path"] = true
	values["initChownData.enabled"] = true
	values["initChownData.image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, constant.LocalDockerRepositoryPort, initChownDataImageName)
	values["datasources.'datasources\\.yaml'.apiVersion"] = 1

	if len(c.prometheusNs) != 0 {
		values["datasources.'datasources\\.yaml'.datasources[0].name"] = "MYDS_Prometheus"
		values["datasources.'datasources\\.yaml'.datasources[0].type"] = "prometheus"
		values["datasources.'datasources\\.yaml'.datasources[0].url"] = "http://prometheus-server." + c.prometheusNs
		values["datasources.'datasources\\.yaml'.datasources[0].access"] = "proxy"
		values["datasources.'datasources\\.yaml'.datasources[0].isDefault"] = true
	}
	if len(c.lokiNs) != 0 {
		values["datasources.'datasources\\.yaml'.datasources[1].name"] = "Loki"
		values["datasources.'datasources\\.yaml'.datasources[1].type"] = "loki"
		values["datasources.'datasources\\.yaml'.datasources[1].url"] = "http://loki." + c.lokiNs + ":3100"
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

	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
}

func NewGrafana(cluster *Cluster, localhostName string, tool *model.ClusterTool, prometheusNs, lokiNs string) (*Grafana, error) {
	p := &Grafana{
		Tool:          tool,
		Cluster:       cluster,
		LocalHostName: localhostName,
		prometheusNs:  prometheusNs,
		lokiNs:        lokiNs,
	}
	return p, nil
}

func (c Grafana) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.GrafanaChartName); err != nil {
		return err
	}
	if err := createRoute(c.Cluster.Namespace, constant.DefaultGrafanaIngressName, constant.DefaultGrafanaIngress, constant.DefaultGrafanaServiceName, 80, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(c.Cluster.Namespace, constant.DefaultGrafanaDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c Grafana) Uninstall() error {
	return uninstall(c.Cluster.Namespace, c.Tool, constant.DefaultGrafanaIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
