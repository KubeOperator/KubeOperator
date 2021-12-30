package prometheus

import (
	"fmt"
	"testing"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
)

func TestGetMetric(t *testing.T) {
	p, err := NewPrometheusService("http://172.16.10.113:32044")
	if err != nil {
		fmt.Println(err)
	}
	exs := []string{"cpu_used", "memory_used"}
	op := dto.QueryOptions{
		Level:         "cluster",
		NodeName:      "",
		Start:         1640756661,
		End:           1640763861,
		Step:          2,
		MetricsFilter: exs,
	}
	value := p.GetNamedMetricsOverTime(&op)

	fmt.Println(value)
}
