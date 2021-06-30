package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Chartmuseum struct {
	Cluster             *Cluster
	Tool                *model.ClusterTool
	LocalHostName       string
	LocalRepositoryPort int
}

func NewChartmuseum(cluster *Cluster, tool *model.ClusterTool) (*Chartmuseum, error) {
	p := &Chartmuseum{
		Tool:                tool,
		Cluster:             cluster,
		LocalHostName:       constant.LocalRepositoryDomainName,
		LocalRepositoryPort: cluster.helmRepoPort,
	}
	return p, nil
}

func (c Chartmuseum) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	versionMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &versionMap)

	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &values)
	values["env.open.DISABLE_API"] = false
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, c.LocalRepositoryPort, versionMap["chartmuseum_image_name"])
	values["image.tag"] = versionMap["chartmuseum_image_tag"]

	if isInstall {
		if va, ok := values["persistence.enabled"]; ok {
			if hasPers, _ := va.(bool); !hasPers {
				delete(values, "nodeSelector.kubernetes\\.io/hostname")
			}
		}

		if _, ok := values["persistence.size"]; ok {
			values["persistence.size"] = fmt.Sprintf("%vGi", values["persistence.size"])
		}
	}

	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
}

func (c Chartmuseum) Install(toolDetail model.ClusterToolDetail) error {
	c.setDefaultValue(toolDetail, true)
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.ChartmuseumChartName, toolDetail.ChartVersion); err != nil {
		return err
	}
	if err := createRoute(c.Cluster.Namespace, constant.DefaultChartmuseumIngressName, constant.DefaultChartmuseumIngress, constant.DefaultChartmuseumServiceName, 8080, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(c.Cluster.Namespace, constant.DefaultChartmuseumDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c Chartmuseum) Upgrade(toolDetail model.ClusterToolDetail) error {
	c.setDefaultValue(toolDetail, false)
	return upgradeChart(c.Cluster.HelmClient, c.Tool, constant.ChartmuseumChartName, toolDetail.ChartVersion)
}

func (c Chartmuseum) Uninstall() error {
	return uninstall(c.Cluster.Namespace, c.Tool, constant.DefaultChartmuseumIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
