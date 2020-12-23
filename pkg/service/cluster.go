package service

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubeconfig"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/KubeOperator/KubeOperator/pkg/util/webkubectl"
)

type ClusterService interface {
	Get(name string) (dto.Cluster, error)
	GetStatus(name string) (dto.ClusterStatus, error)
	GetSecrets(name string) (dto.ClusterSecret, error)
	GetSpec(name string) (dto.ClusterSpec, error)
	GetPlan(name string) (dto.Plan, error)
	GetApiServerEndpoint(name string) (kubernetes.Host, error)
	GetApiServerEndpoints(name string) ([]kubernetes.Host, error)
	GetRouterEndpoint(name string) (dto.Endpoint, error)
	GetWebkubectlToken(name string) (dto.WebkubectlToken, error)
	GetKubeconfig(name string) (string, error)
	Delete(name string) error
	Create(creation dto.ClusterCreate) (*dto.Cluster, error)
	List() ([]dto.Cluster, error)
	Page(num, size int, projectName string) (dto.ClusterPage, error)
	Batch(batch dto.ClusterBatch) error
}

func NewClusterService() ClusterService {
	return &clusterService{
		clusterRepo:                repository.NewClusterRepository(),
		clusterSpecRepo:            repository.NewClusterSpecRepository(),
		clusterNodeRepo:            repository.NewClusterNodeRepository(),
		clusterStatusRepo:          repository.NewClusterStatusRepository(),
		clusterSecretRepo:          repository.NewClusterSecretRepository(),
		clusterStatusConditionRepo: repository.NewClusterStatusConditionRepository(),
		hostRepo:                   repository.NewHostRepository(),
		clusterInitService:         NewClusterInitService(),
		planRepo:                   repository.NewPlanRepository(),
		projectRepository:          repository.NewProjectRepository(),
		projectResourceRepository:  repository.NewProjectResourceRepository(),
		messageService:             NewMessageService(),
	}
}

type clusterService struct {
	clusterRepo                repository.ClusterRepository
	clusterSpecRepo            repository.ClusterSpecRepository
	clusterNodeRepo            repository.ClusterNodeRepository
	clusterStatusRepo          repository.ClusterStatusRepository
	clusterSecretRepo          repository.ClusterSecretRepository
	clusterStatusConditionRepo repository.ClusterStatusConditionRepository
	hostRepo                   repository.HostRepository
	planRepo                   repository.PlanRepository
	clusterInitService         ClusterInitService
	projectRepository          repository.ProjectRepository
	projectResourceRepository  repository.ProjectResourceRepository
	messageService             MessageService
}

func (c clusterService) Get(name string) (dto.Cluster, error) {
	var clusterDTO dto.Cluster
	mo, err := c.clusterRepo.Get(name)
	if err != nil {
		return clusterDTO, err
	}
	clusterDTO.Provider = mo.Spec.Provider
	clusterDTO.Cluster = mo
	clusterDTO.NodeSize = len(mo.Nodes)
	clusterDTO.Status = mo.Status.Phase
	clusterDTO.PreStatus = mo.Status.PrePhase
	clusterDTO.Architectures = mo.Spec.Architectures
	if len(mo.MultiClusterRepositories) > 0 {
		clusterDTO.MultiClusterRepository = mo.MultiClusterRepositories[0].Name
	}

	return clusterDTO, nil
}

func (c clusterService) List() ([]dto.Cluster, error) {
	var clusterDTOS []dto.Cluster
	mos, err := c.clusterRepo.List()
	if err != nil {
		return clusterDTOS, nil
	}
	for _, mo := range mos {
		clusterDTO := dto.Cluster{
			Cluster:       mo,
			NodeSize:      len(mo.Nodes),
			Status:        mo.Status.Phase,
			Provider:      mo.Spec.Provider,
			PreStatus:     mo.Status.PrePhase,
			Architectures: mo.Spec.Architectures,
		}
		if len(mo.MultiClusterRepositories) > 0 {
			clusterDTO.MultiClusterRepository = mo.MultiClusterRepositories[0].Name
		}
		clusterDTOS = append(clusterDTOS, clusterDTO)
	}
	return clusterDTOS, err
}

