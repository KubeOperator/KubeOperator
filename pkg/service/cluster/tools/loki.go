package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	LokiImageName    = "grafana/loki"
	LokiTagAmd64Name = "2.0.0-amd64"
	LokiTagArm64Name = "2.0.0-arm64"

	PromtailImageName    = "grafana/promtail"
	PromtailTagAmd64Name = "2.0.0-amd64"
	PromtailTagArm64Name = "2.0.0-arm64"
)

type Loki struct {
	Cluster       *Cluster
	Tool          *model.ClusterTool
	LocalHostName string
}

func (c Loki) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &values)
	values["loki.image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, constant.LocalDockerRepositoryPort, LokiImageName)
	values["promtail.image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, constant.LocalDockerRepositoryPort, PromtailImageName)

	if c.Cluster.Spec.Architectures == "amd64" {
		values["loki.image.tag"] = LokiTagAmd64Name
		values["promtail.image.tag"] = PromtailTagAmd64Name
	} else {
		values["loki.image.tag"] = LokiTagArm64Name
		values["promtail.image.tag"] = PromtailTagArm64Name
	}

	if _, ok := values["loki.persistence.size"]; ok {
		values["loki.persistence.size"] = fmt.Sprintf("%vGi", values["loki.persistence.size"])
	}
	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
}

func NewLoki(cluster *Cluster, localhostName string, tool *model.ClusterTool) (*Loki, error) {
	p := &Loki{
		Tool:          tool,
		Cluster:       cluster,
		LocalHostName: localhostName,
	}
	return p, nil
}

func (c Loki) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.LokiChartName); err != nil {
		return err
	}
	if err := createRoute(c.Cluster.Namespace, constant.DefaultLokiIngressName, constant.DefaultLokiIngress, constant.DefaultLokiServiceName, 3100, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForStatefulSetsRunning(c.Cluster.Namespace, constant.DefaultLokiStateSetsfulName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c Loki) Uninstall() error {
	return uninstall(c.Cluster.Namespace, c.Tool, constant.DefaultLokiIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
