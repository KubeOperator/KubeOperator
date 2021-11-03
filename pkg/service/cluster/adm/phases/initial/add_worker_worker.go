package initial

import (
	"io"

	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initAddWorkerWorker = "91-add-worker-06-kubernetes-worker.yml"
)

type AddWorkerMasterPhase struct {
}

func (AddWorkerMasterPhase) Name() string {
	return "InitWorker"
}

func (s AddWorkerMasterPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, initAddWorkerWorker, "", writer)
}
