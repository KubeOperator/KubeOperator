package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"io"
)

const (
	initHelm = "11-helm-install.yml"
)

type HelmPhase struct {
}

func (h HelmPhase) Name() string {
	return "InitHelm"
}

func (h HelmPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, initHelm, "", writer)
}
