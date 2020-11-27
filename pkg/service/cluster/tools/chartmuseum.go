package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	ChartmuseumImageName = "chartmuseum/chartmuseum"
	ChartmuseumTag       = "v0.12.0"
)

type Chartmuseum struct {
	Cluster       *Cluster
	Tool          *model.ClusterTool
	LocalhostName string
}

func NewChartmuseum(cluster *Cluster, localhostName string, tool *model.ClusterTool) (*Chartmuseum, error) {
	p := &Chartmuseum{
		Tool:          tool,
		Cluster:       cluster,
		LocalhostName: localhostName,
	}
	return p, nil
}

func (c Chartmuseum) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &values)
	values["env.open.DISABLE_API"] = false
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalhostName, constant.LocalDockerRepositoryPort, ChartmuseumImageName)
	values["image.tag"] = ChartmuseumTag
	if _, ok := values["persistence.size"]; ok {
		values["persistence.size"] = fmt.Sprintf("%vGi", values["persistence.size"])
	}

	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
}

func (c Chartmuseum) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.ChartmuseumChartName); err != nil {
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

func (c Chartmuseum) Uninstall() error {
	return uninstall(c.Cluster.Namespace, c.Tool, constant.DefaultChartmuseumIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
