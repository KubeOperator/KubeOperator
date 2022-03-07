package prepare

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	prepareKubernetesComponents = "03-kubernetes-component.yml"
)

type KubernetesComponentPhase struct {
}

func (s KubernetesComponentPhase) Name() string {
	return "Prepare Kubernetes Component"
}

func (s KubernetesComponentPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, prepareKubernetesComponents, "", fileName)
}
