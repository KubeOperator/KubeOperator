package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
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
		messageService:    NewMessageService(),
	}
}

type clusterTerminalService struct {
	clusterRepo       repository.ClusterRepository
	clusterStatusRepo repository.ClusterStatusRepository
	planRepo          repository.PlanRepository
	messageService    MessageService
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
		waitGroup.Add(1)
		go doBareMetalTerminal(&waitGroup, &cluster)
	case constant.ClusterProviderPlan:
		waitGroup.Add(1)
		go doPlanTerminal(&waitGroup, &cluster)
	default:
		return
	}
	waitGroup.Wait()
	_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterUnInstall, true, ""), cluster.Name, constant.ClusterUnInstall)
	err := c.clusterRepo.Delete(cluster.Name)
	if err != nil {
		log.Error(err)
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUnInstall, false, err.Error()), cluster.Name, constant.ClusterUnInstall)
	} else {
		log.Error(err)
	}
}

const terminalPlaybookName = "99-reset-cluster.yml"

func doPlanTerminal(wg *sync.WaitGroup, cluster *model.Cluster) {
	defer wg.Done()
	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	_, err := k.Destroy()
	if err != nil {
		log.Error(err)
		messageService := NewMessageService()
		_ = messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUnInstall, false, err.Error()), cluster.Name, constant.ClusterUnInstall)
	}
}

func doBareMetalTerminal(wg *sync.WaitGroup, cluster *model.Cluster) {
	defer wg.Done()
	inventory := cluster.ParseInventory()
	k := kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})
	for i := range facts.DefaultFacts {
		k.SetVar(i, facts.DefaultFacts[i])
	}
	vars := cluster.GetKobeVars()
	for key, value := range vars {
		k.SetVar(key, value)
	}
	err := phases.RunPlaybookAndGetResult(k, terminalPlaybookName, nil)
	if err != nil {
		log.Error(err)
		messageService := NewMessageService()
		_ = messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUnInstall, false, err.Error()), cluster.Name, constant.ClusterUnInstall)
	}
}
