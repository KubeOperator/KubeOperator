package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Registry struct {
	Cluster             *Cluster
	Tool                *model.ClusterTool
	LocalHostName       string
	LocalRepositoryPort int
}

func NewRegistry(cluster *Cluster, tool *model.ClusterTool) (*Registry, error) {
	p := &Registry{
		Tool:                tool,
		Cluster:             cluster,
		LocalHostName:       constant.LocalRepositoryDomainName,
		LocalRepositoryPort: cluster.helmRepoPort,
	}
	return p, nil
}

func (r Registry) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)

	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(r.Tool.Vars), &values)
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", r.LocalHostName, r.LocalRepositoryPort, imageMap["registry_image_name"])
	values["image.tag"] = imageMap["registry_image_tag"]
	values["secrets.htpasswd"] = "admin:$2y$05$xOL4vcb.1gGpKBHzW0Vv0O4KV0kOAHLXkXBPHtZFAswoW.hYVGzOy"

	if isInstall {
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
	r.Tool.Vars = string(str)
}

func (r Registry) Install(toolDetail model.ClusterToolDetail) error {
	r.setDefaultValue(toolDetail, true)
	if err := installChart(r.Cluster.HelmClient, r.Tool, constant.DockerRegistryChartName, toolDetail.ChartVersion); err != nil {
		return err
	}

	ingressItem := &Ingress{
		name:    constant.DefaultRegistryIngressName,
		url:     constant.DefaultRegistryIngress,
		service: constant.DefaultRegistryServiceName,
		port:    5000,
		version: r.Cluster.Spec.Version,
	}
	if err := createRoute(r.Cluster.Namespace, ingressItem, r.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(r.Cluster.Namespace, constant.DefaultRegistryDeploymentName, 1, r.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (r Registry) Upgrade(toolDetail model.ClusterToolDetail) error {
	r.setDefaultValue(toolDetail, false)
	return upgradeChart(r.Cluster.HelmClient, r.Tool, constant.DockerRegistryChartName, toolDetail.ChartVersion)
}

func (r Registry) Uninstall() error {
	return uninstall(r.Cluster.Namespace, r.Tool, constant.DefaultRegistryIngressName, r.Cluster.HelmClient, r.Cluster.KubeClient)
}
