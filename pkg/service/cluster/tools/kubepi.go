package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type KubePi struct {
	Cluster             *Cluster
	Tool                *model.ClusterTool
	LocalHostName       string
	LocalRepositoryPort int
}

func NewKubePi(cluster *Cluster, tool *model.ClusterTool) (*KubePi, error) {
	p := &KubePi{
		Tool:                tool,
		Cluster:             cluster,
		LocalHostName:       constant.LocalRepositoryDomainName,
		LocalRepositoryPort: cluster.helmRepoPort,
	}
	return p, nil
}

func (k KubePi) setDefaultValue(toolDetail model.ClusterToolDetail) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)

	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(k.Tool.Vars), &values)
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", k.LocalHostName, k.LocalRepositoryPort, imageMap["kubepi_image_name"])
	values["image.tag"] = imageMap["kubepi_image_tag"]
	str, _ := json.Marshal(&values)
	k.Tool.Vars = string(str)
}

func (k KubePi) Install(toolDetail model.ClusterToolDetail) error {
	k.setDefaultValue(toolDetail)
	if err := installChart(k.Cluster.HelmClient, k.Tool, constant.KubePiChartName, toolDetail.ChartVersion); err != nil {
		return err
	}
	if err := createRoute(k.Cluster.Namespace, constant.DefaultKubePiIngressName, constant.DefaultKubePiIngress, constant.DefaultKubePiServiceName, 80, k.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(k.Cluster.Namespace, constant.DefaultKubePiDeploymentName, 1, k.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (k KubePi) Upgrade(toolDetail model.ClusterToolDetail) error {
	k.setDefaultValue(toolDetail)
	return upgradeChart(k.Cluster.HelmClient, k.Tool, constant.KubePiChartName, toolDetail.ChartVersion)
}

func (k KubePi) Uninstall() error {
	return uninstall(k.Cluster.Namespace, k.Tool, constant.DefaultKubePiIngressName, k.Cluster.HelmClient, k.Cluster.KubeClient)
}
