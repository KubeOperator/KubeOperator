package job

import (
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/service"
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
}
