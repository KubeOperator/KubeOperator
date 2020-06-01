package init

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initMaster = "08-master.yml"
)

type MasterPhase struct {
	KubePodSubnet         string
	KubeServiceSubnet     string
	KubeNetworkNodePrefix string
	KubeMaxPods           string
	KubeProxyMode         string
	NetworkPlugin         string
}

func (MasterPhase) Name() string {
	return "InitEtcd"
}

func (s MasterPhase) Run(b kobe.Interface) (result kobe.Result, err error) {
	if s.KubePodSubnet != "" {
		b.SetVar(facts.KubePodSubnetFactName, s.KubePodSubnet)
	}

	if s.KubeServiceSubnet != "" {
		b.SetVar(facts.KubeServiceSubnetFactName, s.KubeServiceSubnet)
	}
	if s.KubeNetworkNodePrefix != "" {
		b.SetVar(facts.KubeNetworkNodePrefixFactName, s.KubeNetworkNodePrefix)
	}
	if s.KubeMaxPods != "" {
		b.SetVar(facts.KubeMaxPodsFactName, s.KubeMaxPods)
	}
	if s.KubeProxyMode != "" {
		b.SetVar(facts.KubeProxyModeFactName, s.KubeProxyMode)
	}
	if s.NetworkPlugin != "" {
		b.SetVar(facts.NetworkPluginFactName, s.NetworkPlugin)
	}
	return phases.RunPlaybookAndGetResult(b, initMaster)
}
