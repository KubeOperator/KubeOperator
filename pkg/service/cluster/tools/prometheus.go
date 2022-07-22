package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		LocalRepositoryPort: cluster.helmRepoPort,
	}
	return p, nil
}

func (p Prometheus) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)
	values := map[string]interface{}{}
	switch toolDetail.ChartVersion {
	case "11.12.1", "11.5.0":
		values = p.valuse11121Binding(imageMap)
	case "15.0.1", "15.10.1":
		values = p.valuse1501Binding(imageMap, isInstall)
	}

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

	ingressItem := &Ingress{
		name:    constant.DefaultPrometheusIngressName,
		url:     constant.DefaultPrometheusIngress,
		service: constant.DefaultPrometheusServiceName,
		port:    80,
		version: p.Cluster.Version,
	}
	if err := createRoute(p.Cluster.Namespace, ingressItem, p.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(p.Cluster.Namespace, constant.DefaultPrometheusDeploymentName, 1, p.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (p Prometheus) Upgrade(toolDetail model.ClusterToolDetail) error {
	p.setDefaultValue(toolDetail, false)
	return upgradeChart(p.Cluster.HelmClient, p.Tool, constant.PrometheusChartName, toolDetail.ChartVersion)
}

func (p Prometheus) Uninstall() error {
	return uninstall(p.Cluster.Namespace, p.Tool, constant.DefaultPrometheusIngressName, p.Cluster.Version, p.Cluster.HelmClient, p.Cluster.KubeClient)
}

// 11.12.1
func (p Prometheus) valuse11121Binding(imageMap map[string]interface{}) map[string]interface{} {
	values := map[string]interface{}{}
	if len(p.Tool.Vars) != 0 {
		_ = json.Unmarshal([]byte(p.Tool.Vars), &values)
	}
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

	return values
}

// 15.0.1
func (p Prometheus) valuse1501Binding(imageMap map[string]interface{}, isInstall bool) map[string]interface{} {
	values := map[string]interface{}{}
	if len(p.Tool.Vars) != 0 {
		_ = json.Unmarshal([]byte(p.Tool.Vars), &values)
	}

	if !isInstall {
		if err := p.Cluster.KubeClient.AppsV1().Deployments(p.Cluster.Namespace).Delete(context.TODO(), "prometheus-kube-state-metrics", metav1.DeleteOptions{}); err != nil {
			logger.Log.Infof("delete deployment kubeapps-internal-apprepository-controller from %s failed, err: %v", p.Cluster.Namespace, err)
		}
	}

	values["alertmanager.enabled"] = false
	values["pushgateway.enabled"] = false

	values["configmapReload.prometheus.enabled"] = true
	values["configmapReload.prometheus.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalHostName, p.LocalRepositoryPort, imageMap["configmap_image_name"])
	values["configmapReload.prometheus.image.tag"] = imageMap["configmap_image_tag"]
	values["kube-state-metrics.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalHostName, p.LocalRepositoryPort, imageMap["metrics_image_name"])
	values["kube-state-metrics.image.tag"] = imageMap["metrics_image_tag"]

	values["nodeExporter.enabled"] = true
	values["nodeExporter.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalHostName, p.LocalRepositoryPort, imageMap["exporter_image_name"])
	values["nodeExporter.image.tag"] = imageMap["exporter_image_tag"]

	values["server.enabled"] = true
	values["server.service.type"] = "NodePort"
	values["server.image.repository"] = fmt.Sprintf("%s:%d/%s", p.LocalHostName, p.LocalRepositoryPort, imageMap["prometheus_image_name"])
	values["server.image.tag"] = imageMap["prometheus_image_tag"]

	return values
}
