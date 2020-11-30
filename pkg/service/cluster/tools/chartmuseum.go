package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	ChartmuseumImageNameAmd64Name = "chartmuseum/chartmuseum"
	ChartmuseumTagAmd64Name       = "v0.12.0"
	ChartmuseumImageNameArm64Name = "kubeoperator/chartmuseum"
	ChartmuseumTagArm64Name       = "v0.12.0-arm64"
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

	if c.Cluster.Spec.Architectures == "amd64" {
		values["image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalhostName, constant.LocalDockerRepositoryPort, ChartmuseumImageNameAmd64Name)
		values["image.tag"] = ChartmuseumTagAmd64Name
	} else {
		values["image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalhostName, constant.LocalDockerRepositoryPort, ChartmuseumImageNameArm64Name)
		values["image.tag"] = ChartmuseumTagArm64Name
	}

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
