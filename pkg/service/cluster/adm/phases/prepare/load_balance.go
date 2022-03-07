package prepare

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	prepareLoadBalancer = "04-load-balancer.yml"
)

type LoadBalancerPhase struct {
}

func (s LoadBalancerPhase) Name() string {
	return "Install Load Balancer"
}

func (s LoadBalancerPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, prepareLoadBalancer, "", fileName)
}
