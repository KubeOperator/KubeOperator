package service

import (
	"errors"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubeconfig"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/KubeOperator/KubeOperator/pkg/util/webkubectl"
)

type ClusterService interface {
	Get(name string) (dto.Cluster, error)
	GetStatus(name string) (dto.ClusterStatus, error)
	GetSecrets(name string) (dto.ClusterSecret, error)
	GetSpec(name string) (dto.ClusterSpec, error)
	GetPlan(name string) (dto.Plan, error)
	GetApiServerEndpoint(name string) (dto.Endpoint, error)
	GetRouterEndpoint(name string) (dto.Endpoint, error)
	GetWebkubectlToken(name string) (dto.WebkubectlToken, error)
	GetKubeconfig(name string) (string, error)
	Delete(name string) error
	Create(creation dto.ClusterCreate) (dto.Cluster, error)
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
		clusterTerminalService:     NewCLusterTerminalService(),
		projectRepository:          repository.NewProjectRepository(),
		projectResourceRepository:  repository.NewProjectResourceRepository(),
		clusterLogService:          NewClusterLogService(),
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
	clusterTerminalService     ClusterTerminalService
	projectRepository          repository.ProjectRepository
	projectResourceRepository  repository.ProjectResourceRepository
	clusterLogService          ClusterLogService
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
	if cs.KubernetesToken == "" {
		err := c.clusterInitService.GatherKubernetesToken(cluster)
		if err != nil {
			return secret, err
		}
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

func (c clusterService) Create(creation dto.ClusterCreate) (dto.Cluster, error) {
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
	cluster.Spec = spec
	cluster.Status = status
	cluster.Secret = secret
	if cluster.Spec.Provider != constant.ClusterProviderBareMetal {
		cluster.Spec.WorkerAmount = creation.WorkerAmount
		plan, err := c.planRepo.Get(creation.Plan)
		if err != nil {
			return dto.Cluster{}, err
		}
		cluster.PlanID = plan.ID
	}
	workerNo := 1
	masterNo := 1
	for _, nc := range creation.Nodes {
		node := model.ClusterNode{
			ClusterID: cluster.ID,
			Role:      nc.Role,
		}
		switch node.Role {
		case constant.NodeRoleNameMaster:
			node.Name = fmt.Sprintf("%s-%s-%d", cluster.Name, constant.NodeRoleNameMaster, masterNo)
			masterNo++
		case constant.NodeRoleNameWorker:
			node.Name = fmt.Sprintf("%s-%s-%d", cluster.Name, constant.NodeRoleNameWorker, workerNo)
			workerNo++
		}
		host, err := c.hostRepo.Get(nc.HostName)
		if err != nil {
			return dto.Cluster{}, err
		}
		node.HostID = host.ID
		node.Host = host
		cluster.Nodes = append(cluster.Nodes, node)
	}
	if len(cluster.Nodes) > 0 {
		cluster.Spec.KubeRouter = cluster.Nodes[0].Host.Ip
	}
	if err := c.clusterRepo.Save(&cluster); err != nil {
		return dto.Cluster{}, err
	}
	project, err := c.projectRepository.Get(creation.ProjectName)
	if err != nil {
		return dto.Cluster{}, err

	}
	if err := c.projectResourceRepository.Create(model.ProjectResource{
		ResourceId:   cluster.ID,
		ProjectID:    project.ID,
		ResourceType: constant.ResourceCluster,
	}); err != nil {
		return dto.Cluster{}, err
	}
	if err := c.clusterInitService.Init(cluster.Name); err != nil {
		return dto.Cluster{}, err
	}
	return dto.Cluster{Cluster: cluster}, nil
}

func (c clusterService) GetApiServerEndpoint(name string) (dto.Endpoint, error) {
	cluster, err := c.clusterRepo.Get(name)
	var endpoint dto.Endpoint
	if err != nil {
		return endpoint, err
	}
	endpoint.Port = cluster.Spec.KubeApiServerPort
	if cluster.Spec.LbKubeApiserverIp != "" {
		endpoint.Address = cluster.Spec.LbKubeApiserverIp
		return endpoint, nil
	}
	master, err := c.clusterNodeRepo.FistMaster(cluster.ID)
	if err != nil {
		return endpoint, err
	}
	endpoint.Address = master.Host.Ip
	return endpoint, nil
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
	endpoint, err := c.GetApiServerEndpoint(name)
	if err != nil {
		return token, err
	}
	addr := fmt.Sprintf("https://%s:%d", endpoint.Address, endpoint.Port)
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

func (c clusterService) Delete(name string) error {
	cluster, err := c.Get(name)
	if err != nil {
		return err
	}
	switch cluster.Source {
	case constant.ClusterSourceLocal:
		switch cluster.Status {
		case constant.ClusterRunning:
			go c.clusterTerminalService.Terminal(cluster.Cluster)
		case constant.ClusterCreating, constant.ClusterInitializing:
			return errors.New("CLUSTER_DELETE_FAILED")
		case constant.ClusterFailed:
			if cluster.Spec.Provider == constant.ClusterProviderPlan {
				var hosts []model.Host
				db.DB.Where(model.Host{ClusterID: cluster.ID}).Find(&hosts)
				if len(hosts) > 0 {
					go c.clusterTerminalService.Terminal(cluster.Cluster)
				} else {
					_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterUnInstall, true, ""), cluster.Name, constant.ClusterUnInstall)
					err = c.clusterRepo.Delete(name)
					if err != nil {
						return err
					}
				}
			} else {
				go c.clusterTerminalService.Terminal(cluster.Cluster)
			}
		}
	case constant.ClusterSourceExternal:
		_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterUnInstall, true, ""), cluster.Name, constant.ClusterUnInstall)
		err = c.clusterRepo.Delete(name)
		if err != nil {
			return err
		}
	}
	return nil
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
	m, err := c.clusterNodeRepo.FistMaster(cluster.ID)
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
