package kubernetes

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	Host  string
	Token string
	Port  int
}

func NewKubernetesClient(c *Config) (*kubernetes.Clientset, error) {
	var clientSet kubernetes.Clientset
	kubeConf := &rest.Config{
		Host:        fmt.Sprintf("%s:%d", c.Host, c.Port),
		BearerToken: c.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	api, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		return &clientSet, err
	}
	return api, err
}
