package service

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubepi"
	"github.com/sirupsen/logrus"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
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
	GetClusterByProject(projectNames string) ([]dto.ClusterInfo, error)
	CheckExistence(name string) bool
	GetStatus(name string) (dto.ClusterStatus, error)
	GetNodeStatus(cluster, node string) (dto.ClusterStatus, error)
	GetSecrets(name string) (dto.ClusterSecret, error)
	GetSpec(name string) (dto.ClusterSpec, error)
	GetPlan(name string) (dto.Plan, error)
	GetApiServerEndpoint(name string) (kubernetes.Host, error)
	GetApiServerEndpoints(name string) ([]kubernetes.Host, error)
	GetRouterEndpoint(name string) (dto.Endpoint, error)
	GetWebkubectlToken(name string) (dto.WebkubectlToken, error)
	GetKubeconfig(name string) (string, error)
	Create(creation dto.ClusterCreate) (*dto.Cluster, error)
	List() ([]dto.Cluster, error)
	Page(num, size int, isPolling string, user dto.SessionUser, conditions condition.Conditions) (*dto.ClusterPage, error)
	Delete(name string, force bool, uninstall bool) error
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
		ntpServerRepo:              repository.NewNtpServerRepository(),
		systemSettingService:       NewSystemSettingService(),
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
	ntpServerRepo              repository.NtpServerRepository
	systemSettingService       SystemSettingService
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

func (c clusterService) GetClusterByProject(projectNames string) ([]dto.ClusterInfo, error) {
	var (
		projectList []string
		projects    []model.Project
		backdatas   []dto.ClusterInfo
	)
	if len(projectNames) != 0 {
		projectList = strings.Split(projectNames, ",")
	}
	if err := db.DB.Where("name in (?)", projectList).Preload("Clusters").Preload("Clusters.Spec").Find(&projects).Error; err != nil {
		return nil, err
	}
	for _, pro := range projects {
		for _, clu := range pro.Clusters {
			backdatas = append(backdatas, dto.ClusterInfo{Name: clu.Name, Provider: clu.Spec.Provider})
		}
	}
	return backdatas, nil
}

