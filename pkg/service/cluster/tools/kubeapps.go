package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	helm2 "github.com/KubeOperator/KubeOperator/pkg/util/helm"
)

type Kubeapps struct {
	Tool                *model.ClusterTool
	Cluster             *Cluster
	LocalHostName       string
	LocalRepositoryPort int
}

func NewKubeapps(cluster *Cluster, tool *model.ClusterTool) (*Kubeapps, error) {
	p := &Kubeapps{
		Tool:                tool,
		Cluster:             cluster,
		LocalHostName:       constant.LocalRepositoryDomainName,
		LocalRepositoryPort: constant.LocalDockerRepositoryPort,
	}
	return p, nil
}

func (k Kubeapps) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	imageMap := map[string]interface{}{}
	if len(toolDetail.Vars) != 0 {
		if err := json.Unmarshal([]byte(toolDetail.Vars), &imageMap); err != nil {
			log.Errorf("json unmarshal falied : %v", toolDetail.Vars)
		}
	}
	values := map[string]interface{}{}
	switch toolDetail.ChartVersion {
	case "3.7.2":
		values = k.valuseV372Binding(imageMap)
	case "5.0.1":
		values = k.valuseV501Binding(imageMap)
	}
	if isInstall {
		if va, ok := values["postgresql.persistence.enabled"]; ok {
			hasPers, ok := va.(bool)
			if ok {
				if hasPers {
					if va, ok := values["nodeSelector"]; ok {
						values["postgresql.primary.nodeSelector.kubernetes\\.io/hostname"] = va
					}
				}
			}
		}
		if _, ok := values["postgresql.persistence.size"]; ok {
			values["postgresql.persistence.size"] = fmt.Sprintf("%vGi", values["postgresql.persistence.size"])
		}
		delete(values, "nodeSelector")
	}

	str, err := json.Marshal(&values)
	if err != nil {
		log.Errorf("json marshal falied : %v", values)
	}
	k.Tool.Vars = string(str)
}

func (k Kubeapps) Install(toolDetail model.ClusterToolDetail) error {
	k.setDefaultValue(toolDetail, true)
	if err := installChart(k.Cluster.HelmClient, k.Tool, constant.KubeappsChartName, toolDetail.ChartVersion); err != nil {
		return err
	}
	if err := createRoute(k.Cluster.Namespace, constant.DefaultKubeappsIngressName, constant.DefaultKubeappsIngress, constant.DefaultKubeappsServiceName, 80, k.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(k.Cluster.Namespace, constant.DefaultKubeappsDeploymentName, 1, k.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (k Kubeapps) Upgrade(toolDetail model.ClusterToolDetail) error {
	k.setDefaultValue(toolDetail, false)
	return upgradeChart(k.Cluster.HelmClient, k.Tool, constant.KubeappsChartName, toolDetail.ChartVersion)
}

func (k Kubeapps) Uninstall() error {
	return uninstall(k.Cluster.Namespace, k.Tool, constant.DefaultKubeappsIngressName, k.Cluster.HelmClient, k.Cluster.KubeClient)
}

// v3.7.2
func (k Kubeapps) valuseV372Binding(imageMap map[string]interface{}) map[string]interface{} {
	values := map[string]interface{}{}
	if len(k.Tool.Vars) != 0 {
		if err := json.Unmarshal([]byte(k.Tool.Vars), &values); err != nil {
			log.Errorf("json unmarshal falied : %v", (k.Tool.Vars))
		}
	}
	var c helm2.Client
	_, repoIP, _ := c.GetRepoIP("amd64")
	values["global.imageRegistry"] = fmt.Sprintf("%s:%d", k.LocalHostName, k.LocalRepositoryPort)
	values["apprepository.initialRepos[0].name"] = "kubeoperator"
	values["apprepository.initialRepos[0].url"] = fmt.Sprintf("https://%s:8081/repository/kubeapps", repoIP)
	values["useHelm3"] = true
	values["postgresql.enabled"] = true
	values["postgresql.image.repository"] = imageMap["postgresql_image_name"]
	values["postgresql.image.tag"] = imageMap["postgresql_image_tag"]

	return values
}

// v5.0.1
func (k Kubeapps) valuseV501Binding(imageMap map[string]interface{}) map[string]interface{} {
	values := map[string]interface{}{}
	if len(k.Tool.Vars) != 0 {
		if err := json.Unmarshal([]byte(k.Tool.Vars), &values); err != nil {
			log.Errorf("json unmarshal falied : %v", (k.Tool.Vars))
		}
	}
	var c helm2.Client
	_, repoIP, _ := c.GetRepoIP("amd64")
	delete(values, "useHelm3")
	delete(values, "postgresql.enabled")
	delete(values, "postgresql.image.repository")
	delete(values, "postgresql.image.tag")
	values["global.imageRegistry"] = fmt.Sprintf("%s:%d", k.LocalHostName, k.LocalRepositoryPort)
	values["apprepository.initialRepos[0].name"] = "kubeoperator"
	values["apprepository.initialRepos[0].url"] = fmt.Sprintf("https://%s:8081/repository/kubeapps", repoIP)

	return values
}
