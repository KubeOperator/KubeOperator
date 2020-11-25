package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	LokiImageName     = "grafana/loki"
	LokiTag           = "2.0.0"
	PromtailImageName = "grafana/promtail"
	PromtailTag       = "2.0.0"
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
	values["loki.image.tag"] = LokiTag
	values["promtail.image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, constant.LocalDockerRepositoryPort, PromtailImageName)
	values["promtail.image.tag"] = PromtailTag

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
	if err := createRoute(constant.DefaultLokiIngressName, constant.DefaultLokiIngress, constant.DefaultLokiServiceName, 3100, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForStatefulSetsRunning(constant.DefaultLokiStateSetsfulName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c Loki) Uninstall() error {
	return uninstall(c.Tool, constant.DefaultLokiIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
