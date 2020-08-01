package helm

import (
	"context"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/ghodss/yaml"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"io/ioutil"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	helmDriver = "configmap"
)

var log = logger.Default

func nolog(format string, v ...interface{}) {}

type Interface interface {
	Install(name string, chartName string, values map[string]interface{}) (*release.Release, error)
	Uninstall(name string) (*release.UninstallReleaseResponse, error)
	List() ([]*release.Release, error)
}

type Config struct {
	ApiServer   string
	BearerToken string
	Namespace   string
}
type Client struct {
	actionConfig *action.Configuration
	Namespace    string
	settings     *cli.EnvSettings
}

func NewClient(config Config) (*Client, error) {
	client := Client{}
	client.settings = GetSettings()
	cf := genericclioptions.NewConfigFlags(true)
	inscure := true
	cf.APIServer = &config.ApiServer
	cf.BearerToken = &config.BearerToken
	cf.Insecure = &inscure
	if config.Namespace == "" {
		client.Namespace = constant.DefaultNamespace
	} else {
		client.Namespace = config.Namespace
	}
	cf.Namespace = &client.Namespace
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(cf, client.Namespace, helmDriver, nolog)
	if err != nil {
		return nil, err
	}
	client.actionConfig = actionConfig
	return &client, nil
}

func LoadCharts(path string) (*chart.Chart, error) {
	return loader.Load(path)
}

func (c Client) Install(name string, chartName string, values map[string]interface{}) (*release.Release, error) {
	if err := updateRepo(); err != nil {
		return nil, err
	}
	client := action.NewInstall(c.actionConfig)
	client.ReleaseName = name
	client.Namespace = c.Namespace
	p, err := client.ChartPathOptions.LocateChart(chartName, c.settings)
	if err != nil {
		return nil, err
	}
	ct, err := loader.Load(p)
	if err != nil {
		return nil, err
	}

	return client.Run(ct, values)
}
func (c Client) Uninstall(name string) (*release.UninstallReleaseResponse, error) {
	client := action.NewUninstall(c.actionConfig)
	return client.Run(name)
}

func (c Client) List() ([]*release.Release, error) {
	client := action.NewList(c.actionConfig)
	client.All = true
	return client.Run()
}

func GetSettings() *cli.EnvSettings {
	return &cli.EnvSettings{
		PluginsDirectory: helmpath.DataPath("plugins"),
		RegistryConfig:   helmpath.ConfigPath("registry.json"),
		RepositoryConfig: helmpath.ConfigPath("repositories.yaml"),
		RepositoryCache:  helmpath.CachePath("repository"),
	}

}

func updateRepo() error {
	repos, _ := ListRepo()
	flag := false
	for _, r := range repos {
		if r.Name == "nexus" {
			flag = true
		}
	}
	if !flag {
		r := repository.NewSystemSettingRepository()
		s, err := r.Get("ip")
		if err != nil && s.Value == "" {
			return errors.New("can not find local hostname")
		}
		err = addRepo("nexus", fmt.Sprintf("http://%s:8081/repository/applications-amd64", s.Value), "admin", "admin123")
		if err != nil {
			return err
		}
		err = addRepo("nexus", fmt.Sprintf("http://%s:8081/repository/applications-arm64", s.Value), "admin", "admin123")
		if err != nil {
			return err
		}
	} else {
		settings := GetSettings()
		repoFile := settings.RepositoryConfig
		repoCache := settings.RepositoryCache
		f, err := repo.LoadFile(repoFile)
		if err != nil {
			return err
		}
		var repos []*repo.ChartRepository
		for _, cfg := range f.Repositories {
			r, err := repo.NewChartRepository(cfg, getter.All(settings))
			if err != nil {
				return err
			}
			if repoCache != "" {
				r.CachePath = repoCache
			}
			repos = append(repos, r)
		}
		updateCharts(repos)
	}
	return nil
}

func updateCharts(repos []*repo.ChartRepository) {
	log.Debug("Hang tight while we grab the latest from your chart repositories...")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				log.Debugf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				log.Debugf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	log.Debugf("Update Complete. ⎈ Happy Helming!⎈ ")
}

func addRepo(name string, url string, username string, password string) error {
	settings := GetSettings()

	repoFile := settings.RepositoryConfig

	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if f.Has(name) {
		return errors.Errorf("repository name (%s) already exists, please specify a different name", name)
	}

	e := repo.Entry{
		Name:                  name,
		URL:                   url,
		Username:              username,
		Password:              password,
		InsecureSkipTLSverify: true,
	}

	r, err := repo.NewChartRepository(&e, getter.All(settings))
	if err != nil {
		return err
	}
	r.CachePath = settings.RepositoryCache
	if _, err := r.DownloadIndexFile(); err != nil {
		return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
	}

	f.Update(&e)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		return err
	}
	return nil
}

func ListRepo() ([]*repo.Entry, error) {
	settings := GetSettings()
	var repos []*repo.Entry
	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return repos, err
	}
	return f.Repositories, nil
}
