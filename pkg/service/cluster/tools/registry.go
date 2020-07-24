package tools

import (
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
)

const (
	RegistryImageName = "registry"
	RegistryTag       = "2.7.1"
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

func (c Registry) setDefaultValue() {
	systemService := service.NewSystemSettingService()
	locahostName := systemService.GetLocalHostName()
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &values)
	values["image.repository"] = fmt.Sprintf("%s:%d/%s", locahostName, constant.LocalDockerRepositoryPort, RegistryImageName)
	values["image.tag"] = RegistryTag
	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
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
