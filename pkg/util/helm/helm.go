package helm

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/spf13/viper"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
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
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	helmDriver = "configmap"
)

var log = logger.Default

func nolog(format string, v ...interface{}) {}

type Interface interface {
	Install(name string, chartName string, chartVersion string, values map[string]interface{}) (*release.Release, error)
	Upgrade(name string, chartName string, chartVersion string, values map[string]interface{}) (*release.Release, error)
	Uninstall(name string) (*release.UninstallReleaseResponse, error)
	List() ([]*release.Release, error)
	GetRepoIP(arch string) (string, string, error)
}

type Config struct {
	Hosts         []kubernetes.Host
	BearerToken   string
	OldNamespace  string
	Namespace     string
	Architectures string
}
type Client struct {
	installActionConfig   *action.Configuration
	unInstallActionConfig *action.Configuration
	Namespace             string
	settings              *cli.EnvSettings
	Architectures         string
}

func NewClient(config *Config) (*Client, error) {
	var aliveHost kubernetes.Host
	aliveHost, err := kubernetes.SelectAliveHost(config.Hosts)
	if err != nil {
		return nil, err
	}
	client := Client{
		Architectures: config.Architectures,
	}
	client.settings = GetSettings()
	cf := genericclioptions.NewConfigFlags(true)
	inscure := true
	apiServer := fmt.Sprintf("https://%s", aliveHost)
	cf.APIServer = &apiServer
	cf.BearerToken = &config.BearerToken
	cf.Insecure = &inscure
	if config.Namespace == "" {
		client.Namespace = constant.DefaultNamespace
	} else {
		client.Namespace = config.Namespace
	}
	cf.Namespace = &client.Namespace
	installActionConfig := new(action.Configuration)
	if err := installActionConfig.Init(cf, client.Namespace, helmDriver, nolog); err != nil {
		return nil, err
	}
	client.installActionConfig = installActionConfig
	unInstallActionConfig := new(action.Configuration)
	if err := unInstallActionConfig.Init(cf, config.OldNamespace, helmDriver, nolog); err != nil {
		return nil, err
	}
	client.unInstallActionConfig = unInstallActionConfig
	return &client, nil
}

func LoadCharts(path string) (*chart.Chart, error) {
	return loader.Load(path)
}

func (c Client) Install(name, chartName, chartVersion string, values map[string]interface{}) (*release.Release, error) {
	if err := updateRepo(c.Architectures); err != nil {
		return nil, err
	}
	client := action.NewInstall(c.installActionConfig)
	client.ReleaseName = name
	client.Namespace = c.Namespace
	client.ChartPathOptions.InsecureSkipTLSverify = true
	if len(chartVersion) != 0 {
		client.ChartPathOptions.Version = chartVersion
	}
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

func (c Client) Upgrade(name, chartName, chartVersion string, values map[string]interface{}) (*release.Release, error) {
	if err := updateRepo(c.Architectures); err != nil {
		return nil, err
	}
	client := action.NewUpgrade(c.installActionConfig)
	client.Namespace = c.Namespace
	client.DryRun = false
	client.ChartPathOptions.InsecureSkipTLSverify = true
	client.ChartPathOptions.Version = chartVersion
	p, err := client.ChartPathOptions.LocateChart(chartName, c.settings)
	if err != nil {
		return nil, err
	}
	ct, err := loader.Load(p)
	if err != nil {
		return nil, err
	}

	return client.Run(name, ct, values)
}

func (c Client) Uninstall(name string) (*release.UninstallReleaseResponse, error) {
	client := action.NewUninstall(c.unInstallActionConfig)
	return client.Run(name)
}

func (c Client) List() ([]*release.Release, error) {
	client := action.NewList(c.unInstallActionConfig)
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

func updateRepo(arch string) error {
	repos, err := ListRepo()
	if err != nil {
		return err
	}
	flag := false
	for _, r := range repos {
		if r.Name == "nexus" {
			flag = true
		}
	}
	if !flag {
		var c Client
		repoProtocol, repoIP, err := c.GetRepoIP(arch)
		if err != nil {
			return err
		}
		password, err := encrypt.StringDecrypt(viper.GetString("repository.nexus_password"))
		if err != nil {
			return err
		}
		if err = addRepo("nexus", fmt.Sprintf("%s://%s::8081/repository/applications", repoProtocol, repoIP), "admin", password); err != nil {
			log.Errorf("addRepo failed, error: %s", err.Error())
			return err
		}
	}
	settings := GetSettings()
	repoFile := settings.RepositoryConfig
	repoCache := settings.RepositoryCache
	f, err := repo.LoadFile(repoFile)
	if err != nil {
		return err
	}
	var rps []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			return err
		}
		if repoCache != "" {
			r.CachePath = repoCache
		}
		rps = append(rps, r)
	}
	updateCharts(rps)
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
		defer func() {
			if err := fileLock.Unlock(); err != nil {
				log.Errorf("addRepo fileLock.Unlock failed, error: %s", err.Error())
			}
		}()
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

func (c Client) GetRepoIP(arch string) (string, string, error) {
	var repo model.SystemRegistry
	switch arch {
	case "amd64":
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfAMD64).First(&repo).Error; err != nil {
			return "", "", err
		}
		return repo.Protocol, repo.Hostname, nil
	case "arm64":
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfARM64).First(&repo).Error; err != nil {
			return "", "", err
		}
		return repo.Protocol, repo.Hostname, nil
	case "all":
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfARM64).First(&repo).Error; err != nil {
			return "", "", err
		}
		return repo.Protocol, repo.Hostname, nil
	}
	return "", "", fmt.Errorf("no such architecture")
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
