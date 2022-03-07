package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initMetricsServer = "13-metrics-server.yml"
)

type MetricsServerPhase struct {
}

func (m MetricsServerPhase) Name() string {
	return "Npd Init"
}

func (m MetricsServerPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, initMetricsServer, "", fileName)
}
