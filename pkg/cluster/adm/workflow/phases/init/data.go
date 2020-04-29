package phases

import (
	"ko3-gin/pkg/cluster/adm/api"
)

type InitData interface {
	Cfg() *api.InitConfiguration
}
