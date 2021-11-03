package prepare

import (
	"io"

	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	prepareAddWorkerKubernetesComponents = "91-add-worker-03-kubernetes-component.yml"
)

type AddWorkerKubernetesComponentPhase struct {
}

func (s AddWorkerKubernetesComponentPhase) Name() string {
	return "Prepare Kubernetes Component"
}

func (s AddWorkerKubernetesComponentPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, prepareAddWorkerKubernetesComponents, "", writer)
}
