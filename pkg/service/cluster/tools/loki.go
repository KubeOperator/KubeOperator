package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Loki struct {
	Cluster             *Cluster
	Tool                *model.ClusterTool
	LocalHostName       string
	LocalRepositoryPort int
}

func NewLoki(cluster *Cluster, tool *model.ClusterTool) (*Loki, error) {
	p := &Loki{
		Tool:                tool,
		Cluster:             cluster,
		LocalHostName:       constant.LocalRepositoryDomainName,
		LocalRepositoryPort: cluster.helmRepoPort,
	}
	return p, nil
}

func (l Loki) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)

	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(l.Tool.Vars), &values)
	values["loki.image.repository"] = fmt.Sprintf("%s:%d/%s", l.LocalHostName, l.LocalRepositoryPort, imageMap["loki_image_name"])
	values["promtail.image.repository"] = fmt.Sprintf("%s:%d/%s", l.LocalHostName, l.LocalRepositoryPort, imageMap["promtail_image_name"])
	values["loki.image.tag"] = imageMap["loki_image_tag"]
	values["promtail.image.tag"] = imageMap["promtail_image_tag"]

	if isInstall {
		if _, ok := values["loki.persistence.size"]; ok {
			values["loki.persistence.size"] = fmt.Sprintf("%vGi", values["loki.persistence.size"])
		}
		if va, ok := values["loki.persistence.enabled"]; ok {
			if hasPers, _ := va.(bool); !hasPers {
				delete(values, "loki.nodeSelector.kubernetes\\.io/hostname")
			}
		}
	}

	str, _ := json.Marshal(&values)
	l.Tool.Vars = string(str)
}

func (l Loki) Install(toolDetail model.ClusterToolDetail) error {
	l.setDefaultValue(toolDetail, true)
	if err := installChart(l.Cluster.HelmClient, l.Tool, constant.LokiChartName, toolDetail.ChartVersion); err != nil {
		return err
	}

	ingressItem := &Ingress{
		name:    constant.DefaultLokiIngressName,
		url:     constant.DefaultLokiIngress,
		service: constant.DefaultLokiServiceName,
		port:    3100,
		version: l.Cluster.Spec.Version,
	}
	if err := createRoute(l.Cluster.Namespace, ingressItem, l.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForStatefulSetsRunning(l.Cluster.Namespace, constant.DefaultLokiStateSetsfulName, 1, l.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (l Loki) Upgrade(toolDetail model.ClusterToolDetail) error {
	l.setDefaultValue(toolDetail, false)
	return upgradeChart(l.Cluster.HelmClient, l.Tool, constant.LokiChartName, toolDetail.ChartVersion)
}

func (l Loki) Uninstall() error {
	return uninstall(l.Cluster.Namespace, l.Tool, constant.DefaultLokiIngressName, l.Cluster.HelmClient, l.Cluster.KubeClient)
}