func (c clusterService) CheckExistence(name string) bool {
	count := 1
	_ = db.DB.Model(&model.Cluster{}).Where("name = ?", name).Count(&count)
	return count == 1
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

func (c clusterService) Page(num, size int, isPolling string, user dto.SessionUser, conditions condition.Conditions) (*dto.ClusterPage, error) {
	var (
		page             dto.ClusterPage
		clusters         []model.Cluster
		clusterResources []model.ProjectResource
		projects         []model.Project
	)
	d := db.DB.Model(model.Cluster{})
	if err := dbUtil.WithConditions(&d, model.Cluster{}, conditions); err != nil {
		return nil, err
	}
	if err := db.DB.Where("resource_type = 'CLUSTER'").Find(&clusterResources).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Find(&projects).Error; err != nil {
		return nil, err
	}

	if user.IsAdmin {
		if err := d.Count(&page.Total).Order("created_at ASC").
			Preload("Status").
			Preload("Spec").
			Preload("Nodes").
			Preload("MultiClusterRepositories").
			Offset((num - 1) * size).Limit(size).Find(&clusters).Error; err != nil {
			return nil, err
		}
	} else {
		var (
			clusterIds []string
			resources  []model.ProjectResource
		)
		for _, pro := range projects {
			if pro.Name == user.CurrentProject {
				for _, res := range clusterResources {
					if res.ProjectID == pro.ID {
						resources = append(resources, res)
					}
				}
			}
		}
		if user.IsRole(constant.RoleProjectManager) {
			for _, pm := range resources {
				clusterIds = append(clusterIds, pm.ResourceID)
			}
		} else {
			var resourceIds []string
			for _, pm := range resources {
				resourceIds = append(resourceIds, pm.ResourceID)
			}
			var clusterMembers []model.ClusterMember
			if err := db.DB.Raw("SELECT DISTINCT cluster_id FROM ko_cluster_member WHERE cluster_id in (?) AND user_id = ?", resourceIds, user.UserId).Scan(&clusterMembers).Error; err != nil {
				return nil, err
			}
			for _, pm := range clusterMembers {
				clusterIds = append(clusterIds, pm.ClusterID)
			}
		}
		if err := db.DB.Model(&model.Cluster{}).
			Where("id in (?)", clusterIds).
			Count(&page.Total).
			Offset((num - 1) * size).
			Limit(size).
			Order("created_at ASC").
			Preload("Status").
			Preload("Spec").
			Preload("Nodes").
			Preload("MultiClusterRepositories").
			Find(&clusters).Error; err != nil {
			return nil, err
		}
	}

	for _, mo := range clusters {
		status := mo.Status.Phase
		message := ""
		if (mo.Status.Phase == constant.ClusterRunning || mo.Status.Phase == constant.ClusterNotReady) && !(isPolling == "true") {
			isOK := false
			isOK, message = GetClusterStatusByAPI(fmt.Sprintf("%s:%d", mo.Spec.LbKubeApiserverIp, mo.Spec.KubeApiServerPort))
			if !isOK {
				status = constant.ClusterNotReady
				_ = db.DB.Model(&model.ClusterStatus{}).Where("id = ?", mo.StatusID).Updates(map[string]interface{}{"Phase": constant.ClusterNotReady, "Message": message})
			}
			if isOK && mo.Status.Phase == constant.ClusterNotReady {
				status = constant.ClusterRunning
				_ = db.DB.Model(&model.ClusterStatus{}).Where("id = ?", mo.StatusID).Updates(map[string]interface{}{"Phase": constant.ClusterRunning, "Message": ""})
			}
		}
		for _, res := range clusterResources {
			if mo.ID == res.ResourceID {
				for _, pro := range projects {
					if pro.ID == res.ProjectID {
						clusterDTO := dto.Cluster{
							Cluster:       mo,
							ProjectName:   pro.Name,
							NodeSize:      len(mo.Nodes),
							Status:        status,
							Provider:      mo.Spec.Provider,
							PreStatus:     mo.Status.PrePhase,
							Architectures: mo.Spec.Architectures,
							Message:       message,
						}
						if len(mo.MultiClusterRepositories) > 0 {
							clusterDTO.MultiClusterRepository = mo.MultiClusterRepositories[0].Name
						}
						page.Items = append(page.Items, clusterDTO)
						break
					}
				}
				break
			}
		}
	}
	return &page, nil
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

func (c clusterService) GetNodeStatus(clusterName, nodeName string) (dto.ClusterStatus, error) {
	var status dto.ClusterStatus
	node, err := c.clusterNodeRepo.Get(clusterName, nodeName)
	if err != nil {
		return status, err
	}
	cs, err := c.clusterStatusRepo.Get(node.StatusID)
	if err != nil {
		return status, err
	}
	status.ClusterStatus = cs
	status.ClusterStatus.Phase = node.Status
	status.ClusterStatus.PrePhase = node.PreStatus
	status.ClusterStatus.Message = node.Message
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

var maxNodePodNumMap = map[int]int{
	24: 110,
	25: 64,
	26: 32,
	27: 16,
}

func (c clusterService) Create(creation dto.ClusterCreate) (*dto.Cluster, error) {
	loginfo, _ := json.Marshal(creation)
	logger.Log.WithFields(logrus.Fields{"cluster_creation": string(loginfo)}).Debugf("start to create the cluster %s", creation.Name)

	cluster := model.Cluster{
		Name:         creation.Name,
		NodeNameRule: creation.NodeNameRule,
		Source:       constant.ClusterSourceLocal,
	}
	spec := model.ClusterSpec{
		RuntimeType:              creation.RuntimeType,
		DockerStorageDir:         creation.DockerStorageDIr,
		ContainerdStorageDir:     creation.ContainerdStorageDIr,
		NetworkType:              creation.NetworkType,
		CiliumVersion:            creation.CiliumVersion,
		CiliumTunnelMode:         creation.CiliumTunnelMode,
		CiliumNativeRoutingCidr:  creation.CiliumNativeRoutingCidr,
		Version:                  creation.Version,
		Provider:                 creation.Provider,
		FlannelBackend:           creation.FlannelBackend,
		CalicoIpv4poolIpip:       creation.CalicoIpv4poolIpip,
		KubeProxyMode:            creation.KubeProxyMode,
		NodeportAddress:          creation.NodeportAddress,
		KubeServiceNodePortRange: creation.KubeServiceNodePortRange,
		EnableDnsCache:           creation.EnableDnsCache,
		DnsCacheVersion:          creation.DnsCacheVersion,
		IngressControllerType:    creation.IngressControllerType,
		Architectures:            creation.Architectures,
		KubeDnsDomain:            creation.KubeDnsDomain,
		KubernetesAudit:          creation.KubernetesAudit,
		DockerSubnet:             creation.DockerSubnet,
		HelmVersion:              creation.HelmVersion,
		NetworkInterface:         creation.NetworkInterface,
		NetworkCidr:              creation.NetworkCidr,
		SupportGpu:               creation.SupportGpu,
		YumOperate:               creation.YumOperate,
		LbMode:                   creation.LbMode,
		LbKubeApiserverIp:        creation.LbKubeApiserverIp,
		KubeApiServerPort:        creation.KubeApiServerPort,
		KubePodSubnet:            creation.KubePodSubnet,
		KubeServiceSubnet:        creation.KubeServiceSubnet,
		MaxNodeNum:               creation.MaxNodeNum,
		MasterScheduleType:       creation.MasterScheduleType,
	}

	nodeMask, err := getNodeCIDRMaskSize(creation.MaxNodePodNum)
	if err != nil {
		return nil, err
	}

	spec.KubeMaxPods = maxNodePodNumMap[nodeMask]
	spec.KubeNetworkNodePrefix = nodeMask

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
	var project model.Project
	if err := tx.Where("name = ?", creation.ProjectName).First(&project).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("can not load project %s reason %s", project.Name, err.Error())
	}
	cluster.ProjectID = project.ID
	if err := tx.Create(&cluster).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if creation.SupportGpu == constant.StatusEnabled {
		gpuInfo := &model.ClusterGpu{
			ClusterID: cluster.ID,
			Status:    constant.StatusEnabled,
		}
		if err := tx.Create(&gpuInfo).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	switch spec.Provider {
	case constant.ClusterProviderPlan:
		spec.WorkerAmount = creation.WorkerAmount
		var plan model.Plan
		if err := tx.Where("name = ?", creation.Plan).First(&plan).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not query plan %s reason %s", creation.Plan, err.Error())
		}
		cluster.PlanID = plan.ID
		if err := tx.Save(&cluster).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		clusterResource := model.ClusterResource{
			ResourceID:   plan.ID,
			ClusterID:    cluster.ID,
			ResourceType: constant.ResourcePlan,
		}
		if err := tx.Create(&clusterResource).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not create cluster  %s resource reason %s", cluster.Name, err.Error())
		}
	case constant.ClusterProviderBareMetal:
		workerNo := 1
		masterNo := 1
		firstMasterIP := ""
		for _, nc := range creation.Nodes {
			n := model.ClusterNode{
				ClusterID: cluster.ID,
				Role:      nc.Role,
			}

			var host model.Host
			if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("name = ?", nc.HostName).First(&host).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("can not find host %s cluster id ", nc.HostName)
			}
			host.ClusterID = cluster.ID
			if err := tx.Save(&host).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("can not update host %s cluster id ", nc.HostName)
			}

			clusterResource := model.ClusterResource{
				ClusterID:    cluster.ID,
				ResourceID:   host.ID,
				ResourceType: constant.ResourceHost,
			}
			if err := tx.Create(&clusterResource).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("can bind host %s to cluster", nc.HostName)
			}

			switch cluster.NodeNameRule {
			case constant.NodeNameRuleDefault:
				if n.Role == constant.NodeRoleNameMaster {
					n.Name = fmt.Sprintf("%s-%s-%d", cluster.Name, constant.NodeRoleNameMaster, masterNo)
					if len(firstMasterIP) == 0 {
						firstMasterIP = n.Host.Ip
					}
					masterNo++
				} else {
					n.Name = fmt.Sprintf("%s-%s-%d", cluster.Name, constant.NodeRoleNameWorker, workerNo)
					workerNo++
				}
			case constant.NodeNameRuleIP:
				n.Name = host.Ip
			case constant.NodeNameRuleHostName:
				n.Name = host.Name
			}

			n.HostID = host.ID
			if err := tx.Create(&n).Error; err != nil {
				return nil, fmt.Errorf("can not create  node %s reason %s", n.Name, err.Error())
			}
			n.Host = host
			cluster.Nodes = append(cluster.Nodes, n)
		}
		if cluster.Spec.LbMode == constant.LbModeInternal {
			cluster.Spec.LbKubeApiserverIp = firstMasterIP
		}
		spec.KubeRouter = firstMasterIP
	}
	if err := tx.Save(&spec).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	projectResource := model.ProjectResource{
		ResourceID:   cluster.ID,
		ProjectID:    project.ID,
		ResourceType: constant.ResourceCluster,
	}
	if err := tx.Create(&projectResource).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("can not create project  %s resource reason %s", project.Name, err.Error())
	}

	var (
		manifest model.ClusterManifest
		toolVars []model.VersionHelp
	)
	if err := tx.Where("name = ?", spec.Version).First(&manifest).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("can find manifest version: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(manifest.ToolVars), &toolVars); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("unmarshal manifest.toolvar error %s", err.Error())
	}
	for _, tool := range cluster.PrepareTools() {
		for _, item := range toolVars {
			if tool.Name == item.Name {
				tool.Version = item.Version
				break
			}
		}
		tool.ClusterID = cluster.ID
		err := tx.Create(&tool).Error
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("can not prepare cluster tool %s reason %s", tool.Name, err.Error())
		}
	}

	if spec.Architectures == "amd64" {
		for _, istio := range cluster.PrepareIstios() {
			istio.ClusterID = cluster.ID
			err := tx.Create(&istio).Error
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("can not prepare cluster istio %s reason %s", istio.Name, err.Error())
			}
		}
	}
	tx.Commit()

	logger.Log.Infof("init db data of cluster %s successful, now start to create cluster", cluster.Name)
	if err := c.clusterInitService.Init(cluster.Name); err != nil {
		return nil, err
	}
	return &dto.Cluster{Cluster: cluster}, nil
}

