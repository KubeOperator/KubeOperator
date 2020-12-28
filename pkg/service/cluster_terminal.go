package service

import (
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
)

type ClusterTerminalService interface {
	Terminal(cluster model.Cluster) error
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

func (c clusterTerminalService) Terminal(cluster model.Cluster) error {
	// if cluster.Status.Phase == constant.ClusterTerminating {
	// 	return
	// }
	cluster.Status.Phase = constant.ClusterTerminating
	cluster.Status.ClusterStatusConditions = []model.ClusterStatusCondition{}
	condition := model.ClusterStatusCondition{
		Name:          "DeleteCluster",
		Status:        constant.ConditionUnknown,
		OrderNum:      0,
		LastProbeTime: time.Now(),
	}
	cluster.Status.ClusterStatusConditions = append(cluster.Status.ClusterStatusConditions, condition)
	if err := c.clusterStatusRepo.Save(&cluster.Status); err != nil {
		return err
	}

	if cluster.Spec.Provider == constant.ClusterProviderBareMetal {
		if err := doBareMetalTerminal(&cluster); err != nil {
			c.errClusterDelete(&cluster, "uninstall cluster err: "+err.Error())
			return err
		}
	} else {
		if err := doPlanTerminal(&cluster); err != nil {
			c.errClusterDelete(&cluster, "uninstall cluster err: "+err.Error())
			return err
		}
	}

	err := c.clusterRepo.Delete(cluster.Name)
	if err != nil {
		log.Error(err)
		c.errClusterDelete(&cluster, "delete cluster err: "+err.Error())
		return err
	}
	_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterUnInstall, true, ""), cluster.Name, constant.ClusterUnInstall)
	return nil
}

const terminalPlaybookName = "99-reset-cluster.yml"

func doPlanTerminal(cluster *model.Cluster) error {
	logId, _, err := ansible.CreateAnsibleLogWriter(cluster.Name)
	if err != nil {
		log.Error(err)
	}
	cluster.LogId = logId
	_ = db.DB.Save(cluster)

	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	if _, err := k.Destroy(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func doBareMetalTerminal(cluster *model.Cluster) error {
	logId, writer, err := ansible.CreateAnsibleLogWriter(cluster.Name)
	if err != nil {
		log.Error(err)
	}
	cluster.LogId = logId
	_ = db.DB.Save(cluster)

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
	if err := phases.RunPlaybookAndGetResult(k, terminalPlaybookName, writer); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c *clusterTerminalService) errClusterDelete(cluster *model.Cluster, errStr string) {
	cluster.Status.Phase = constant.ClusterFailed
	cluster.Status.Message = errStr
	if len(cluster.Status.ClusterStatusConditions) == 1 {
		cluster.Status.ClusterStatusConditions[0].Status = constant.ConditionFalse
		cluster.Status.ClusterStatusConditions[0].Message = errStr
	}
	_ = c.clusterStatusRepo.Save(&cluster.Status)
	_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUpgrade, false, errStr), cluster.Name, constant.ClusterUpgrade)
}
