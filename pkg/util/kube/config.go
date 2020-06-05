package kube

import (
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"strconv"
)

type KubeConfig struct {
	host  string
	token string
	port  int
}

func NewKubeConfig(host string, token string, port int) *KubeConfig {
	return &KubeConfig{
		host:  host,
		token: token,
		port:  port,
	}
}

func (c *KubeConfig) InitK8SClient() (*kubernetes.Clientset, error) {
	var clientSet kubernetes.Clientset
	kubeConf, err := config.GetConfig()
	if err != nil {
		return &clientSet, err
	}
	kubeConf.Host = c.host + ":" + strconv.Itoa(c.port)
	kubeConf.BearerToken = c.token
	api, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		return &clientSet, err
	}
	return api, err
}
