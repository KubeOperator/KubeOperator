package phase

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	playbookNameContainerd = "02-containerd.yml"
)

type ContainerdRuntimePhase struct {
}

func (s ContainerdRuntimePhase) Name() string {
	return "Install Containerd"
}

func (s ContainerdRuntimePhase) Run(b kobe.Interface) (result kobe.Result, err error) {
	return phases.RunPlaybookAndGetResult(b, playbookNameContainerd)
}
