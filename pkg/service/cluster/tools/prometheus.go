package tools

import (
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/grafana"
)

type Prometheus struct {
	Tool    *model.ClusterTool
	Cluster *Cluster
}

func NewPrometheus(cluster *Cluster, tool *model.ClusterTool) (*Prometheus, error) {
	p := &Prometheus{
		Tool:    tool,
		Cluster: cluster,
	}
	return p, nil
}

func (p Prometheus) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(p.Tool.Vars), &values)
	values["alertmanager.enabled"] = false
	values["pushgateway.enabled"] = false
	str, _ := json.Marshal(&values)
	p.Tool.Vars = string(str)
}

func (c Prometheus) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.PrometheusChartName); err != nil {
		return err
	}
	if err := createRoute(constant.DefaultPrometheusIngressName, constant.DefaultPrometheusIngress, constant.DefaultPrometheusServiceName, 80, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(constant.DefaultPrometheusDeploymentName, 1, c.Cluster.KubeClient); err != nil {
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
	return uninstall(c.Tool, constant.DefaultPrometheusIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
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
