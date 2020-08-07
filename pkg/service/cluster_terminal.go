package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"sync"
)

type ClusterTerminalService interface {
	Terminal(cluster model.Cluster)
}

func NewCLusterTerminalService() ClusterTerminalService {
	return clusterTerminalService{
		clusterRepo:       repository.NewClusterRepository(),
		clusterStatusRepo: repository.NewClusterStatusRepository(),
		planRepo:          repository.NewPlanRepository(),
	}
}

type clusterTerminalService struct {
	clusterRepo       repository.ClusterRepository
	clusterStatusRepo repository.ClusterStatusRepository
	planRepo          repository.PlanRepository
}

func (c clusterTerminalService) Terminal(cluster model.Cluster) {
	if cluster.Status.Phase == constant.ClusterTerminating {
		return
	}
	cluster.Status.Phase = constant.ClusterTerminating
	_ = c.clusterStatusRepo.Save(&cluster.Status)

	var waitGroup sync.WaitGroup
	switch cluster.Spec.Provider {
	case constant.ClusterProviderBareMetal:
		go doBareMetalTerminal(&waitGroup, &cluster)
		waitGroup.Add(1)
	case constant.ClusterProviderPlan:
		go doPlanTerminal(&waitGroup, &cluster)
		waitGroup.Add(1)
	default:
		return
	}
	waitGroup.Wait()
	err := c.clusterRepo.Delete(cluster.Name)
	if err != nil {
		log.Error(err)
	}
}

const terminalPlaybookName = "99-reset-cluster.yml"

func doPlanTerminal(wg *sync.WaitGroup, cluster *model.Cluster) {
	defer wg.Done()
	//k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	//_, err := k.Destroy()
	//if err != nil {
	//	log.Error(err)
	//}
}

func doBareMetalTerminal(wg *sync.WaitGroup, cluster *model.Cluster) {
	defer wg.Done()
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
	err := phases.RunPlaybookAndGetResult(k, terminalPlaybookName)
	if err != nil {
		log.Error(err)
	}
}
