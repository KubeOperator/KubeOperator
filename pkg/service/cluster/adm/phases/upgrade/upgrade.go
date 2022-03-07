package upgrade

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	upgradeCluster = "92-upgrade-cluster.yml"
)

type UpgradeClusterPhase struct {
	Version string
}

func (upgrade UpgradeClusterPhase) Name() string {
	return "upgradeCluster"
}

func (upgrade UpgradeClusterPhase) Run(b kobe.Interface, fileName string) error {
	if upgrade.Version != "" {
		b.SetVar("kube_upgrade_version", upgrade.Version)
	}
	return phases.RunPlaybookAndGetResult(b, upgradeCluster, "", fileName)
}
