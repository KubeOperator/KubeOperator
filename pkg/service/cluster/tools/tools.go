package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"time"
)

type Interface interface {
	Install() error
	Uninstall() error
}

type Cluster struct {
	model.Cluster
	HelmClient helm.Interface
	KubeClient *kubernetes.Clientset
}

func NewCluster(cluster model.Cluster, endpoint dto.Endpoint, secret model.ClusterSecret) (*Cluster, error) {
	c := Cluster{
		Cluster: cluster,
	}

	helmClient, err := helm.NewClient(helm.Config{
		ApiServer:   fmt.Sprintf("https://%s:%d", endpoint.Address, endpoint.Port),
		BearerToken: secret.KubernetesToken,
		Namespace:   constant.DefaultNamespace,
	})
	if err != nil {
		return nil, err
	}
	c.HelmClient = helmClient
	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Host:  endpoint.Address,
		Token: secret.KubernetesToken,
		Port:  endpoint.Port,
	})
	if err != nil {
		return nil, err
	}
	c.KubeClient = kubeClient
	return &c, nil
}

func NewClusterTool(tool *model.ClusterTool, cluster model.Cluster, endpoint dto.Endpoint, secret model.ClusterSecret) (Interface, error) {
	systemRepo := repository.NewSystemSettingRepository()
	localIP, err := systemRepo.Get("ip")
	if err != nil || localIP.Value == "" {
		return nil, errors.New("invalid system setting: ip")
	}
	c, err := NewCluster(cluster, endpoint, secret)
	if err != nil {
		return nil, err
	}
	switch tool.Name {
	case "prometheus":
		return NewPrometheus(c, localIP.Value, tool)
	case "efk":
		return NewEFK(c, localIP.Value, tool)
	case "registry":
		return NewRegistry(c, localIP.Value, tool)
	case "dashboard":
		return NewDashboard(c, localIP.Value, tool)
	case "chartmuseum":
		return NewChartmuseum(c, localIP.Value, tool)
	case "kubeapps":
		return NewKubeapps(c, localIP.Value, tool)
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

func installChart(h helm.Interface, tool *model.ClusterTool, chartName string) error {
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
	_, err = h.Install(tool.Name, chartName, m)
	if err != nil {
		return err
	}
	return nil
}

func preCreateRoute(ingressName string, kubeClient *kubernetes.Clientset) error {
	ingress, _ := kubeClient.NetworkingV1beta1().
		Ingresses(constant.DefaultNamespace).
		Get(context.TODO(), ingressName, metav1.GetOptions{})
	if ingress.Name != "" {
		err := kubeClient.NetworkingV1beta1().Ingresses(constant.DefaultNamespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func createRoute(ingressName string, ingressUrl string, serviceName string, port int, kubeClient *kubernetes.Clientset) error {
	if err := preCreateRoute(ingressName, kubeClient); err != nil {
		return err
	}
	service, err := kubeClient.CoreV1().
		Services(constant.DefaultNamespace).
		Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	ingress := v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: constant.DefaultNamespace,
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
	_, err = kubeClient.NetworkingV1beta1().Ingresses(constant.DefaultNamespace).Create(context.TODO(), &ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func waitForRunning(deploymentName string, minReplicas int32, kubeClient *kubernetes.Clientset) error {
	kubeClient.CoreV1()
	err := wait.Poll(5*time.Second, 30*time.Minute, func() (done bool, err error) {
		d, err := kubeClient.AppsV1().Deployments(constant.DefaultNamespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
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

func uninstall(tool *model.ClusterTool, ingressName string, h helm.Interface, kubeClient *kubernetes.Clientset) error {
	rs, err := h.List()
	if err != nil {
		return err
	}
	for _, r := range rs {
		if r.Name == tool.Name {
			_, _ = h.Uninstall(tool.Name)
		}
	}
	_ = kubeClient.NetworkingV1beta1().Ingresses(constant.DefaultNamespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	return nil
}
