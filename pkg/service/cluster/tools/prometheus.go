package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/grafana"
)

type Prometheus struct {
	Tool                *model.ClusterTool
	Cluster             *Cluster
	LocalHostName       string
	LocalRepositoryPort int
}

func NewPrometheus(cluster *Cluster, tool *model.ClusterTool) (*Prometheus, error) {
	p := &Prometheus{
		Tool:                tool,
		Cluster:             cluster,
		LocalHostName:       constant.LocalRepositoryDomainName,
		LocalRepositoryPort: constant.LocalDockerRepositoryPort,
	}
	return p, nil
}

func (p Prometheus) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)

	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(p.Tool.Vars), &values)
	values["alertmanager.enabled"] = false
	values["pushgateway.enabled"] = false
	values["configmapReload.prometheus.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalHostName, p.LocalRepositoryPort, imageMap["configmap_image_name"])
	values["configmapReload.prometheus.image.tag"] = imageMap["configmap_image_tag"]
	values["kube-state-metrics.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalHostName, p.LocalRepositoryPort, imageMap["metrics_image_name"])
	values["kube-state-metrics.image.tag"] = imageMap["metrics_image_tag"]
	values["nodeExporter.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalHostName, p.LocalRepositoryPort, imageMap["exporter_image_name"])
	values["nodeExporter.image.tag"] = imageMap["exporter_image_tag"]
	values["server.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalHostName, p.LocalRepositoryPort, imageMap["prometheus_image_name"])
	values["server.image.tag"] = imageMap["prometheus_image_tag"]

	if isInstall {
		if _, ok := values["server.retention"]; ok {
			values["server.retention"] = fmt.Sprintf("%vd", values["server.retention"])
		}
		if _, ok := values["server.persistentVolume.size"]; ok {
			values["server.persistentVolume.size"] = fmt.Sprintf("%vGi", values["server.persistentVolume.size"])
		}
		if va, ok := values["server.persistentVolume.enabled"]; ok {
			if hasPers, _ := va.(bool); !hasPers {
				delete(values, "server.nodeSelector.kubernetes\\.io/hostname")
			}
		}
	}
	str, _ := json.Marshal(&values)
	p.Tool.Vars = string(str)
}

func (p Prometheus) Install(toolDetail model.ClusterToolDetail) error {
	p.setDefaultValue(toolDetail, true)
	if err := installChart(p.Cluster.HelmClient, p.Tool, constant.PrometheusChartName, toolDetail.ChartVersion); err != nil {
		return err
	}
	if err := createRoute(p.Cluster.Namespace, constant.DefaultPrometheusIngressName, constant.DefaultPrometheusIngress, constant.DefaultPrometheusServiceName, 80, p.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(p.Cluster.Namespace, constant.DefaultPrometheusDeploymentName, 1, p.Cluster.KubeClient); err != nil {
		return err
	}

	if err := p.createGrafanaDataSource(); err != nil {
		return err
	}

	if err := p.createGrafanaDashboard(); err != nil {
		return err
	}
	return nil
}

func (p Prometheus) Upgrade(toolDetail model.ClusterToolDetail) error {
	p.setDefaultValue(toolDetail, false)
	return upgradeChart(p.Cluster.HelmClient, p.Tool, constant.PrometheusChartName, toolDetail.ChartVersion)
}

func (p Prometheus) Uninstall() error {
	gClient := grafana.NewClient()
	if err := gClient.DeleteDashboard(p.Cluster.Name); err != nil {
		return err
	}
	if err := gClient.DeleteDataSource(p.Cluster.Name); err != nil {
		return err
	}
	return uninstall(p.Cluster.Namespace, p.Tool, constant.DefaultPrometheusIngressName, p.Cluster.HelmClient, p.Cluster.KubeClient)
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
