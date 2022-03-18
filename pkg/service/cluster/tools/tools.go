package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/pkg/errors"

	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"helm.sh/helm/v3/pkg/strvals"

	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"

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
	helmRepoPort int
	HelmClient   helm.Interface
	KubeClient   *kubernetes.Clientset
}

type Ingress struct {
	name    string
	url     string
	service string
	port    int
	version string
}

func NewCluster(cluster model.Cluster, hosts []kubernetesUtil.Host, oldNamespace, namespace string) (*Cluster, error) {
	c := Cluster{
		Cluster: cluster,
	}
	var registery model.SystemRegistry
	if cluster.Spec.Architectures == constant.ArchAMD64 {
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfAMD64).First(&registery).Error; err != nil {
			return nil, errors.New("load image pull port failed")
		}
	} else {
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfARM64).First(&registery).Error; err != nil {
			return nil, errors.New("load image pull port failed")
		}
	}
	c.helmRepoPort = registery.RegistryPort
	c.Namespace = namespace
	helmClient, err := helm.NewClient(&helm.Config{
		Hosts:         hosts,
		BearerToken:   cluster.Secret.KubernetesToken,
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
		Token: cluster.Secret.KubernetesToken,
	})
	if err != nil {
		return nil, err
	}
	c.KubeClient = kubeClient
	return &c, nil
}

func NewClusterTool(tool *model.ClusterTool, cluster model.Cluster, hosts []kubernetesUtil.Host, oldNamespace, namespace string, enable bool) (Interface, error) {
	c, err := NewCluster(cluster, hosts, oldNamespace, namespace)
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
			lokiNs, _ := getGrafanaSourceNs(cluster, "loki")
			return NewGrafana(c, tool, prometheusNs, lokiNs)
		} else {
			return NewGrafana(c, tool, "", "")
		}
	case "registry":
		return NewRegistry(c, tool)
	case "dashboard":
		return NewDashboard(c, tool)
	case "gatekeeper":
		return NewGatekeeper(c, tool)
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
			logger.Log.Infof("uninstall %s before installation", tool.Name)
			_, err := h.Uninstall(tool.Name)
			if err != nil {
				return err
			}
		}
	}
	logger.Log.Infof("uninstall %s before installation successful", tool.Name)
	return nil
}

func installChart(h helm.Interface, tool *model.ClusterTool, chartName, chartVersion string) error {
	err := preInstallChart(h, tool)
	if err != nil {
		return err
	}
	valueMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(tool.Vars), &valueMap)
	m, err := MergeValueMap(valueMap)
	if err != nil {
		return err
	}
	logger.Log.Infof("start install tool %s with chartName: %s, chartVersion: %s", tool.Name, chartName, chartVersion)
	_, err = h.Install(tool.Name, chartName, chartVersion, m)
	if err != nil {
		return err
	}
	logger.Log.Infof("install tool %s successful", tool.Name)
	return nil
}

func upgradeChart(h helm.Interface, tool *model.ClusterTool, chartName, chartVersion string) error {
	valueMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(tool.Vars), &valueMap); err != nil {
		return err
	}
	m, err := MergeValueMap(valueMap)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("merge value map failed: %v", err))
	}
	logger.Log.Infof("start upgrade tool %s with chartName: %s, chartVersion: %s", tool.Name, chartName, chartVersion)
	_, err = h.Upgrade(tool.Name, chartName, chartVersion, m)
	if err != nil {
		return err
	}
	logger.Log.Infof("upgrade tool %s successful", tool.Name)
	return nil
}

