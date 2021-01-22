package job

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"sync"
)

var log = logger.Default

type RefreshHostInfo struct {
	hostService service.HostService
}

func NewRefreshHostInfo() *RefreshHostInfo {
	return &RefreshHostInfo{
		hostService: service.NewHostService(),
	}
}

func (r *RefreshHostInfo) Run() {
	var hosts []model.Host
	var wg sync.WaitGroup
	sem := make(chan struct{}, 2) // 信号量
	db.DB.Model(&model.Host{}).Find(&hosts)
	for _, host := range hosts {
		if host.Status == constant.ClusterCreating || host.Status == constant.Initializing {
			continue
		}
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			log.Infof("gather host [%s] info", name)
			_, err := r.hostService.Sync(name)
			if err != nil {
				log.Errorf("gather host info error: %s", err.Error())
			}
		}(host.Name)
	}
	wg.Wait()
}
