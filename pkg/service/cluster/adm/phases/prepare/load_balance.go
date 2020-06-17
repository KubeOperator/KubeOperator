package prepare

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	prepareLoadBalancer = "04-load-balancer.yml"
)

type LoadBalancerPhase struct {
	LbKubeApiserverIp string
}

func (s LoadBalancerPhase) Name() string {
	return "Install Load Balancer"
}

func (s LoadBalancerPhase) Run(b kobe.Interface) error {
	if s.LbKubeApiserverIp != "" {
		b.SetVar(facts.LbKubeApiserverPortFactName, s.LbKubeApiserverIp)
	}
	return phases.RunPlaybookAndGetResult(b, prepareLoadBalancer)
}