func (c clusterService) Page(num, size int, projectName string) (dto.ClusterPage, error) {
	var page dto.ClusterPage
	total, mos, err := c.clusterRepo.Page(num, size, projectName)
	if err != nil {
		return page, nil
	}
	for _, mo := range mos {
		clusterDTO := dto.Cluster{
			Cluster:       mo,
			NodeSize:      len(mo.Nodes),
			Status:        mo.Status.Phase,
			Provider:      mo.Spec.Provider,
			PreStatus:     mo.Status.PrePhase,
			Architectures: mo.Spec.Architectures,
		}
		if len(mo.MultiClusterRepositories) > 0 {
			clusterDTO.MultiClusterRepository = mo.MultiClusterRepositories[0].Name
		}
		page.Items = append(page.Items, clusterDTO)
	}
	page.Total = total
	return page, err
}

func (c clusterService) GetSecrets(name string) (dto.ClusterSecret, error) {
	var secret dto.ClusterSecret
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return secret, err
	}
	cs, err := c.clusterSecretRepo.Get(cluster.SecretID)
	if err != nil {
		return secret, err
	}
	secret.ClusterSecret = cs

	return secret, nil
}

func (c clusterService) GetStatus(name string) (dto.ClusterStatus, error) {
	var status dto.ClusterStatus
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return status, err
	}
	cs, err := c.clusterStatusRepo.Get(cluster.StatusID)
	if err != nil {
		return status, err
	}
	status.ClusterStatus = cs
	return status, nil
}

func (c clusterService) GetSpec(name string) (dto.ClusterSpec, error) {
	var spec dto.ClusterSpec
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return spec, err
	}
	cs, err := c.clusterSpecRepo.Get(cluster.SpecID)
	if err != nil {
		return spec, err
	}
	spec.ClusterSpec = cs
	return spec, nil
}

func (c clusterService) GetPlan(name string) (dto.Plan, error) {
	var plan dto.Plan
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return plan, err
	}
	p, err := c.planRepo.GetById(cluster.PlanID)
	if err != nil {
		return plan, err
	}
	plan.Plan = p
	return plan, nil
}

