package kubernetes

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Config struct {
	Host  string
	Token string
	Port  int
}

func NewKubernetesClient(c *Config) (*kubernetes.Clientset, error) {
	kubeConf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	kubeConf.Host = address
	kubeConf.BearerToken = c.Token
	api, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		return nil, err
	}
	return api, err
}
