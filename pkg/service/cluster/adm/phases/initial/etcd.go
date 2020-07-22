package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initEtcd = "06-etcd.yml"
)

type EtcdPhase struct {
}

func (s EtcdPhase) Name() string {
	return "InitEtcd"
}

func (s EtcdPhase) Run(b kobe.Interface) error {
	return phases.RunPlaybookAndGetResult(b, initEtcd)
}
