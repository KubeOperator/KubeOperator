package prepare

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	playbookNameContainerRuntime = "02-runtime.yml"
)

type ContainerRuntimePhase struct {
	ContainerRuntime     string
	DockerStorageDir     string
	ContainerdStorageDir string
}

func (s ContainerRuntimePhase) Name() string {
	return "Install Container Runtime"
}

func (s ContainerRuntimePhase) Run(b kobe.Interface) error {
	if s.ContainerRuntime != "" {
		b.SetVar(facts.ContainerRuntimeFactName, s.ContainerRuntime)
	}
	if s.DockerStorageDir != "" {
		b.SetVar(facts.DockerStorageDirFactName, s.DockerStorageDir)
	}
	if s.ContainerdStorageDir != "" {
		b.SetVar(facts.ContainerdStorageDirFactName, s.ContainerdStorageDir)
	}
	return phases.RunPlaybookAndGetResult(b, playbookNameContainerRuntime)
}
