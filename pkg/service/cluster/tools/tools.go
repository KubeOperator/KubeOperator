package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/db"

	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

type Interface interface {
	Install(toolDetail model.ClusterToolDetail) error
	Upgrade(toolDetail model.ClusterToolDetail) error
	Uninstall() error
}

type Cluster struct {
	OldNamespace string
	Namespace    string
	model.Cluster
	HelmClient helm.Interface
	KubeClient *kubernetes.Clientset
}

func NewCluster(cluster model.Cluster, hosts []kubernetesUtil.Host, secret model.ClusterSecret, oldNamespace, namespace string) (*Cluster, error) {
	c := Cluster{
		Cluster: cluster,
	}
	c.Namespace = namespace
	helmClient, err := helm.NewClient(&helm.Config{
		Hosts:         hosts,
		BearerToken:   secret.KubernetesToken,
		OldNamespace:  oldNamespace,
		Namespace:     namespace,
		Architectures: cluster.Spec.Architectures,
	})
	if err != nil {
		return nil, err
	}
	c.HelmClient = helmClient
	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Hosts: hosts,
		Token: secret.KubernetesToken,
	})
	if err != nil {
		return nil, err
	}
	c.KubeClient = kubeClient
	return &c, nil
}

func NewClusterTool(tool *model.ClusterTool, cluster model.Cluster, hosts []kubernetesUtil.Host, secret model.ClusterSecret, oldNamespace, namespace string, enable bool) (Interface, error) {
	c, err := NewCluster(cluster, hosts, secret, oldNamespace, namespace)
	if err != nil {
		return nil, err
	}
	switch tool.Name {
	case "prometheus":
		return NewPrometheus(c, tool)
	case "logging":
		return NewEFK(c, tool)
	case "loki":
		return NewLoki(c, tool)
	case "grafana":
		if enable {
			prometheusNs, err := getGrafanaSourceNs(cluster, "prometheus")
			if err != nil {
				return nil, err
			}
			lokiNs, err := getGrafanaSourceNs(cluster, "loki")
			if err != nil {
				return nil, err
			}
			return NewGrafana(c, tool, prometheusNs, lokiNs)
		} else {
			return NewGrafana(c, tool, "", "")
		}
	case "registry":
		return NewRegistry(c, tool)
	case "dashboard":
		return NewDashboard(c, tool)
	case "chartmuseum":
		return NewChartmuseum(c, tool)
	case "kubeapps":
		return NewKubeapps(c, tool)
	}
	return nil, nil
}

func MergeValueMap(source map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	var valueStrings []string
	for k, v := range source {
		str := fmt.Sprintf("%s=%v", k, v)
		valueStrings = append(valueStrings, str)
	}
	for _, str := range valueStrings {
		err := strvals.ParseInto(str, result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func preInstallChart(h helm.Interface, tool *model.ClusterTool) error {
	rs, err := h.List()
	if err != nil {
		return err
	}
	for _, r := range rs {
		if r.Name == tool.Name {
			_, err := h.Uninstall(tool.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func installChart(h helm.Interface, tool *model.ClusterTool, chartName, chartVersion string) error {
	err := preInstallChart(h, tool)
	if err != nil {
		return err
	}
	valueMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(tool.Vars), &valueMap); err != nil {
		return err
	}
	m, err := MergeValueMap(valueMap)
	if err != nil {
		return err
	}
	_, err = h.Install(tool.Name, chartName, chartVersion, m)
	if err != nil {
		return err
	}
	return nil
}

func upgradeChart(h helm.Interface, tool *model.ClusterTool, chartName, chartVersion string) error {
	valueMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(tool.Vars), &valueMap); err != nil {
		return err
	}
	m, err := MergeValueMap(valueMap)
	if err != nil {
		return err
	}
	_, err = h.Upgrade(tool.Name, chartName, chartVersion, m)
	if err != nil {
		return err
	}
	return nil
}

func preCreateRoute(namespace string, ingressName string, kubeClient *kubernetes.Clientset) error {
	ingress, _ := kubeClient.NetworkingV1beta1().
		Ingresses(namespace).
		Get(context.TODO(), ingressName, metav1.GetOptions{})
	if ingress.Name != "" {
		err := kubeClient.NetworkingV1beta1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func createRoute(namespace string, ingressName string, ingressUrl string, serviceName string, port int, kubeClient *kubernetes.Clientset) error {
	if err := preCreateRoute(namespace, ingressName, kubeClient); err != nil {
		return err
	}
	service, err := kubeClient.CoreV1().
		Services(namespace).
		Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	ingress := v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: namespace,
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: ingressUrl,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Backend: v1beta1.IngressBackend{
										ServiceName: service.Name,
										ServicePort: intstr.FromInt(port),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = kubeClient.NetworkingV1beta1().Ingresses(namespace).Create(context.TODO(), &ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func waitForRunning(namespace string, deploymentName string, minReplicas int32, kubeClient *kubernetes.Clientset) error {
	kubeClient.CoreV1()
	err := wait.Poll(5*time.Second, 10*time.Minute, func() (done bool, err error) {
		d, err := kubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		if d.Status.ReadyReplicas > minReplicas-1 {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func waitForStatefulSetsRunning(namespace string, statefulSetsName string, minReplicas int32, kubeClient *kubernetes.Clientset) error {
	kubeClient.CoreV1()
	err := wait.Poll(5*time.Second, 10*time.Minute, func() (done bool, err error) {
		d, err := kubeClient.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetsName, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		if d.Status.ReadyReplicas > minReplicas-1 {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func uninstall(namespace string, tool *model.ClusterTool, ingressName string, h helm.Interface, kubeClient *kubernetes.Clientset) error {
	rs, err := h.List()
	if err != nil {
		return err
	}
	for _, r := range rs {
		if r.Name == tool.Name {
			_, _ = h.Uninstall(tool.Name)
		}
	}
	_ = kubeClient.NetworkingV1beta1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	return nil
}

func getGrafanaSourceNs(cluster model.Cluster, sourceFrom string) (string, error) {
	var sourceData model.ClusterTool
	if err := db.DB.
		Where("cluster_id = ? AND status = ? AND name = ?", cluster.ID, "Running", sourceFrom).
		Find(&sourceData).Error; err != nil {
		return "", err
	}
	sourceVars := map[string]interface{}{}
	if err := json.Unmarshal([]byte(sourceData.Vars), &sourceVars); err != nil {
		return "", err
	}
	sp, ok := sourceVars["namespace"]
	if !ok {
		return "", fmt.Errorf("获取prometheus ns 失败")
	}
	return sp.(string), nil
}
