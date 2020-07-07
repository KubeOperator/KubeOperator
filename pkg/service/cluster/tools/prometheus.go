package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/grafana"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type Prometheus struct {
	Cluster       model.Cluster
	HelmClient    helm.Interface
	KubeClient    *kubernetes.Clientset
	GrafanaClient grafana.Interface
	Tool          *model.ClusterTool
}

func NewPrometheus(cluster dto.ClusterWithEndpoint, tool *model.ClusterTool) (*Prometheus, error) {
	p := &Prometheus{
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

func (p Prometheus) Install() error {
	if err := p.installChart(); err != nil {
		return err
	}

	if err := p.createRoute(); err != nil {
		return err
	}

	if err := p.createGrafanaDataSource(); err != nil {
		return err
	}
	url, err := p.createGrafanaDashboard()
	if err != nil {
		return err
	}
	p.Tool.Url = url
	return nil
}

func (p Prometheus) installChart() error {
	valueMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(p.Tool.Vars), &valueMap)
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
	_, err := p.HelmClient.Install(p.Tool.Name, constant.PrometheusChartName, map[string]interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (p Prometheus) createRoute() error {
	services, err := p.KubeClient.CoreV1().Services(constant.DefaultNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: fmt.Sprintf("release=%s", p.Tool.Name)})
	if err != nil {
		return err
	}
	serviceName := ""
	for _, svc := range services.Items {
		if strings.Contains(svc.Name, "prometheus-server") {
			serviceName = svc.Name
		}
	}
	ingress := v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prometheus-ingress",
			Namespace: constant.DefaultNamespace,
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: constant.DefaultPrometheusIngress,
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

func (p Prometheus) createGrafanaDataSource() error {
	url := fmt.Sprintf("http://localhost:8080/proxy/prometheus/%s/", p.Cluster.Name)
	return p.GrafanaClient.CreateDataSource(p.Cluster.Name, url)

}
func (p Prometheus) createGrafanaDashboard() (string, error) {
	return p.GrafanaClient.CreateDashboard(p.Cluster.Name)
}
