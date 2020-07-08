package tools

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Registry struct {
	Cluster *Cluster
	Tool    *model.ClusterTool
}

func NewRegistry(cluster *Cluster, tool *model.ClusterTool) (*Registry, error) {
	p := &Registry{
		Tool:    tool,
		Cluster: cluster,
	}
	return p, nil
}

func (c Registry) Install() error {
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.DockerRegistryChartName); err != nil {
		return err
	}
	if err := createRoute(constant.DefaultRegistryIngressName, constant.DefaultRegistryIngress, constant.DefaultRegistryServiceName, 5000, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(constant.DefaultRegistryDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c Registry) Uninstall() error {
	return uninstall(c.Tool, constant.DefaultRegistryIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
