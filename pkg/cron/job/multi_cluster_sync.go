package job

import (
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"sync"
)

type MultiClusterSyncJob struct {
	multiClusterRepositoryService service.MultiClusterRepositoryService
}

func NewMultiClusterSyncJob() *MultiClusterSyncJob {
	return &MultiClusterSyncJob{
		multiClusterRepositoryService: service.NewMultiClusterRepositoryService(),
	}
}

func (m *MultiClusterSyncJob) Run() {
	repos, err := m.multiClusterRepositoryService.List()
	if err != nil {
		log.Error(err)
		return
	}
	wg := &sync.WaitGroup{}
	for _, repo := range repos {
		if repo.SyncEnable {
			s, err := service.NewMultiClusterRepositorySync(&repo.MultiClusterRepository, []string{})
			if err != nil {
				log.Error(err)
				return
			}
			go func() {
				wg.Add(1)
				s.Sync()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
