package job

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"sync"
	"time"
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
	if len(repos) > 0 {
		log.Infof("scan job to sync...")
	}
	wg := &sync.WaitGroup{}
	for _, repo := range repos {
		interval := (time.Now().UnixNano() - repo.LastSyncTime.UnixNano()) / 1e6
		if repo.SyncStatus == constant.StatusPending && repo.SyncEnable && interval > repo.SyncInterval*time.Minute.Milliseconds() {
			log.Infof("repository %s need to sync", repo.Name)
			relations, err := m.multiClusterRepositoryService.GetClusterRelations(repo.Name)
			if err != nil {
				log.Error(err)
				return
			}
			if !(len(relations) > 0) {
				log.Info("repository not have related cluster. skip it")
				return
			}
			clusterNames := func() []string {
				var result []string
				for _, r := range relations {
					result = append(result, r.ClusterName)
				}
				return result
			}()
			s, err := service.NewMultiClusterRepositorySync(&repo.MultiClusterRepository, clusterNames)
			if err != nil {
				log.Error(err)
				return
			}
			go func() {
				wg.Add(1)
				s.Sync()
				log.Infof("repository %s sync completed", repo.Name)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
