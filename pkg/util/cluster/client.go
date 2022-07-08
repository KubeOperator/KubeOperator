package cluster

import (
	"fmt"
	"sync"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/net"
	"github.com/pkg/errors"
	extensionClientSet "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func NewClusterClient(cluster *model.Cluster) (*kubernetes.Clientset, error) {
	var hosts []string
	port := cluster.SpecConf.KubeApiServerPort
	hosts = append(hosts, fmt.Sprintf("%s:%d", cluster.SpecConf.LbKubeApiserverIp, port))
	for _, node := range cluster.Nodes {
		if node.Role == constant.NodeRoleNameMaster {
			hosts = append(hosts, fmt.Sprintf("%s:%d", node.Host.Ip, port))
		}
	}
	availableHost, err := SelectAliveHost(hosts)
	if err != nil {
		return nil, err
	}
	conf, err := LoadConnConf(cluster, availableHost)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return client, err
	}
	return client, nil
}

func LoadConnConf(cluster *model.Cluster, availableHost string) (*rest.Config, error) {
	var connConf rest.Config
	switch cluster.SpecConf.AuthenticationMode {
	case constant.AuthenticationModeBearer:
		connConf.Host = availableHost
		connConf.BearerToken = cluster.Secret.KubernetesToken
	case constant.AuthenticationModeCertificate:
		connConf.CertData = []byte(cluster.Secret.CertDataStr)
		connConf.KeyData = []byte(cluster.Secret.KeyDataStr)
	case constant.AuthenticationModeConfigFile:
		apiConfig, err := PauseConfigApi(&cluster.Secret.ConfigContent)
		if err != nil {
			return nil, err
		}
		getter := func() (*api.Config, error) {
			return apiConfig, nil
		}
		itemConfig, err := clientcmd.BuildConfigFromKubeconfigGetter("", getter)
		if err != nil {
			return nil, err
		}
		connConf = *itemConfig
	}
	return &connConf, nil
}

func LoadAvailableHost(cluster *model.Cluster) (string, error) {
	var hosts []string
	port := cluster.SpecConf.KubeApiServerPort
	hosts = append(hosts, fmt.Sprintf("%s:%d", cluster.SpecConf.LbKubeApiserverIp, port))
	for _, node := range cluster.Nodes {
		if node.Role == constant.NodeRoleNameMaster {
			hosts = append(hosts, fmt.Sprintf("%s:%d", node.Host.Ip, port))
		}
	}
	return SelectAliveHost(hosts)
}

func NewClusterExtensionClient(cluster *model.Cluster) (*extensionClientSet.Clientset, error) {
	var hosts []string
	port := cluster.SpecConf.KubeApiServerPort
	hosts = append(hosts, fmt.Sprintf("%s:%d", cluster.SpecConf.LbKubeApiserverIp, port))
	for _, node := range cluster.Nodes {
		if node.Role == constant.NodeRoleNameMaster {
			hosts = append(hosts, fmt.Sprintf("%s:%d", node.Host.Ip, port))
		}
	}
	availableHost, err := SelectAliveHost(hosts)
	if err != nil {
		return nil, err
	}
	conf, err := LoadConnConf(cluster, availableHost)
	if err != nil {
		return nil, err
	}

	client, err := extensionClientSet.NewForConfig(conf)
	if err != nil {
		return client, err
	}
	return client, nil
}

func SelectAliveHost(hosts []string) (string, error) {
	var aliveHost string
	aliveHostCh := make(chan string, len(hosts)+1)
	wg := &sync.WaitGroup{}
	for i := range hosts {
		wg.Add(1)
		go func(h string) {
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
		return "", errors.New("NO_MASTER_AVAILABLE")
	}
	return aliveHost, nil
}

func PauseConfigApi(conf *string) (*api.Config, error) {
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
