package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	defaultNamespace = "default"
	helmDriver       = "configmap"
)

func nolog(format string, v ...interface{}) {}

type Interface interface {
	Install(name string, chart *chart.Chart, values map[string]interface{})
	Uninstall(name string) (*release.UninstallReleaseResponse, error)
	List(name string) ([]*release.Release, error)
}

type Config struct {
	ApiServer   string
	BearerToken string
}
type Client struct {
	actionConfig *action.Configuration
}

func NewClient(config Config) (*Client, error) {
	client := Client{}
	cf := genericclioptions.NewConfigFlags(true)
	inscure := true
	cf.APIServer = &config.ApiServer
	cf.BearerToken = &config.BearerToken
	cf.Insecure = &inscure
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(cf, defaultNamespace, helmDriver, nolog)
	if err != nil {
		return nil, err
	}
	client.actionConfig = actionConfig
	return &client, nil
}

func LoadCharts(path string) (*chart.Chart, error) {
	return loader.Load(path)
}

func (c Client) Install(name string, chart *chart.Chart, values map[string]interface{}) (*release.Release, error) {
	client := action.NewInstall(c.actionConfig)
	client.ReleaseName = name
	return client.Run(chart, values)
}
func (c Client) Uninstall(name string) (*release.UninstallReleaseResponse, error) {
	client := action.NewUninstall(c.actionConfig)
	return client.Run(name)
}

func (c Client) List(name string) ([]*release.Release, error) {
	client := action.NewList(c.actionConfig)
	client.All = true
	return client.Run()
}
