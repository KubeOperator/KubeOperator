package tools

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type EFK struct {
	Cluster *Cluster
	Tool    *model.ClusterTool
}

func NewEFK(cluster *Cluster, tool *model.ClusterTool) (*EFK, error) {
	p := &EFK{
		Tool:    tool,
		Cluster: cluster,
	}
	return p, nil
}

func (c EFK) Install() error {
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.EFKChartName); err != nil {
		return err
	}
	if err := createRoute(constant.DefaultEFKIngressName, constant.DefaultEFKIngress, constant.DefaultEFKServiceName, 8080, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(constant.DefaultEFKDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c EFK) Uninstall() error {
	return uninstall(c.Tool, constant.DefaultEFKIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
