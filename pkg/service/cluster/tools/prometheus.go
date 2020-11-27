package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/grafana"
)

const (
	PrometheusConfigMapReloadImageName = "jimmidyson/configmap-reload"
	PrometheusConfigMapReloadTag       = "v0.3.0"
	KubeStateMetricsImageArm64Name     = "carlosedp/kube-state-metrics"
	KubeStateMetricsImageAmd64Name     = "coreos/kube-state-metrics"
	KubeStateMetricsTag                = "v1.9.5"
	NodeExporterImageName              = "prom/node-exporter"
	NodeExporterTag                    = "v0.18.1"
	ServerImageName                    = "prom/prometheus"
	ServerTag                          = "v2.18.1"
)

type Prometheus struct {
	Tool          *model.ClusterTool
	Cluster       *Cluster
	LocalhostName string
}

func NewPrometheus(cluster *Cluster, localhostName string, tool *model.ClusterTool) (*Prometheus, error) {
	p := &Prometheus{
		Tool:          tool,
		Cluster:       cluster,
		LocalhostName: localhostName,
	}
	return p, nil
}

func (p Prometheus) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(p.Tool.Vars), &values)
	values["alertmanager.enabled"] = false
	values["pushgateway.enabled"] = false
	values["configmapReload.prometheus.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalhostName, constant.LocalDockerRepositoryPort, PrometheusConfigMapReloadImageName)
	values["configmapReload.prometheus.image.tag"] = PrometheusConfigMapReloadTag
	values["kube-state-metrics.image.tag"] = KubeStateMetricsTag
	values["nodeExporter.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalhostName, constant.LocalDockerRepositoryPort, NodeExporterImageName)
	values["nodeExporter.image.tag"] = NodeExporterTag
	values["server.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalhostName, constant.LocalDockerRepositoryPort, ServerImageName)
	values["server.image.tag"] = ServerTag
	switch p.Cluster.Spec.Architectures {
	case "amd64":
		values["kube-state-metrics.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalhostName, constant.LocalDockerRepositoryPort, KubeStateMetricsImageAmd64Name)
	case "arm64":
		values["kube-state-metrics.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalhostName, constant.LocalDockerRepositoryPort, KubeStateMetricsImageArm64Name)
	}
	if _, ok := values["server.retention"]; ok {
		values["server.retention"] = fmt.Sprintf("%vd", values["server.retention"])
	}
	if _, ok := values["server.persistentVolume.size"]; ok {
		values["server.persistentVolume.size"] = fmt.Sprintf("%vGi", values["server.persistentVolume.size"])
	}
	str, _ := json.Marshal(&values)
	p.Tool.Vars = string(str)
}

func (c Prometheus) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.PrometheusChartName); err != nil {
		return err
	}
	if err := createRoute(c.Cluster.Namespace, constant.DefaultPrometheusIngressName, constant.DefaultPrometheusIngress, constant.DefaultPrometheusServiceName, 80, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(c.Cluster.Namespace, constant.DefaultPrometheusDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}

	if err := c.createGrafanaDataSource(); err != nil {
		return err
	}

	if err := c.createGrafanaDashboard(); err != nil {
		return err
	}
	return nil
}

func (c Prometheus) Uninstall() error {
	return uninstall(c.Cluster.Namespace, c.Tool, constant.DefaultPrometheusIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}

func (p Prometheus) createGrafanaDataSource() error {
	grafanaClient := grafana.NewClient()
	url := fmt.Sprintf("http://server:8080/proxy/prometheus/%s/", p.Cluster.Name)
	return grafanaClient.CreateDataSource(p.Cluster.Name, url)

}
func (p Prometheus) createGrafanaDashboard() error {
	grafanaClient := grafana.NewClient()
	u, err := grafanaClient.CreateDashboard(p.Cluster.Name)
	if err != nil {
		return err
	}
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(p.Tool.Vars), &values)
	values["url"] = u
	str, _ := json.Marshal(&values)
	p.Tool.Vars = string(str)
	return nil
}
