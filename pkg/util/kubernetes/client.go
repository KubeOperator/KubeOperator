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
	var clientSet kubernetes.Clientset
	kubeConf, err := config.GetConfig()
	if err != nil {
		return &clientSet, err
	}
	kubeConf.Host = fmt.Sprintf("%s:%d", c.Host, c.Port)
	kubeConf.BearerToken = c.Token
	api, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		return &clientSet, err
	}
	return api, err
}