func (c clusterService) Create(creation dto.ClusterCreate) (*dto.Cluster, error) {
	cluster := model.Cluster{
		Name:   creation.Name,
		Source: constant.ClusterSourceLocal,
	}
	spec := model.ClusterSpec{
		RuntimeType:           creation.RuntimeType,
		DockerStorageDir:      creation.DockerStorageDIr,
		ContainerdStorageDir:  creation.ContainerdStorageDIr,
		NetworkType:           creation.NetworkType,
		KubePodSubnet:         creation.KubePodSubnet,
		KubeServiceSubnet:     creation.KubeServiceSubnet,
		Version:               creation.Version,
		Provider:              creation.Provider,
		FlannelBackend:        creation.FlannelBackend,
		CalicoIpv4poolIpip:    creation.CalicoIpv4poolIpip,
		KubeMaxPods:           creation.KubeMaxPods,
		KubeProxyMode:         creation.KubeProxyMode,
		IngressControllerType: creation.IngressControllerType,
		Architectures:         creation.Architectures,
		KubernetesAudit:       creation.KubernetesAudit,
		DockerSubnet:          creation.DockerSubnet,
		KubeApiServerPort:     constant.DefaultApiServerPort,
		HelmVersion:           creation.HelmVersion,
		NetworkInterface:      creation.NetworkInterface,
		SupportGpu:            creation.SupportGpu,
	}

	status := model.ClusterStatus{Phase: constant.ClusterWaiting}
	secret := model.ClusterSecret{
		KubeadmToken: clusterUtil.GenerateKubeadmToken(),
	}
	tx := db.DB.Begin()
	if err := tx.Create(&spec).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Create(&status).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Create(&secret).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	cluster.SpecID = spec.ID
	cluster.StatusID = status.ID
	cluster.SecretID = secret.ID
	if err := tx.Create(&cluster).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	switch spec.Provider {
	case constant.ClusterProviderPlan:
		spec.WorkerAmount = creation.WorkerAmount
		var plan model.Plan
		if err := tx.Where(model.Plan{Name: creation.Plan}).First(&plan).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not query plan %s reason %s", creation.Plan, err.Error())
		}
		cluster.PlanID = plan.ID
		if err := tx.Save(&cluster).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	case constant.ClusterProviderBareMetal:
		workerNo := 1
		masterNo := 1
		for _, nc := range creation.Nodes {
			n := model.ClusterNode{
				ClusterID: cluster.ID,
				Role:      nc.Role,
			}
			switch n.Role {
			case constant.NodeRoleNameMaster:
				n.Name = fmt.Sprintf("%s-%s-%d", cluster.Name, constant.NodeRoleNameMaster, masterNo)
				masterNo++
			case constant.NodeRoleNameWorker:
				n.Name = fmt.Sprintf("%s-%s-%d", cluster.Name, constant.NodeRoleNameWorker, workerNo)
				workerNo++
			}
			var host model.Host
			if err := tx.Where(model.Host{Name: nc.HostName}).First(&host).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("can not query host %s reason %s", nc.HostName, err.Error())
			}
			host.ClusterID = cluster.ID
			if err := tx.Save(&host).Error; err != nil {
				return nil, fmt.Errorf("can not save  host %s reason %s", nc.HostName, err.Error())
			}
			n.HostID = host.ID
			if err := tx.Create(&n).Error; err != nil {
				return nil, fmt.Errorf("can not create  node %s reason %s", n.Name, err.Error())
			}
			n.Host = host
			cluster.Nodes = append(cluster.Nodes, n)
		}
		if len(cluster.Nodes) > 0 {
			spec.KubeRouter = cluster.Nodes[0].Host.Ip
		}
	}
	if err := tx.Save(&spec).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	var project model.Project
	if err := tx.Where(model.Project{Name: creation.ProjectName}).First(&project).Error; err != nil {
		return nil, fmt.Errorf("can not load project %s reason %s", project.Name, err.Error())
	}
	projectResource := model.ProjectResource{
		ResourceID:   cluster.ID,
		ProjectID:    project.ID,
		ResourceType: constant.ResourceCluster,
	}
	if err := tx.Create(&projectResource).Error; err != nil {
		return nil, fmt.Errorf("can not create project  %s resource reason %s", project.Name, err.Error())
	}
	for _, tool := range cluster.PrepareTools() {
		tool.ClusterID = cluster.ID
		err := tx.Create(&tool).Error
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not prepare cluster tool %s reason %s", tool.Name, err.Error())
		}
	}
	tx.Commit()
	if err := c.clusterInitService.Init(cluster.Name); err != nil {
		return nil, err
	}
	return &dto.Cluster{Cluster: cluster}, nil
}

func (c *clusterService) Delete(name string) error {
	cluster, err := c.Get(name)
	if err != nil {
		return fmt.Errorf("can not get cluster %s reason %s", name, err)
	}
	if cluster.Dirty {
		if err := db.DB.Delete(&cluster.Cluster).Error; err != nil {
			return fmt.Errorf("can not delete cluster %s reason %s", name, err)
		}
		return nil
	}
	switch cluster.Source {
	case constant.ClusterSourceLocal:
		switch cluster.Status {
		case constant.StatusRunning, constant.StatusLost, constant.StatusFailed:
			cluster.Cluster.Status.Phase = constant.StatusTerminating
			if err := db.DB.Save(&cluster.Cluster.Status).Error; err != nil {
				return fmt.Errorf("can not update cluster %s status", cluster.Name)
			}
			switch cluster.Spec.Provider {
			case constant.ClusterProviderBareMetal:
				go func() {
					log.Infof("start uninstall cluster %s", cluster.Name)
					if err := c.uninstallCluster(&cluster.Cluster); err != nil {
						log.Errorf("uninstall cluster %s error %s", cluster.Name, err.Error())
					}
					if err := db.DB.Delete(&cluster.Cluster).Error; err != nil {
						log.Errorf("delete cluster error %s", err.Error())
					}
					_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUnInstall, true, ""), cluster.Name, constant.ClusterUnInstall)
				}()
			case constant.ClusterProviderPlan:
				go func() {
					log.Infof("start destroy cluster %s", cluster.Name)
					if err := c.destroyCluster(&cluster.Cluster); err != nil {
						log.Errorf("destroy cluster %s error %s", cluster.Name, err.Error())
					}
					if err := db.DB.Delete(&cluster.Cluster).Error; err != nil {
						log.Errorf("delete cluster error %s", err.Error())
					}
					_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUnInstall, true, ""), cluster.Name, constant.ClusterUnInstall)
				}()
			}
		case constant.StatusCreating, constant.StatusInitializing:
			return fmt.Errorf("can not delete cluster %s in this  status %s", cluster.Name, cluster.Status)
		case constant.StatusTerminating:
			return fmt.Errorf("cluster %s already in status %s", cluster.Name, cluster.Status)
		}
	case constant.ClusterSourceExternal:
		if err := db.DB.Delete(&cluster.Cluster).Error; err != nil {
			return err
		}
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUnInstall, true, ""), cluster.Name, constant.ClusterUnInstall)
	}
	return nil
}

