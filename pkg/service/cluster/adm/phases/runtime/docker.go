package phase

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	playbookNameDocker = "02-docker.yml"
)

type DockerRuntimePhase struct {
}

func (s DockerRuntimePhase) Name() string {
	return "Install Docker"
}

func (s DockerRuntimePhase) Run(b kobe.Interface) (result kobe.Result, err error) {
	return phases.RunPlaybookAndGetResult(b, playbookNameDocker)
}
