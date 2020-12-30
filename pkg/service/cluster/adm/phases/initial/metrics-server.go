package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"io"
)

const (
	initMetricsServer = "13-metrics-server.yml"
)

type MetricsServerPhase struct {
}

func (m MetricsServerPhase) Name() string {
	return "Npd Init"
}

func (m MetricsServerPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, initMetricsServer, "", writer)
}
