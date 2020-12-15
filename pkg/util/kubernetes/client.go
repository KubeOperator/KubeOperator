package kubernetes

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/net"
	extensionClientSet "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
)

type Host string

type Config struct {
	Hosts []Host
	Token string
}

var log = logger.Default

func NewKubernetesClient(c *Config) (*kubernetes.Clientset, error) {
	var aliveHost Host
	aliveHost, err := SelectAliveHost(c.Hosts)
	if err != nil {
		return nil, err
	}
	kubeConf := &rest.Config{
		Host:        string(aliveHost),
		BearerToken: c.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	return kubernetes.NewForConfig(kubeConf)
}

func NewKubernetesExtensionClient(c *Config) (*extensionClientSet.Clientset, error) {
	var aliveHost Host
	aliveHost, err := SelectAliveHost(c.Hosts)
	if err != nil {
		return nil, err
	}
	kubeConf := &rest.Config{
		Host:        string(aliveHost),
		BearerToken: c.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	return extensionClientSet.NewForConfig(kubeConf)
}
func SelectAliveHost(hosts []Host) (Host, error) {
	var aliveHost Host
	aliveHostCh := make(chan Host,len(hosts)+1)
	wg := &sync.WaitGroup{}
	for i := range hosts {
		wg.Add(1)
		go func(h Host) {
			defer wg.Done()
			err := net.TcpPing(string(h), true)
			if err != nil {
				log.Warnf("dial host %s error %s",h, err.Error())
				return
			}
			aliveHostCh <-h
		}(hosts[i])
	}
	go func() {
		wg.Wait()
		aliveHostCh <- ""
	}()
	aliveHost=<-aliveHostCh
	if aliveHost==""{
		return "", fmt.Errorf("no alive host in %v", hosts)
	}
	return aliveHost, nil
}