const terminalPlaybookName = "99-reset-cluster.yml"

func (c *clusterService) uninstallCluster(cluster *model.Cluster) error {
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
		return err
	}
	return nil
}

func (c *clusterService) destroyCluster(cluster *model.Cluster) error {
	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	_, err := k.Destroy()
	if err != nil {
		return err
	}
	return nil
}

func (c clusterService) GetApiServerEndpoint(name string) (kubernetes.Host, error) {
	var result kubernetes.Host
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return "", err
	}
	port := cluster.Spec.KubeApiServerPort
	if cluster.Spec.LbKubeApiserverIp != "" {
		result = kubernetes.Host(fmt.Sprintf("%s:%d", cluster.Spec.LbKubeApiserverIp, port))
		return result, nil
	}
	master, err := c.clusterNodeRepo.FirstMaster(cluster.ID)
	if err != nil {
		return "", err
	}
	result = kubernetes.Host(fmt.Sprintf("%s:%d", master.Host.Ip, port))
	return result, nil
}

func (c clusterService) GetApiServerEndpoints(name string) ([]kubernetes.Host, error) {
	var result []kubernetes.Host
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return nil, err
	}
	port := cluster.Spec.KubeApiServerPort
	if cluster.Spec.LbKubeApiserverIp != "" {
		result = append(result, kubernetes.Host(fmt.Sprintf("%s:%d", cluster.Spec.LbKubeApiserverIp, port)))
		return result, nil
	}
	masters, err := c.clusterNodeRepo.AllMaster(cluster.ID)
	if err != nil {
		return nil, err
	}
	for i := range masters {
		result = append(result, kubernetes.Host(fmt.Sprintf("%s:%d", masters[i].Host.Ip, port)))
	}
	return result, nil
}

func (c clusterService) GetRouterEndpoint(name string) (dto.Endpoint, error) {
	cluster, err := c.clusterRepo.Get(name)
	var endpoint dto.Endpoint
	if err != nil {
		return endpoint, err
	}
	endpoint.Address = cluster.Spec.KubeRouter
	return endpoint, nil
}

func (c clusterService) GetWebkubectlToken(name string) (dto.WebkubectlToken, error) {
	var token dto.WebkubectlToken
	endpoints, err := c.GetApiServerEndpoints(name)
	if err != nil {
		return token, err
	}
	aliveHost, err := kubernetes.SelectAliveHost(endpoints)
	if err != nil {
		return token, err
	}
	addr := fmt.Sprintf("https://%s", aliveHost)
	secret, err := c.GetSecrets(name)
	if err != nil {
		return token, nil
	}
	t, err := webkubectl.GetConnectToken(name, addr, secret.KubernetesToken)
	token.Token = t
	if err != nil {
		return token, nil
	}

	return token, nil
}

func (c clusterService) Batch(batch dto.ClusterBatch) error {
	switch batch.Operation {
	case constant.BatchOperationDelete:
		for _, item := range batch.Items {
			err := c.Delete(item.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c clusterService) GetKubeconfig(name string) (string, error) {
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return "", err
	}
	m, err := c.clusterNodeRepo.FirstMaster(cluster.ID)
	if err != nil {
		return "", err
	}
	cfg := m.ToSSHConfig()
	s, err := ssh.New(&cfg)
	if err != nil {
		return "", err
	}
	bf, err := kubeconfig.ReadKubeConfigFile(s)
	if err != nil {
		return "", err
	}
	return string(bf), nil
}
