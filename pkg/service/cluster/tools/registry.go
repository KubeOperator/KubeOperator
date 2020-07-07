package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type Registry struct {
	Cluster    model.Cluster
	HelmClient helm.Interface
	KubeClient *kubernetes.Clientset
	Tool       *model.ClusterTool
}

func NewRegistry(cluster dto.ClusterWithEndpoint, tool *model.ClusterTool) (*Registry, error) {
	p := &Registry{
		Tool: tool,
	}
	p.Cluster = cluster.Cluster
	helmClient, err := helm.NewClient(helm.Config{
		ApiServer:   fmt.Sprintf("https://%s:%d", cluster.Endpoint.Address, cluster.Endpoint.Port),
		BearerToken: cluster.Cluster.Secret.KubernetesToken,
		Namespace:   constant.DefaultNamespace,
	})
	if err != nil {
		return nil, err
	}
	p.HelmClient = helmClient
	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Host:  cluster.Endpoint.Address,
		Token: cluster.Cluster.Secret.KubernetesToken,
		Port:  cluster.Endpoint.Port,
	})
	if err != nil {
		return nil, err
	}
	p.KubeClient = kubeClient
	return p, nil
}

func (c Registry) Install() error {
	if err := c.installChart(); err != nil {
		return err
	}

	if err := c.createRoute(); err != nil {
		return err
	}

	return nil
}
func (c Registry) installChart() error {
	valueMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &valueMap)
	var valueStrings []string
	for k, v := range valueMap {
		str := fmt.Sprintf("%s=%v", k, v)
		valueStrings = append(valueStrings, str)
	}
	valueMap = map[string]interface{}{}
	for _, str := range valueStrings {
		err := strvals.ParseIntoString(str, valueMap)
		if err != nil {
			return err
		}
	}
	_, err := c.HelmClient.Install(c.Tool.Name, constant.DockerRegistryChartName, valueMap)
	if err != nil {
		return nil
	}
	return nil
}

func (c Registry) createRoute() error {
	services, err := c.KubeClient.CoreV1().Services(constant.DefaultNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: fmt.Sprintf("release=%s", c.Tool.Name)})
	if err != nil {
		return err
	}
	serviceName := ""
	for _, svc := range services.Items {
		if strings.Contains(svc.Name, "docker-registry") {
			serviceName = svc.Name
		}
	}
	ingress := v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "docker-registry-ingress",
			Namespace: constant.DefaultNamespace,
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: constant.DefaultLoggingIngress,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Backend: v1beta1.IngressBackend{
										ServiceName: serviceName,
										ServicePort: intstr.FromInt(80),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = c.KubeClient.NetworkingV1beta1().Ingresses(constant.DefaultNamespace).Create(context.TODO(), &ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
