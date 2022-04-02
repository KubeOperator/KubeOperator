package service

import (
	"context"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubernetesService interface {
	Get(req dto.SourceSearch) (interface{}, error)
	GetMetric(cluster string) (interface{}, error)
	CreateNs(req dto.SourceNsCreate) error
	CreateSc(req dto.SourceScCreate) error
	CreatePv(req dto.SourcePvCreate) error
	CreateSecret(req dto.SourceSecretCreate) error
	Delete(req dto.SourceDelete) error
}

type kubernetesService struct {
	clusterService ClusterService
}

func NewKubernetesService() KubernetesService {
	return &kubernetesService{
		clusterService: NewClusterService(),
	}
}

func (k kubernetesService) Get(req dto.SourceSearch) (interface{}, error) {
	secret, err := k.clusterService.GetSecrets(req.Cluster)
	if err != nil {
		return dto.SourceList{}, err
	}
	client, err := kubeUtil.NewKubernetesClient(&secret.KubeConf)
	if err != nil {
		return dto.SourceList{}, err
	}

	switch req.Kind {
	case "namespacelist":
		pods, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		return pods, err
	case "podlist":
		pods, err := client.CoreV1().Pods(req.Namespace).List(context.TODO(), metav1.ListOptions{})
		return pods, err
	case "nodelist":
		nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		return nodes, err
	case "eventlist":
		if req.Limit != 0 {
			events, err := client.CoreV1().Events(req.Namespace).List(context.TODO(), metav1.ListOptions{Limit: req.Limit, Continue: req.Continue})
			return events, err
		} else {
			events, err := client.CoreV1().Events(req.Namespace).List(context.TODO(), metav1.ListOptions{})
			return events, err
		}
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
	case "pvlist":
		if req.Limit != 0 {
			storageclass, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{Limit: req.Limit, Continue: req.Continue})
			return storageclass, err
		} else {
			storageclass, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
			return storageclass, err
		}
	case "deploymentlist":
		deployments, err := client.AppsV1().Deployments(req.Namespace).List(context.TODO(), metav1.ListOptions{})
		return deployments, err
	}

	return dto.SourceList{}, nil
}

func (k kubernetesService) GetMetric(cluster string) (interface{}, error) {
	secret, err := k.clusterService.GetSecrets(cluster)
	if err != nil {
		return dto.SourceList{}, err
	}
	mclient, err := kubeUtil.NewMetricClient(&secret.KubeConf)
	if err != nil {
		return dto.SourceList{}, err
	}

	ms, err := mclient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	return ms, err
}

func (k kubernetesService) CreateNs(req dto.SourceNsCreate) error {
	secret, err := k.clusterService.GetSecrets(req.Cluster)
	if err != nil {
		return err
	}
	client, err := kubeUtil.NewKubernetesClient(&secret.KubeConf)
	if err != nil {
		return err
	}
	_, err = client.CoreV1().Namespaces().Create(context.TODO(), &req.Info, metav1.CreateOptions{})
	return err
}

func (k kubernetesService) CreateSc(req dto.SourceScCreate) error {
	secret, err := k.clusterService.GetSecrets(req.Cluster)
	if err != nil {
		return err
	}
	client, err := kubeUtil.NewKubernetesClient(&secret.KubeConf)
	if err != nil {
		return err
	}
	_, err = client.StorageV1().StorageClasses().Create(context.TODO(), &req.Info, metav1.CreateOptions{})
	return err
}

func (k kubernetesService) CreatePv(req dto.SourcePvCreate) error {
	secret, err := k.clusterService.GetSecrets(req.Cluster)
	if err != nil {
		return err
	}
	client, err := kubeUtil.NewKubernetesClient(&secret.KubeConf)
	if err != nil {
		return err
	}
	_, err = client.CoreV1().PersistentVolumes().Create(context.TODO(), &req.Info, metav1.CreateOptions{})
	return err
}

func (k kubernetesService) CreateSecret(req dto.SourceSecretCreate) error {
	secret, err := k.clusterService.GetSecrets(req.Cluster)
	if err != nil {
		return err
	}
	client, err := kubeUtil.NewKubernetesClient(&secret.KubeConf)
	if err != nil {
		return err
	}
	_, err = client.CoreV1().Secrets(req.Namespace).Create(context.TODO(), &req.Info, metav1.CreateOptions{})
	return err
}

func (k kubernetesService) Delete(req dto.SourceDelete) error {
	secret, err := k.clusterService.GetSecrets(req.Cluster)
	if err != nil {
		return err
	}
	client, err := kubeUtil.NewKubernetesClient(&secret.KubeConf)
	if err != nil {
		return err
	}

	switch req.Kind {
	case "pv":
		err := client.CoreV1().PersistentVolumes().Delete(context.TODO(), req.Name, metav1.DeleteOptions{})
		return err
	case "storageclass":
		err := client.StorageV1().StorageClasses().Delete(context.TODO(), req.Name, metav1.DeleteOptions{})
		return err
	case "secret":
		err := client.CoreV1().Secrets(req.Namespace).Delete(context.TODO(), req.Name, metav1.DeleteOptions{})
		return err
	case "namespace":
		err := client.CoreV1().Namespaces().Delete(context.TODO(), req.Name, metav1.DeleteOptions{})
		return err
	}
	return nil
}