func getNodeCIDRMaskSize(maxNodePodNum int) (int, error) {
	nodeCidrOccupy := math.Ceil(math.Log2(float64(maxNodePodNum)))
	nodeCIDRMaskSize := 32 - int(nodeCidrOccupy)
	return nodeCIDRMaskSize, nil
}

func (c *clusterService) Delete(name string, force bool, uninstall bool) error {
	logger.Log.Infof("start to delete cluster %s, isforce: %v", name, force)
	go c.deleteKubePi(name)
	cluster, err := c.Get(name)
	if err != nil {
		return fmt.Errorf("can not get cluster %s reason %s", name, err)
	}

	// ko 导入集群执行删除时，是否卸载，卸载则走正常手动卸载模式，否则走导入集群删除逻辑直接删除数据库数据，但是需要删除主机和资源绑定信息
	if cluster.Source == constant.ClusterSourceKoExternal {
		if uninstall {
			cluster.Source = constant.ClusterSourceLocal
		} else {
			cluster.Source = constant.ClusterSourceExternal
			if err := db.DB.Model(&model.ClusterSpec{}).Where("id = ?", cluster.Spec.ID).Update("provider", constant.ClusterProviderPlan).Error; err != nil {
				return err
			}
		}
	}

	switch cluster.Source {
	case constant.ClusterSourceLocal:
		switch cluster.Status {
		case constant.StatusRunning, constant.StatusLost, constant.StatusFailed, constant.StatusNotReady:
			cluster.Cluster.Status.Phase = constant.StatusTerminating
			cluster.Cluster.Status.ClusterStatusConditions = []model.ClusterStatusCondition{}
			condition := model.ClusterStatusCondition{
				Name:          "DeleteCluster",
				Status:        constant.ConditionUnknown,
				OrderNum:      0,
				LastProbeTime: time.Now(),
			}
			cluster.Cluster.Status.ClusterStatusConditions = append(cluster.Cluster.Status.ClusterStatusConditions, condition)
			if err := c.clusterStatusRepo.Save(&cluster.Cluster.Status); err != nil {
				return fmt.Errorf("can not update cluster %s status", cluster.Name)
			}
			switch cluster.Spec.Provider {
			case constant.ClusterProviderBareMetal:
				go c.uninstallCluster(&cluster.Cluster, force)
			case constant.ClusterProviderPlan:
				go c.destroyCluster(&cluster.Cluster, force)
			}
		case constant.StatusCreating, constant.StatusInitializing:
			return fmt.Errorf("can not delete cluster %s in this  status %s", cluster.Name, cluster.Status)
		case constant.StatusTerminating:
			return fmt.Errorf("cluster %s already in status %s", cluster.Name, cluster.Status)
		}
	case constant.ClusterSourceExternal:
		_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterDelete, true, ""), cluster.Name, constant.ClusterDelete)
		if err := db.DB.Delete(&cluster.Cluster).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c *clusterService) errClusterDelete(cluster *model.Cluster, errStr error) {
	logger.Log.Infof("cluster %s delete failed: %+v", cluster.Name, errStr)
	cluster.Status.Phase = constant.ClusterFailed
	cluster.Status.Message = errStr.Error()
	if len(cluster.Status.ClusterStatusConditions) == 1 {
		cluster.Status.ClusterStatusConditions[0].Status = constant.ConditionFalse
		cluster.Status.ClusterStatusConditions[0].Message = errStr.Error()
	}
	_ = c.clusterStatusRepo.Save(&cluster.Status)
	_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUnInstall, false, errStr.Error()), cluster.Name, constant.ClusterUnInstall)
}

