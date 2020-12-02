package v2

import (
	"os"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/repo"
)

var repositoryConfigLock sync.RWMutex

func (h *HelmV2) RepositoryIndex() error {

	repositoryConfigLock.RLock()
	f, err := loadRepositoryConfig()
	repositoryConfigLock.RUnlock()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, c := range f.Repositories {
		r, err := repo.NewChartRepository(c, getterProviders())
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(r *repo.ChartRepository) {
			if err := r.DownloadIndexFile(repositoryCache); err != nil {
				h.logger.Log("error", "unable to get an update from the chart repository", "url", r.Config.URL, "err", err)
			}
			wg.Done()
		}(r)
	}
	wg.Wait()
	return nil
}

func (h *HelmV2) RepositoryAdd(name, url, username, password, certFile, keyFile, caFile string) error {
	repositoryConfigLock.Lock()
	defer repositoryConfigLock.Unlock()

	f, err := loadRepositoryConfig()
	if err != nil {
		return err
	}

	c := &repo.Entry{
		Name:     name,
		URL:      url,
		Username: username,
		Password: password,
		CertFile: certFile,
		KeyFile:  keyFile,
		CAFile:   caFile,
	}
	f.Add(c)

	if f.Has(name) {
		return errors.New("chart repository with name %s already exists")
	}

	r, err := repo.NewChartRepository(c, getterProviders())
	if err != nil {
		return err
	}
	if err = r.DownloadIndexFile(repositoryCache); err != nil {
		return err
	}

	return f.WriteFile(repositoryConfig, 0644)
}

func (h *HelmV2) RepositoryRemove(name string) error {
	repositoryConfigLock.Lock()
	defer repositoryConfigLock.Unlock()

	f, err := repo.LoadRepositoriesFile(repositoryConfig)
	if err != nil {
		return err
	}
	f.Remove(name)

	return f.WriteFile(repositoryConfig, 0644)
}

func (h *HelmV2) RepositoryImport(path string) error {
	s, err := repo.LoadRepositoriesFile(path)
	if err != nil {
		return err
	}

	repositoryConfigLock.Lock()
	defer repositoryConfigLock.Unlock()

	t, err := loadRepositoryConfig()
	if err != nil {
		return err
	}

	for _, c := range s.Repositories {
		if t.Has(c.Name) {
			h.logger.Log("error", "repository with name already exists", "name", c.Name, "url", c.URL)
			continue
		}
		r, err := repo.NewChartRepository(c, getterProviders())
		if err != nil {
			h.logger.Log("error", err, "name", c.Name, "url", c.URL)
			continue
		}
		if err := r.DownloadIndexFile(repositoryCache); err != nil {
			h.logger.Log("error", err, "name", c.Name, "url", c.URL)
			continue
		}

		t.Add(c)
		h.logger.Log("info", "successfully imported repository", "name", c.Name, "url", c.URL)
	}

	return t.WriteFile(repositoryConfig, 0644)
}

func loadRepositoryConfig() (*repo.RepoFile, error) {
	r, err := repo.LoadRepositoriesFile(repositoryConfig)
	if err != nil && !os.IsNotExist(errors.Cause(err)) {
		return nil, err
	}
	return r, nil
}
