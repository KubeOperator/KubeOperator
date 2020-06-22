package service

import (
	"context"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strings"
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
	monitor, err := c.ClusterService.GetMonitor(clusterName)
	if err != nil {
		return err
	}
	endpoint, err := c.ClusterService.GetEndpoint(clusterName)
	if err != nil {
		return err
	}
	secret, err := c.ClusterService.GetSecrets(clusterName)
	if err != nil {
		return err
	}
	spec, err := c.ClusterService.GetSpec(clusterName)
	if err != nil {
		return err
	}
	m := monitor.ClusterMonitor
	m.Domain = fmt.Sprintf("prometheus.%s", spec.AppDomain)
	m.Status = constant.ClusterInitializing
	if err := c.ClusterMonitorRepo.Save(&m); err != nil {
		return err
	}
	c.Do(endpoint, secret, &m)
	return nil
}

func (c clusterMonitorService) Do(endpoint string, secret dto.ClusterSecret, monitor *model.ClusterMonitor) {
	helmClient, err := helm.NewClient(helm.Config{
		ApiServer:   fmt.Sprintf("https://%s:%d", endpoint, 8443),
		BearerToken: secret.KubernetesToken,
	})
	if err != nil {
		c.errorHandler(err, monitor)
		return
	}
	if err := installMonitor(helmClient); err != nil {
		c.errorHandler(err, monitor)
		return
	}
	k8sClient, err := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
		Host:  endpoint,
		Token: secret.KubernetesToken,
		Port:  8443,
	})
	if err := createMonitorIngress(k8sClient, monitor.Domain); err != nil {
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
	chart, err := helm.LoadCharts("../../resource/charts/prometheus-11.6.0.tgz")
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