const terminalPlaybookName = "99-reset-cluster.yml"

func (c *clusterService) uninstallCluster(cluster *model.Cluster, force bool) {
	logger.Log.Infof("start to uninstall cluster %s, isforce: %v", cluster.Name, force)
	logId, writer, err := ansible.CreateAnsibleLogWriter(cluster.Name)
	if err != nil {
		logger.Log.Error(err)
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
	k.SetVar(facts.ClusterNameFactName, cluster.Name)
	ntps, _ := c.ntpServerRepo.GetAddressStr()
	k.SetVar(facts.NtpServerName, ntps)

	vars := cluster.GetKobeVars()
	for key, value := range vars {
		k.SetVar(key, value)
	}
	err = phases.RunPlaybookAndGetResult(k, terminalPlaybookName, "", writer)
	if err != nil {
		if force {
			logger.Log.Errorf("destroy cluster %s error %s", cluster.Name, err.Error())
			if err := db.DB.Delete(&cluster).Error; err != nil {
				c.errClusterDelete(cluster, err)
			}
		} else {
			c.errClusterDelete(cluster, err)
		}
		return
	}
	_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterUnInstall, true, ""), cluster.Name, constant.ClusterUnInstall)
	logger.Log.Infof("start clearing cluster data %s", cluster.Name)
	if err := db.DB.Delete(&cluster).Error; err != nil {
		logger.Log.Errorf("delete cluster error %s", err.Error())
		return
	}
}

