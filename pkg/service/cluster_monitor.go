package service

import (
	"context"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/KubeOperator/KubeOperator/pkg/util/grafana"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

type ClusterMonitorService interface {
	Init(clusterName string) error
}

func NewClusterMonitorService() ClusterMonitorService {
	return &clusterMonitorService{
		ClusterService:     NewClusterService(),
		ClusterMonitorRepo: repository.NewClusterMonitorRepository(),
	}
}

type clusterMonitorService struct {
	ClusterMonitorRepo repository.ClusterMonitorRepository
	ClusterService     ClusterService
}

func (c clusterMonitorService) Init(clusterName string) error {
	cluster, err := c.ClusterService.Get(clusterName)
	if err != nil {
		return err
	}
	monitor, err := c.ClusterService.GetMonitor(clusterName)
	if err != nil {
		return err
	}
	endpoint, err := c.ClusterService.GetApiServerEndpoint(clusterName)
	if err != nil {
		return err
	}
	secret, err := c.ClusterService.GetSecrets(clusterName)
	if err != nil {
		return err
	}

	m := monitor.ClusterMonitor
	m.Domain = fmt.Sprintf("prometheus.%s", cluster.Spec.AppDomain)
	m.Status = constant.ClusterInitializing
	m.Enable = true
	if err := c.ClusterMonitorRepo.Save(&m); err != nil {
		return err
	}
	go c.Do(cluster.Cluster, endpoint, secret, &m)
	return nil
}

func (c clusterMonitorService) Do(cluster model.Cluster, endpoint dto.Endpoint, secret dto.ClusterSecret, monitor *model.ClusterMonitor) {
	helmClient, err := helm.NewClient(helm.Config{
		ApiServer:   fmt.Sprintf("https://%s:%d", endpoint.Address, endpoint.Port),
		BearerToken: secret.KubernetesToken,
	})
	time.Sleep(10 * time.Second)
	if err != nil {
		c.errorHandler(err, monitor)
		return
	}
	if err := installMonitor(helmClient); err != nil {
		c.errorHandler(err, monitor)
		return
	}
	k8sClient, err := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
		Host:  endpoint.Address,
		Port:  endpoint.Port,
		Token: secret.KubernetesToken,
	})
	if err := createMonitorIngress(k8sClient, monitor.Domain); err != nil {
		c.errorHandler(err, monitor)
		return
	}
	grafanaClient := grafana.NewClient()
	if err := createGrafanaDataSource(cluster.Name, grafanaClient); err != nil {
		c.errorHandler(err, monitor)
		return
	}
	url, err := createGrafanaDashboard(cluster.Name, grafanaClient)
	monitor.DashboardUrl = url
	if err != nil {
		c.errorHandler(err, monitor)
		return
	}
	monitor.Status = constant.ClusterRunning
	_ = c.ClusterMonitorRepo.Save(monitor)
}

func (c clusterMonitorService) errorHandler(err error, monitor *model.ClusterMonitor) {
	monitor.Status = constant.ClusterFailed
	monitor.Message = err.Error()
	_ = c.ClusterMonitorRepo.Save(monitor)
}

func installMonitor(helmClient helm.Interface) error {
	chart, err := helm.LoadCharts("resource/charts/prometheus-11.6.0.tgz")
	if err != nil {
		return err
	}
	values := map[string]interface{}{
		"alertmanager": map[string]interface{}{
			"enabled": false,
		},
		"server": map[string]interface{}{
			"persistentVolume": map[string]interface{}{
				"enabled": false,
			},
		},
		"pushgateway": map[string]interface{}{
			"enabled": false,
		},
	}
	_, err = helmClient.Install("monitor", chart, values)
	if err != nil {
		return err
	}
	return nil
}

func createMonitorIngress(client *kubernetes.Clientset, domain string) error {
	var prometheusServiceName string
	svcs, err := client.CoreV1().Services(constant.DefaultNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "release=monitor"})
	if err != nil {
		return err
	}
	for _, svc := range svcs.Items {
		if strings.Contains(svc.Name, "prometheus-server") {
			prometheusServiceName = svc.Name
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
					Host: domain,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Backend: v1beta1.IngressBackend{
										ServiceName: prometheusServiceName,
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
	_, err = client.NetworkingV1beta1().Ingresses(constant.DefaultNamespace).Create(context.TODO(), &ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func createGrafanaDataSource(clusterName string, grafanaClient grafana.Interface) error {
	url := fmt.Sprintf("http://localhost:8080/proxy/prometheus/%s/", clusterName)
	return grafanaClient.CreateDataSource(clusterName, url)

}
func createGrafanaDashboard(clusterName string, grafanaClient grafana.Interface) (string, error) {
	return grafanaClient.CreateDashboard(clusterName)
}
