package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

var log = logger.Default

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
		LocalRepositoryPort: constant.LocalDockerRepositoryPort,
	}
	return p, nil
}

func (c Chartmuseum) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	versionMap := map[string]interface{}{}
	if len(toolDetail.Vars) != 0 {
		if err := json.Unmarshal([]byte(toolDetail.Vars), &versionMap); err != nil {
			log.Errorf("json unmarshal falied : %v", toolDetail.Vars)
		}
	}

	values := map[string]interface{}{}
	if len(c.Tool.Vars) != 0 {
		if err := json.Unmarshal([]byte(c.Tool.Vars), &values); err != nil {
			log.Errorf("json unmarshal falied : %v", c.Tool.Vars)
		}
	}
	values["env.open.DISABLE_API"] = false
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, c.LocalRepositoryPort, versionMap["chartmuseum_image_name"])
	values["image.tag"] = versionMap["chartmuseum_image_tag"]

	if isInstall {
		if va, ok := values["persistence.enabled"]; ok {
			hasPers, ok := va.(bool)
			if ok {
				if !hasPers {
					delete(values, "nodeSelector.kubernetes\\.io/hostname")
				}
			}
		}

		if _, ok := values["persistence.size"]; ok {
			values["persistence.size"] = fmt.Sprintf("%vGi", values["persistence.size"])
		}
	}

	str, err := json.Marshal(&values)
	if err != nil {
		log.Errorf("json marshal falied : %v", values)
	}
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
