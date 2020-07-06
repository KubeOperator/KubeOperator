package tools

import (
	"context"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/grafana"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strings"
)

const elasticsearch = "elasticsearch.tar.gz"

type Logging struct {
	Cluster       model.Cluster
	HelmClient    helm.Interface
	KubeClient    *kubernetes.Clientset
	GrafanaClient grafana.Interface
	Tool          *model.ClusterTool
}

func NewLogging(cluster dto.ClusterWithEndpoint, tool *model.ClusterTool) (*Logging, error) {
	p := &Logging{
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
	p.GrafanaClient = grafana.NewClient()
	return p, nil
}

func (p Logging) Install() error {
	if err := p.installChart(); err != nil {
		return err
	}

	if err := p.createRoute(); err != nil {
		return err
	}

	return nil
}

func (p Logging) installChart() error {
	chart, err := helm.LoadCharts(elasticsearch)
	if err != nil {
		return err
	}
	_, err = p.HelmClient.Install("elasticsearch", chart, map[string]interface{}{})
	return nil
}

func (p Logging) createRoute() error {
	services, err := p.KubeClient.CoreV1().Services(constant.DefaultNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: fmt.Sprintf("release=%s", p.Tool.Name)})
	if err != nil {
		return err
	}
	serviceName := ""
	for _, svc := range services.Items {
		if strings.Contains(svc.Name, "elasticsearch-logging") {
			serviceName = svc.Name
		}
	}

	ingress := v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "elasticsearch-logging-ingress",
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
	_, err = p.KubeClient.NetworkingV1beta1().Ingresses(constant.DefaultNamespace).Create(context.TODO(), &ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
