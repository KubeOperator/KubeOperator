package v3

import (
	"os"
	"sync"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/repo"
)

var repositoryConfigLock sync.RWMutex

func (h *HelmV3) RepositoryIndex() error {

	repositoryConfigLock.RLock()
	f, err := loadRepositoryConfig()
	repositoryConfigLock.RUnlock()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, c := range f.Repositories {
		r, err := newChartRepository(c)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(r *repo.ChartRepository) {
			if _, err := r.DownloadIndexFile(); err != nil {
				h.logger.Log("error", "unable to get an update from the chart repository", "url", r.Config.URL, "err", err)
			}
			wg.Done()
		}(r)
	}
	wg.Wait()
	return nil
}

func (h *HelmV3) RepositoryAdd(name, url, username, password, certFile, keyFile, caFile string) error {
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
		return errors.New("chart repository with name '%s' already exists")
	}

	r, err := newChartRepository(c)
	if err != nil {
		return err
	}
	if _, err = r.DownloadIndexFile(); err != nil {
		return err
	}

	return f.WriteFile(repositoryConfig, 0644)
}

func (h *HelmV3) RepositoryRemove(name string) error {
	repositoryConfigLock.Lock()
	defer repositoryConfigLock.Unlock()

	f, err := repo.LoadFile(repositoryConfig)
	if err != nil {
		return err
	}
	f.Remove(name)

	return f.WriteFile(repositoryConfig, 0644)
}

func (h *HelmV3) RepositoryImport(path string) error {
	s, err := repo.LoadFile(path)
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
		r, err := newChartRepository(c)
		if err != nil {
			h.logger.Log("error", err, "name", c.Name, "url", c.URL)
			continue
		}
		if _, err := r.DownloadIndexFile(); err != nil {
			h.logger.Log("error", err, "name", c.Name, "url", c.URL)
			continue
		}

		t.Add(c)
		h.logger.Log("info", "successfully imported repository", "name", c.Name, "url", c.URL)
	}

	return t.WriteFile(repositoryConfig, 0644)
}

// newChartRepository constructs a new `repo.ChartRepository`
// for the given `repo.Entry`. It exists to stay in control
// of the cache path and getters while duplicating as less
// code as possible.
func newChartRepository(e *repo.Entry) (*repo.ChartRepository, error) {
	cr, err := repo.NewChartRepository(e, getterProviders())
	if err != nil {
		return nil, err
	}
	cr.CachePath = repositoryCache
	return cr, err
}

func loadRepositoryConfig() (*repo.File, error) {
	r, err := repo.LoadFile(repositoryConfig)
	if err != nil && !os.IsNotExist(errors.Cause(err)) {
		return nil, err
	}
	return r, nil
}
