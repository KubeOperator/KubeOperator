package kubernetes

import (
	"fmt"
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/net"
	extensionClientSet "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Host string

var log = logger.Default

func NewKubernetesClient(conf *string) (*kubernetes.Clientset, error) {
	apiConfig, err := pauseConfigApi(conf)
	if err != nil {
		return nil, err
	}

	getter := func() (*api.Config, error) {
		return apiConfig, nil
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", getter)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func NewMetricClient(conf *string) (*metricsclientset.Clientset, error) {
	apiConfig, err := pauseConfigApi(conf)
	if err != nil {
		return nil, err
	}

	getter := func() (*api.Config, error) {
		return apiConfig, nil
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", getter)
	if err != nil {
		return nil, err
	}
	return metricsclientset.NewForConfig(config)
}

func NewKubernetesExtensionClient(conf *string) (*extensionClientSet.Clientset, error) {
	apiConfig, err := pauseConfigApi(conf)
	if err != nil {
		return nil, err
	}

	getter := func() (*api.Config, error) {
		return apiConfig, nil
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", getter)
	if err != nil {
		return nil, err
	}
	return extensionClientSet.NewForConfig(config)
}

func pauseConfigApi(conf *string) (*api.Config, error) {
	config, err := clientcmd.Load([]byte(*conf))
	if err != nil {
		return nil, err
	}
	for key, obj := range config.AuthInfos {
		config.AuthInfos[key] = obj
	}
	for key, obj := range config.Clusters {
		config.Clusters[key] = obj
	}
	for key, obj := range config.Contexts {
		config.Contexts[key] = obj
	}

	if config.AuthInfos == nil {
		config.AuthInfos = map[string]*api.AuthInfo{}
	}
	if config.Clusters == nil {
		config.Clusters = map[string]*api.Cluster{}
	}
	if config.Contexts == nil {
		config.Contexts = map[string]*api.Context{}
	}
	return config, nil
}

func SelectAliveHost(hosts []Host) (Host, error) {
	var aliveHost Host
	aliveHostCh := make(chan Host, len(hosts)+1)
	wg := &sync.WaitGroup{}
	for i := range hosts {
		wg.Add(1)
		go func(h Host) {
			defer wg.Done()
			err := net.TcpPing(string(h))
			if err != nil {
				log.Warnf("dial host %s error %s", h, err.Error())
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
		return "", fmt.Errorf("no alive host in %v", hosts)
	}
	return aliveHost, nil
}