func (c *clusterService) destroyCluster(cluster *model.Cluster, force bool) {
	logger.Log.Infof("start to destroy cluster %s, isforce: %v", cluster.Name, force)
	logId, _, err := ansible.CreateAnsibleLogWriter(cluster.Name)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%+v", err))
	}
	cluster.LogId = logId
	_ = db.DB.Save(cluster)
	plan, _ := c.planRepo.GetById(cluster.PlanID)
	k := kotf.NewTerraform(&kotf.Config{Cluster: cluster.Name})
	_, err = k.Destroy(plan.Region.Vars)
	if err != nil {
		if force {
			logger.Log.Errorf("destroy cluster %s error %s", cluster.Name, err.Error())
			if err := db.DB.Delete(&cluster).Error; err != nil {
				c.errClusterDelete(cluster, err)
			}
		} else {
			c.errClusterDelete(cluster, err)
		}
		return
	}
	_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterUnInstall, true, ""), cluster.Name, constant.ClusterUnInstall)
	logger.Log.Infof("start clearing cluster data %s", cluster.Name)
	if err := db.DB.Delete(&cluster).Error; err != nil {
		c.errClusterDelete(cluster, err)
		return
	}
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
	bf, err := kubeconfig.ReadKubeConfigFile(s, m.Host.Credential.Username)
	if err != nil {
		return "", err
	}
	configStr := string(bf)

	lbAddr := fmt.Sprintf("%s:%d", cluster.Spec.LbKubeApiserverIp, cluster.Spec.KubeApiServerPort)
	newStr := strings.ReplaceAll(configStr, "127.0.0.1:8443", lbAddr)

	return newStr, nil
}

func (c clusterService) deleteKubePi(name string) {
	logger.Log.Infof("start to delete kubepi client info")
	ss, err := c.systemSettingService.ListByTab("KUBEPI")
	if err != nil {
		logger.Log.Errorf("get kubepi login info failed, err: %v", err)
		return
	}
	apiServer, err := c.GetApiServerEndpoint(name)
	if err != nil {
		logger.Log.Errorf("get api server endpoint failed, err: %v", err)
		return
	}
	kubepiClient := kubepi.GetClient()
	if _, ok := ss.Vars["KUBEPI_USERNAME"]; ok {
		username := ss.Vars["KUBEPI_USERNAME"]
		password := ss.Vars["KUBEPI_PASSWORD"]
		if username != "" && password != "" {
			kubepiClient = kubepi.GetClient(kubepi.WithUsernameAndPassword(username, password))
		}
	}
	if err := kubepiClient.Close(name, string(apiServer)); err != nil {
		logger.Log.Errorf("close kubepi client failed, err: %v", err)
		return
	}

	logger.Log.Infof("delete kubepi client info success")
}
