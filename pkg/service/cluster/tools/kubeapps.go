package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	helm2 "github.com/KubeOperator/KubeOperator/pkg/util/helm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		LocalRepositoryPort: cluster.helmRepoPort,
	}
	return p, nil
}

func (k Kubeapps) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)
	values := map[string]interface{}{}
	switch toolDetail.ChartVersion {
	case "3.7.2":
		values = k.valuseV372Binding(imageMap)
	case "5.0.1":
		values = k.valuseV501Binding(imageMap)
	case "7.6.2":
		values = k.valuseV762Binding(imageMap, isInstall)
	}
	if isInstall {
		var c helm2.Client
		repoIP, _, repoPort, _, _ := c.GetRepoIP("amd64")
		values["apprepository.initialRepos[0].name"] = "kubeoperator"
		values["apprepository.initialRepos[0].url"] = fmt.Sprintf("http://%s:%d/repository/kubeapps", repoIP, repoPort)

		if va, ok := values["postgresql.persistence.enabled"]; ok {
			if hasPers, _ := va.(bool); hasPers {
				if va, ok := values["nodeSelector"]; ok {
					values["postgresql.primary.nodeSelector.kubernetes\\.io/hostname"] = va
				}
			}
		}
		if _, ok := values["postgresql.persistence.size"]; ok {
			values["postgresql.persistence.size"] = fmt.Sprintf("%vGi", values["postgresql.persistence.size"])
		}
		delete(values, "nodeSelector")
	}

	str, _ := json.Marshal(&values)
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
	_ = json.Unmarshal([]byte(k.Tool.Vars), &values)
	values["global.imageRegistry"] = fmt.Sprintf("%s:%d", k.LocalHostName, k.LocalRepositoryPort)
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
		_ = json.Unmarshal([]byte(k.Tool.Vars), &values)
	}
	delete(values, "useHelm3")
	delete(values, "postgresql.enabled")
	delete(values, "postgresql.image.repository")
	delete(values, "postgresql.image.tag")
	values["global.imageRegistry"] = fmt.Sprintf("%s:%d", k.LocalHostName, k.LocalRepositoryPort)

	return values
}

// v7.6.2
func (k Kubeapps) valuseV762Binding(imageMap map[string]interface{}, isInstall bool) map[string]interface{} {
	values := map[string]interface{}{}
	if len(k.Tool.Vars) != 0 {
		_ = json.Unmarshal([]byte(k.Tool.Vars), &values)
	}

	values["global.imageRegistry"] = fmt.Sprintf("%s:%d", k.LocalHostName, k.LocalRepositoryPort)

	if !isInstall {
		delete(values, "apprepository.initialRepos[0].name")
		delete(values, "apprepository.initialRepos[0].url")

		if err := k.Cluster.KubeClient.AppsV1().Deployments(k.Cluster.Namespace).Delete(context.TODO(), "kubeapps-internal-apprepository-controller", metav1.DeleteOptions{}); err != nil {
			logger.Log.Infof("delete deployment kubeapps-internal-apprepository-controller from %s failed, err: %v", k.Cluster.Namespace, err)
		}
		if err := k.Cluster.KubeClient.AppsV1().Deployments(k.Cluster.Namespace).Delete(context.TODO(), "kubeapps", metav1.DeleteOptions{}); err != nil {
			logger.Log.Infof("delete deployment kubeapps-internal-apprepository-controller from %s failed, err: %v", k.Cluster.Namespace, err)
		}
		if err := k.Cluster.KubeClient.AppsV1().Deployments(k.Cluster.Namespace).Delete(context.TODO(), "kubeapps-internal-assetsvc", metav1.DeleteOptions{}); err != nil {
			logger.Log.Infof("delete deployment kubeapps-internal-assetsvc from %s failed, err: %v", k.Cluster.Namespace, err)
		}
		if err := k.Cluster.KubeClient.AppsV1().Deployments(k.Cluster.Namespace).Delete(context.TODO(), "kubeapps-internal-dashboard", metav1.DeleteOptions{}); err != nil {
			logger.Log.Info("delete deploymentkubeapps-internal-assetsvc from %s failed, err: %v", k.Cluster.Namespace, err)
		}
		if err := k.Cluster.KubeClient.AppsV1().Deployments(k.Cluster.Namespace).Delete(context.TODO(), "kubeapps-internal-kubeops", metav1.DeleteOptions{}); err != nil {
			logger.Log.Info("delete deployment kubeapps-internal-kubeops from %s failed, err: %v", k.Cluster.Namespace, err)
		}

		postgresqlSecret, err := k.Cluster.KubeClient.CoreV1().Secrets(k.Cluster.Namespace).Get(context.TODO(), "kubeapps-db", metav1.GetOptions{})
		if err != nil {
			logger.Log.Info("get kubeapps-db secrets from %s failed, err: %v", k.Cluster.Namespace, err)
			return values
		}
		password := postgresqlSecret.Data["postgresql-password"]
		values["postgresql.postgresqlPassword"] = string(password)
	}

	return values
}
