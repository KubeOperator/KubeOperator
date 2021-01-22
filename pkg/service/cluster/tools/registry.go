package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Registry struct {
	Cluster       *Cluster
	Tool          *model.ClusterTool
	LocalhostName string
}

func NewRegistry(cluster *Cluster, localhostName string, tool *model.ClusterTool) (*Registry, error) {
	p := &Registry{
		Tool:          tool,
		Cluster:       cluster,
		LocalhostName: localhostName,
	}
	return p, nil
}

func (r Registry) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)

	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(r.Tool.Vars), &values)
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", r.LocalhostName, constant.LocalDockerRepositoryPort, imageMap["registry_image_name"])
	values["image.tag"] = imageMap["registry_image_tag"]

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
	if err := createRoute(r.Cluster.Namespace, constant.DefaultRegistryIngressName, constant.DefaultRegistryIngress, constant.DefaultRegistryServiceName, 5000, r.Cluster.KubeClient); err != nil {
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
