package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initWorker = "08-kubernetes-worker.yml"
)

type WorkerPhase struct {
}

func (WorkerPhase) Name() string {
	return "InitWorker"
}

func (s WorkerPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, initWorker, "", fileName)
}