func preCreateRoute(namespace string, ingressName string, version string, kubeClient *kubernetes.Clientset) error {
	if isApiV1(version) {
		ingress, _ := kubeClient.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
		if ingress.Name != "" {
			if err := kubeClient.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
	} else {
		ingress, _ := kubeClient.NetworkingV1beta1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
		if ingress.Name != "" {
			if err := kubeClient.NetworkingV1beta1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
	}
	logger.Log.Infof("operation before create route %s successful", ingressName)
	return nil
}

func createRoute(namespace string, ingressInfo *Ingress, kubeClient *kubernetes.Clientset) error {
	if err := preCreateRoute(namespace, ingressInfo.name, ingressInfo.version, kubeClient); err != nil {
		return err
	}
	service, err := kubeClient.CoreV1().
		Services(namespace).
		Get(context.TODO(), ingressInfo.service, metav1.GetOptions{})
	if err != nil {
		return err
	}

	ingressInfo.service = service.Name
	if isApiV1(ingressInfo.version) {
		ingress := newNetworkV1(namespace, ingressInfo)
		if _, err = kubeClient.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{}); err != nil {
			return err
		}
	} else {
		ingress := newNetworkV1bate1(namespace, ingressInfo)
		if _, err = kubeClient.NetworkingV1beta1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{}); err != nil {
			return err
		}
	}
	logger.Log.Infof("create route %s successful", ingressInfo.name)
	return nil
}

func newNetworkV1(namespace string, ingressInfo *Ingress) *netv1.Ingress {
	pathType := netv1.PathTypePrefix
	ingress := netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressInfo.name,
			Namespace: namespace,
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					Host: ingressInfo.url,
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: ingressInfo.service,
											Port: netv1.ServiceBackendPort{
												Number: int32(ingressInfo.port),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return &ingress
}

func newNetworkV1bate1(namespace string, ingressInfo *Ingress) *netv1beta1.Ingress {
	ingress := netv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressInfo.name,
			Namespace: namespace,
		},
		Spec: netv1beta1.IngressSpec{
			Rules: []netv1beta1.IngressRule{
				{
					Host: ingressInfo.url,
					IngressRuleValue: netv1beta1.IngressRuleValue{
						HTTP: &netv1beta1.HTTPIngressRuleValue{
							Paths: []netv1beta1.HTTPIngressPath{
								{
									Backend: netv1beta1.IngressBackend{
										ServiceName: ingressInfo.service,
										ServicePort: intstr.FromInt(int(ingressInfo.port)),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return &ingress
}

func isApiV1(version1 string) bool {
	version1 = strings.ReplaceAll(version1, "v", "")

	verTag1 := strings.Split(version1, ".")
	if len(verTag1) < 3 {
		return false
	}
	itemVersion, _ := strconv.Atoi(verTag1[1])
	return itemVersion > 18
}

func waitForRunning(namespace string, deploymentName string, minReplicas int32, kubeClient *kubernetes.Clientset) error {
	logger.Log.Infof("installation and configuration successful, now waiting for %s running", deploymentName)
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
	logger.Log.Infof("installation and configuration successful, now waiting for %s running", statefulSetsName)
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

func uninstall(namespace string, tool *model.ClusterTool, ingressName string, version string, h helm.Interface, kubeClient *kubernetes.Clientset) error {
	rs, err := h.List()
	if err != nil {
		return err
	}
	for _, r := range rs {
		if r.Name == tool.Name {
			_, _ = h.Uninstall(tool.Name)
		}
	}

	if isApiV1(version) {
		if err := kubeClient.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{}); err != nil {
			logger.Log.Errorf("uninstall tool %s of namespace %s failed, err: %v", tool.Name, namespace, err)
		}
	} else {
		if err := kubeClient.NetworkingV1beta1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{}); err != nil {
			logger.Log.Errorf("uninstall tool %s of namespace %s failed, err: %v", tool.Name, namespace, err)
		}
	}

	logger.Log.Infof("uninstall tool %s of namespace %s successful", tool.Name, namespace)
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
	_ = json.Unmarshal([]byte(sourceData.Vars), &sourceVars)
	sp, ok := sourceVars["namespace"]
	if !ok {
		return "", fmt.Errorf("load namespace of prometheus failed")
	}
	return sp.(string), nil
}
