package kube

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Interface interface {
	getPodsByNamespace(namespace string) ([]v1.Pod, error)
	getAllPods() ([]v1.Pod, error)
	getNameSpaces() ([]v1.Namespace, error)
	GetEvents() ([]v1.Event, error)
}

type Kube struct {
	config *KubeConfig
	client kubernetes.Clientset
}

func NewKube(k *KubeConfig) (*Kube, error) {
	client, err := k.InitK8SClient()
	if err != nil {
		return nil, err
	}
	return &Kube{
		config: k,
		client: *client,
	}, nil
}

func (k *Kube) getNameSpaces() ([]v1.Namespace, error) {
	namespaceList, err := k.client.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return namespaceList.Items, err
}

func (k *Kube) getPodsByNamespace(namespace string) ([]v1.Pod, error) {
	podList, err := k.client.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList.Items, err
}

func (k *Kube) GetEvents() ([]v1.Event, error) {
	var events []v1.Event
	namespaceList, err := k.client.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return events, err
	}
	for _, namespace := range namespaceList.Items {
		eventList, err := k.client.CoreV1().Events(namespace.Name).List(context.TODO(), metaV1.ListOptions{})
		if err != nil {
			break
		}
		events = append(events, eventList.Items...)
	}
	return events, err
}
