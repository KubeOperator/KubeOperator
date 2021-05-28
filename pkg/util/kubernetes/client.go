package kubernetes

import (
	"fmt"
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/net"
	"github.com/pkg/errors"
	extensionClientSet "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Host string

type Config struct {
	Hosts []Host
	Token string
}

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
	client, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		return client, errors.Wrap(err, fmt.Sprintf("new kubernetes client with config failed: %v", err))
	}
	return client, nil
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
	client, err := extensionClientSet.NewForConfig(kubeConf)
	if err != nil {
		return client, errors.Wrap(err, fmt.Sprintf("new extension kubernetes client with config failed: %v", err))
	}
	return client, nil
}
func SelectAliveHost(hosts []Host) (Host, error) {
	var aliveHost Host
	aliveHostCh := make(chan Host, len(hosts)+1)
	wg := &sync.WaitGroup{}
	for i := range hosts {
		wg.Add(1)
		go func(h Host) {
			defer wg.Done()
			err := net.TcpPing(string(h), true)
			if err != nil {
				logger.Log.Warnf("dial host %s falied: %+v", h, err)
				return
			}
			aliveHostCh <- h
		}(hosts[i])
	}
	go func() {
		wg.Wait()
		aliveHostCh <- ""
	}()
	aliveHost = <-aliveHostCh
	if aliveHost == "" {
		return "", errors.Wrap(errors.New("no alive host"), "")
	}
	return aliveHost, nil
}
