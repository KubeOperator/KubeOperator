package service

import (
	"context"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"helm.sh/helm/v3/pkg/time"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

type KubernetesService interface {
	Get(req dto.SourceSearch) (interface{}, error)
	GetMetric(cluster string) (interface{}, error)
	CreateSc(req dto.SourceScCreate) error
	CreateSecret(req dto.SourceSecretCreate) error
	CordonNode(req dto.Cordon) error
	EvictPod(req dto.Evict) error
	Delete(req dto.SourceDelete) error
}

type kubernetesService struct {
	clusterService ClusterService
	clusterRepo    repository.ClusterRepository
}

func NewKubernetesService() KubernetesService {
	return &kubernetesService{
		clusterService: NewClusterService(),
		clusterRepo:    repository.NewClusterRepository(),
	}
}

func (k kubernetesService) Get(req dto.SourceSearch) (interface{}, error) {
	cluster, err := k.clusterRepo.GetWithPreload(req.Cluster, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return "", err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return "", err
	}

	switch req.Kind {
	case "namespacelist":
		ns, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		return ns, err
	case "podlist":
		pods, err := client.CoreV1().Pods(req.Namespace).List(context.TODO(), metav1.ListOptions{})
		return pods, err
	case "deploymentlist":
		pods, err := client.AppsV1().Deployments(req.Namespace).List(context.TODO(), metav1.ListOptions{})
		return pods, err
	case "nodelist":
		nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		return nodes, err
	case "secret":
		secret, err := client.CoreV1().Secrets(req.Namespace).Get(context.TODO(), req.Name, metav1.GetOptions{})
		return secret, err
	case "storageclasslist":
		if req.Limit != 0 {
			storageclass, err := client.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{Limit: req.Limit, Continue: req.Continue})
			return storageclass, err
		} else {
			storageclass, err := client.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
			return storageclass, err
		}
	case "overviewdatas":
		var overDatas OverViewData
		deployments, err := client.AppsV1().Deployments(req.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return overDatas, err
		}
		overDatas.Deployments = len(deployments.Items)

		pods, err := client.CoreV1().Pods(req.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return overDatas, err
		}
		overDatas.Pods = len(pods.Items)
		for _, po := range pods.Items {
			overDatas.Containers += len(po.Spec.Containers)
		}

		ns, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return overDatas, err
		}
		overDatas.Namespaces = len(ns.Items)

		nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return overDatas, err
		}
		overDatas.Nodes = len(nodes.Items)

		return overDatas, err
	}

	return dto.SourceList{}, nil
}

type OverViewData struct {
	Deployments int `json:"deployments"`
	Nodes       int `json:"nodes"`
	Namespaces  int `json:"namespaces"`
	Pods        int `json:"pods"`
	Containers  int `json:"containers"`
}

func (k kubernetesService) GetMetric(clusterName string) (interface{}, error) {
	cluster, err := k.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return "", err
	}
	host, err := clusterUtil.LoadAvailableHost(&cluster)
	if err != nil {
		return "", err
	}
	config, err := clusterUtil.LoadConnConf(&cluster, host)
	if err != nil {
		return "", err
	}
	mclient, err := metricsclientset.NewForConfig(config)
	if err != nil {
		return "", err
	}

	ms, err := mclient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	return ms, err
}

func (k kubernetesService) CreateSc(req dto.SourceScCreate) error {
	cluster, err := k.clusterRepo.GetWithPreload(req.Cluster, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return err
	}
	_, err = client.StorageV1().StorageClasses().Create(context.TODO(), &req.Info, metav1.CreateOptions{})
	return err
}

func (k kubernetesService) CordonNode(req dto.Cordon) error {
	cluster, err := k.clusterRepo.GetWithPreload(req.Cluster, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return err
	}

	node, err := client.CoreV1().Nodes().Get(context.TODO(), req.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	node.Spec.Unschedulable = req.SetUnschedulable
	_, err = client.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	return err
}

func (k kubernetesService) EvictPod(req dto.Evict) error {
	cluster, err := k.clusterRepo.GetWithPreload(req.Cluster, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return err
	}

	rmPod := &policyv1beta1.Eviction{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "policy/v1beta1",
			Kind:       "Eviction",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              req.Name,
			Namespace:         req.Namespace,
			CreationTimestamp: metav1.Time(time.Now()),
		},
	}
	if err := client.CoreV1().Pods(req.Namespace).EvictV1beta1(context.TODO(), rmPod); err != nil {
		return err
	}
	return err
}

func (k kubernetesService) CreateSecret(req dto.SourceSecretCreate) error {
	cluster, err := k.clusterRepo.GetWithPreload(req.Cluster, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return err
	}
	_, err = client.CoreV1().Secrets(req.Namespace).Create(context.TODO(), &req.Info, metav1.CreateOptions{})
	return err
}

func (k kubernetesService) Delete(req dto.SourceDelete) error {
	cluster, err := k.clusterRepo.GetWithPreload(req.Cluster, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	client, err := clusterUtil.NewClusterClient(&cluster)
	if err != nil {
		return err
	}

	switch req.Kind {
	case "storageclass":
		err := client.StorageV1().StorageClasses().Delete(context.TODO(), req.Name, metav1.DeleteOptions{})
		return err
	case "pod":
		err := client.CoreV1().Pods(req.Namespace).Delete(context.TODO(), req.Name, metav1.DeleteOptions{})
		return err
	case "secret":
		err := client.CoreV1().Secrets(req.Namespace).Delete(context.TODO(), req.Name, metav1.DeleteOptions{})
		return err
	}
	return nil
}
