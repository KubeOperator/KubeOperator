package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	RegistryImageName = "registry"
	RegistryTag       = "2.7.1"
)

type Registry struct {
	Cluster       *Cluster
	Tool          *model.ClusterTool
	LocalhostName string
}

func NewRegistry(cluster *Cluster, localhostName string, tool *model.ClusterTool) (*Registry, error) {
	p := &Registry{
		Tool:          tool,
		Cluster:       cluster,
		LocalhostName: localhostName,
	}
	return p, nil
}

func (c Registry) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &values)
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalhostName, constant.LocalDockerRepositoryPort, RegistryImageName)
	values["image.tag"] = RegistryTag

	if _, ok := values["persistence.size"]; ok {
		values["persistence.size"] = fmt.Sprintf("%vGi", values["persistence.size"])
	}

	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
}

func (c Registry) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.DockerRegistryChartName); err != nil {
		return err
	}
	if err := createRoute(c.Cluster.Namespace, constant.DefaultRegistryIngressName, constant.DefaultRegistryIngress, constant.DefaultRegistryServiceName, 5000, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(c.Cluster.Namespace, constant.DefaultRegistryDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c Registry) Uninstall() error {
	return uninstall(c.Cluster.Namespace, c.Tool, constant.DefaultRegistryIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
