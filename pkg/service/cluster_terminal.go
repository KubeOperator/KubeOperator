package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

type ClusterTerminalService interface {
	Terminal(cluster model.Cluster)
}

func NewCLusterTerminalService() ClusterTerminalService {
	return clusterTerminalService{
		clusterRepo:       repository.NewClusterRepository(),
		clusterStatusRepo: repository.NewClusterStatusRepository(),
	}
}

type clusterTerminalService struct {
	clusterRepo       repository.ClusterRepository
	clusterStatusRepo repository.ClusterStatusRepository
}

func (c clusterTerminalService) Terminal(cluster model.Cluster) {
	cluster.Status.Phase = constant.ClusterTerminating
	_ = c.clusterStatusRepo.Save(&cluster.Status)
	doTerminal(cluster)
	_ = c.clusterRepo.Delete(cluster.Name)
}

const terminalPlaybookName = "99-reset-cluster.yml"

func doTerminal(cluster model.Cluster) {
	inventory := cluster.ParseInventory()
	k := kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})
	for name, _ := range facts.DefaultFacts {
		k.SetVar(name, facts.DefaultFacts[name])
	}
	vars := cluster.GetKobeVars()
	for key, value := range vars {
		k.SetVar(key, value)
	}
	_ = phases.RunPlaybookAndGetResult(k, terminalPlaybookName)
}
