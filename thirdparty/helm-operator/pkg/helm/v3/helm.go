package v3

import (
	"fmt"

	"github.com/go-kit/kit/log"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/cmd/util"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

const VERSION = "v3"

var (
	repositoryConfig = helmpath.ConfigPath("repositories.yaml")
	repositoryCache  = helmpath.CachePath("repository")
	pluginsDir       = helmpath.DataPath("plugins")
)

type HelmOptions struct {
	Driver    string
	Namespace string
}

type HelmV3 struct {
	kubeConfig *rest.Config
	logger     log.Logger
}

type infoLogFunc func(string, ...interface{})

// New creates a new HelmV3 client
func New(logger log.Logger, kubeConfig *rest.Config) helm.Client {
	// Add CRDs to the scheme. They are missing by default but required
	// by Helm v3.
	if err := apiextv1beta1.AddToScheme(scheme.Scheme); err != nil {
		// This should never happen.
		panic(err)
	}
	return &HelmV3{
		kubeConfig: kubeConfig,
		logger:     logger,
	}
}

func (h *HelmV3) Version() string {
	return VERSION
}

// infoLogFunc allows us to pass our logger to components
// that expect a klog.Infof function.
func (h *HelmV3) infoLogFunc(namespace string, releaseName string) infoLogFunc {
	return func(format string, args ...interface{}) {
		message := fmt.Sprintf(format, args...)
		h.logger.Log("info", message, "targetNamespace", namespace, "release", releaseName)
	}
}

func newActionConfig(config *rest.Config, logFunc infoLogFunc, namespace, driver string) (*action.Configuration, error) {

	restClientGetter := newConfigFlags(config, namespace)
	kubeClient := &kube.Client{
		Factory: util.NewFactory(restClientGetter),
		Log:     logFunc,
	}
	client, err := kubeClient.Factory.KubernetesClientSet()
	if err != nil {
		return nil, err
	}

	store, err := newStorageDriver(client, logFunc, namespace, driver)
	if err != nil {
		return nil, err
	}

	return &action.Configuration{
		RESTClientGetter: restClientGetter,
		Releases:         store,
		KubeClient:       kubeClient,
		Log:              logFunc,
	}, nil
}

func newConfigFlags(config *rest.Config, namespace string) *genericclioptions.ConfigFlags {
	return &genericclioptions.ConfigFlags{
		Namespace:   &namespace,
		APIServer:   &config.Host,
		CAFile:      &config.CAFile,
		BearerToken: &config.BearerToken,
	}
}

func newStorageDriver(client *kubernetes.Clientset, logFunc infoLogFunc, namespace, d string) (*storage.Storage, error) {
	switch d {
	case "secret", "secrets", "":
		s := driver.NewSecrets(client.CoreV1().Secrets(namespace))
		s.Log = logFunc
		return storage.Init(s), nil
	case "configmap", "configmaps":
		c := driver.NewConfigMaps(client.CoreV1().ConfigMaps(namespace))
		c.Log = logFunc
		return storage.Init(c), nil
	case "memory":
		m := driver.NewMemory()
		return storage.Init(m), nil
	default:
		return nil, fmt.Errorf("unsupported storage driver '%s'", d)
	}
}

func getterProviders() getter.Providers {
	return getter.All(&cli.EnvSettings{
		RepositoryConfig: repositoryConfig,
		RepositoryCache:  repositoryCache,
		PluginsDirectory: pluginsDir,
	})
}
