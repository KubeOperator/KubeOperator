package tools

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Chartmuseum struct {
	Cluster *Cluster
	Tool    *model.ClusterTool
}

func NewChartmuseum(cluster *Cluster, tool *model.ClusterTool) (*Chartmuseum, error) {
	p := &Chartmuseum{
		Tool:    tool,
		Cluster: cluster,
	}
	return p, nil
}

func (c Chartmuseum) Install() error {
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.ChartmuseumChartName); err != nil {
		return err
	}
	if err := createRoute(constant.DefaultChartmuseumIngressName, constant.DefaultChartmuseumIngress, constant.DefaultChartmuseumServiceName, 8080, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(constant.DefaultChartmuseumDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c Chartmuseum) Uninstall() error {
	return uninstall(c.Tool, constant.DefaultChartmuseumIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
