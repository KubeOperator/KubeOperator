package prepare

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	playbookNameContainerRuntime = "02-runtime.yml"
)

type ContainerRuntimePhase struct {
	Upgrade bool
}

func (s ContainerRuntimePhase) Name() string {
	return "Install Container Runtime"
}

func (s ContainerRuntimePhase) Run(b kobe.Interface, fileName string) error {
	var tag string
	if s.Upgrade {
		tag = "upgrade"
	}

	return phases.RunPlaybookAndGetResult(b, playbookNameContainerRuntime, tag, fileName)
}
