package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initEtcd = "06-etcd.yml"
)

type EtcdPhase struct {
	Upgrade bool
}

func (s EtcdPhase) Name() string {
	return "InitEtcd"
}

func (s EtcdPhase) Run(b kobe.Interface, fileName string) error {
	var tag string
	if s.Upgrade {
		tag = "upgrade"
	}
	return phases.RunPlaybookAndGetResult(b, initEtcd, tag, fileName)
}
