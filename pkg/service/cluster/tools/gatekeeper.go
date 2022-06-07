package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Gatekeeper struct {
	Cluster             *Cluster
	Tool                *model.ClusterTool
	LocalHostName       string
	LocalRepositoryPort int
}

func NewGatekeeper(cluster *Cluster, tool *model.ClusterTool) (*Gatekeeper, error) {
	p := &Gatekeeper{
		Tool:                tool,
		Cluster:             cluster,
		LocalHostName:       constant.LocalRepositoryDomainName,
		LocalRepositoryPort: cluster.helmRepoPort,
	}
	return p, nil
}

func (k Gatekeeper) setDefaultValue(toolDetail model.ClusterToolDetail) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)

	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(k.Tool.Vars), &values)
	values["postInstall.labelNamespace.image.repository"] = fmt.Sprintf("%s:%d/%s", k.LocalHostName, k.LocalRepositoryPort, imageMap["post_image_repo"])
	values["postInstall.labelNamespace.image.tag"] = imageMap["post_image_tag"]
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", k.LocalHostName, k.LocalRepositoryPort, imageMap["image_repo"])
	values["image.crdRepository"] = fmt.Sprintf("%s:%d/%s", k.LocalHostName, k.LocalRepositoryPort, imageMap["crd_image_repo"])
	values["image.release"] = imageMap["image_release"]
	str, _ := json.Marshal(&values)
	k.Tool.Vars = string(str)
}

func (k Gatekeeper) Install(toolDetail model.ClusterToolDetail) error {
	k.setDefaultValue(toolDetail)
	if err := installChart(k.Cluster.HelmClient, k.Tool, constant.GatekeeperChartName, toolDetail.ChartVersion); err != nil {
		return err
	}

	ingressItem := &Ingress{
		name:    constant.DefaultGatekeeperIngressName,
		url:     constant.DefaultGatekeeperIngress,
		service: constant.DefaultGatekeeperServiceName,
		port:    80,
		version: k.Cluster.Version,
	}
	if err := createRoute(k.Cluster.Namespace, ingressItem, k.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(k.Cluster.Namespace, constant.DefaultGatekeeperDeploymentName, 1, k.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (k Gatekeeper) Upgrade(toolDetail model.ClusterToolDetail) error {
	k.setDefaultValue(toolDetail)
	return upgradeChart(k.Cluster.HelmClient, k.Tool, constant.GatekeeperChartName, toolDetail.ChartVersion)
}

func (k Gatekeeper) Uninstall() error {
	return uninstall(k.Cluster.Namespace, k.Tool, constant.DefaultGatekeeperIngressName, k.Cluster.Version, k.Cluster.HelmClient, k.Cluster.KubeClient)
}
