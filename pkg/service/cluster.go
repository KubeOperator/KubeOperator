package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
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
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubeconfig"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/KubeOperator/KubeOperator/pkg/util/webkubectl"
)

type ClusterService interface {
	Get(name string) (dto.Cluster, error)
	GetClusterByProject(projectNames string) ([]dto.ClusterInfo, error)
	CheckExistence(name string) bool
	GetStatus(name string) (*dto.TaskLog, error)
	GetSecrets(name string) (dto.ClusterSecret, error)
	GetApiServerEndpoint(name string) (string, error)
	GetApiServerEndpoints(name string) ([]string, error)
	GetRouterEndpoint(name string) (dto.Endpoint, error)
	GetWebkubectlToken(name string) (dto.WebkubectlToken, error)
	GetKubeconfig(name string) (string, error)
	Create(creation dto.ClusterCreate) (*dto.Cluster, error)
	ReCreate(name string) error
	List() ([]dto.Cluster, error)
	Page(num, size int, isPolling string, user dto.SessionUser, conditions condition.Conditions) (*dto.ClusterPage, error)
	Delete(name string, force bool, uninstall bool) error
}

func NewClusterService() ClusterService {
	return &clusterService{
		clusterRepo:          repository.NewClusterRepository(),
		clusterNodeRepo:      repository.NewClusterNodeRepository(),
		clusterSecretRepo:    repository.NewClusterSecretRepository(),
		clusterInitService:   NewClusterInitService(),
		planRepo:             repository.NewPlanRepository(),
		msgService:           NewMsgService(),
		ntpServerRepo:        repository.NewNtpServerRepository(),
		systemSettingService: NewSystemSettingService(),
		clusterIaasService:   NewClusterIaasService(),
		tasklogService:       NewTaskLogService(),
	}
}

type clusterService struct {
	clusterRepo          repository.ClusterRepository
	clusterNodeRepo      repository.ClusterNodeRepository
	clusterSecretRepo    repository.ClusterSecretRepository
	planRepo             repository.PlanRepository
	clusterInitService   ClusterInitService
	msgService           MsgService
	ntpServerRepo        repository.NtpServerRepository
	systemSettingService SystemSettingService
	clusterIaasService   ClusterIaasService
	tasklogService       TaskLogService
}

func (c clusterService) Get(name string) (dto.Cluster, error) {
	var clusterDTO dto.Cluster
	mo, err := c.clusterRepo.GetWithPreload(name, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Nodes"})
	if err != nil {
		return clusterDTO, err
	}
	tasklog, _ := c.tasklogService.GetByID(mo.CurrentTaskID)
	clusterDTO.Cluster.TaskLog = tasklog
	clusterDTO.Provider = mo.Provider
	clusterDTO.Cluster = mo
	clusterDTO.NodeSize = len(mo.Nodes)
	clusterDTO.Architectures = mo.Architectures
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
	if err := db.DB.Where("name in (?)", projectList).Preload("Clusters").Preload("Clusters.SpecConf").Find(&projects).Error; err != nil {
		return nil, err
	}
	for _, pro := range projects {
		for _, clu := range pro.Clusters {
			backdatas = append(backdatas, dto.ClusterInfo{Name: clu.Name, Provider: clu.Provider})
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
			Cluster:  mo,
			NodeSize: len(mo.Nodes),
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
			Preload("SpecConf").
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
			Preload("SpecConf").
			Preload("Nodes").
			Preload("MultiClusterRepositories").
			Find(&clusters).Error; err != nil {
			return nil, err
		}
	}

	for i := 0; i < len(clusters); i++ {
		if (clusters[i].Status == constant.StatusRunning || clusters[i].Status == constant.ClusterNotReady) && !(isPolling == "true") {
			isOK := false
			isOK, clusters[i].Message = GetClusterStatusByAPI(fmt.Sprintf("%s:%d", clusters[i].SpecConf.LbKubeApiserverIp, clusters[i].SpecConf.KubeApiServerPort))
			if !isOK {
				_ = db.DB.Model(&model.Cluster{}).Where("id = ?", clusters[i].ID).Updates(map[string]interface{}{"Status": constant.ClusterNotReady, "Message": clusters[i].Message})
			}
			if isOK && clusters[i].Status == constant.ClusterNotReady {
				_ = db.DB.Model(&model.Cluster{}).Where("id = ?", clusters[i].ID).Updates(map[string]interface{}{"Status": constant.StatusRunning, "Message": ""})
			}
		}
		for _, res := range clusterResources {
			if clusters[i].ID == res.ResourceID {
				for _, pro := range projects {
					if pro.ID == res.ProjectID {
						clusterDTO := dto.Cluster{
							Cluster:     clusters[i],
							ProjectName: pro.Name,
							NodeSize:    len(clusters[i].Nodes),
						}
						if len(clusters[i].MultiClusterRepositories) > 0 {
							clusterDTO.MultiClusterRepository = clusters[i].MultiClusterRepositories[0].Name
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

// func (c clusterService) GetLogBy(name string) (*dto.TaskLog, error) {
// 	var cluster model.Cluster
// 	cluster, err := c.clusterRepo.GetWithPreload(name, []string{"TaskLog", "TaskLog.Details"})
// 	return &dto.TaskLog{TaskLog: cluster.TaskLog}, err
// }

func (c clusterService) GetStatus(name string) (*dto.TaskLog, error) {
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return nil, err
	}
	if cluster.CurrentTaskID == "" {
		return &dto.TaskLog{}, nil
	}
	tasklog, err := c.tasklogService.GetByID(cluster.CurrentTaskID)
	if err != nil {
		return nil, err
	}
	sort.Slice(tasklog.Details, func(i, j int) bool {
		return tasklog.Details[i].StartTime < tasklog.Details[j].StartTime
	})
	return &dto.TaskLog{TaskLog: tasklog}, nil
}

func (c *clusterService) ReCreate(name string) error {
	cluster, err := c.clusterRepo.GetWithPreload(name, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}
	tasklog, err := c.tasklogService.GetByID(cluster.CurrentTaskID)
	if err != nil {
		return err
	}
	cluster.TaskLog = tasklog
	if err := c.tasklogService.RestartTask(&cluster, constant.TaskLogTypeClusterCreate); err != nil {
		return err
	}
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, cluster.TaskLog.ID)
	if err != nil {
		return err
	}
	cluster.Status = constant.StatusWaiting
	_ = c.clusterRepo.Save(&cluster)

	logger.Log.WithFields(logrus.Fields{
		"log_id": cluster.TaskLog.ID,
	}).Debugf("get ansible writer log of cluster %s successful, now start to init the cluster", cluster.Name)

	go c.clusterInitService.Init(cluster, writer)
	return nil
}

func (c *clusterService) Delete(name string, force bool, uninstall bool) error {
	logger.Log.Infof("start to delete cluster %s, isforce: %v", name, force)
	go c.deleteKubePi(name)
	cluster, err := c.clusterRepo.GetWithPreload(name, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Nodes", "Nodes.Host", "Nodes.Host.Credential", "Nodes.Host.Zone", "MultiClusterRepositories"})
	if err != nil {
		return fmt.Errorf("can not get cluster %s reason %s", name, err)
	}

	// ko 导入集群执行删除时，是否卸载，卸载则走正常手动卸载模式，否则走导入集群删除逻辑直接删除数据库数据，但是需要删除主机和资源绑定信息
	if cluster.Source == constant.ClusterSourceKoExternal {
		if uninstall {
			cluster.Source = constant.ClusterSourceLocal
		} else {
			cluster.Source = constant.ClusterSourceExternal
			if err := db.DB.Model(&model.Cluster{}).Where("id = ?", cluster.ID).Update("provider", constant.ClusterProviderPlan).Error; err != nil {
				return err
			}
		}
	}

	switch cluster.Source {
	case constant.ClusterSourceLocal:
		switch cluster.Status {
		case constant.StatusRunning, constant.StatusLost, constant.StatusFailed, constant.StatusNotReady:
			tasklog, err := c.tasklogService.NewTerminalTask(cluster.ID, constant.TaskLogTypeClusterDelete)
			if err != nil {
				return fmt.Errorf("can not update cluster %s status", cluster.Name)
			}
			cluster.TaskLog = *tasklog
			cluster.CurrentTaskID = tasklog.ID
			cluster.Status = constant.StatusTerminating
			_ = c.clusterRepo.Save(&cluster)
			switch cluster.Provider {
			case constant.ClusterProviderBareMetal:
				go c.uninstallCluster(&cluster, force)
			case constant.ClusterProviderPlan:
				go c.destroyCluster(&cluster, force)
			}
		case constant.StatusCreating, constant.StatusInitializing:
			return fmt.Errorf("can not delete cluster %s in this  status %s", cluster.Name, cluster.Status)
		case constant.StatusTerminating:
			return fmt.Errorf("cluster %s already in status %s", cluster.Name, cluster.Status)
		}
	case constant.ClusterSourceExternal:
		_ = c.msgService.SendMsg(constant.ClusterDelete, constant.System, cluster, true, map[string]string{"detailName": cluster.Name})
		if err := db.DB.Delete(&cluster).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c *clusterService) errClusterDelete(cluster *model.Cluster, errStr error) {
	logger.Log.Infof("cluster %s delete failed: %+v", cluster.Name, errStr)
	cluster.Status = constant.StatusFailed
	cluster.Message = errStr.Error()
	_ = c.clusterRepo.Save(cluster)
	_ = c.tasklogService.End(&cluster.TaskLog, false, errStr.Error())

	_ = c.msgService.SendMsg(constant.ClusterDelete, constant.System, cluster, false, map[string]string{"errMsg": errStr.Error(), "detailName": cluster.Name})
}

const terminalPlaybookName = "99-reset-cluster.yml"

func (c *clusterService) uninstallCluster(cluster *model.Cluster, force bool) {
	logger.Log.Infof("start to uninstall cluster %s, isforce: %v", cluster.Name, force)
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, cluster.TaskLog.ID)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%+v", err))
	}

	inventory := cluster.ParseInventory()
	k := kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})
	for i := range facts.DefaultFacts {
		k.SetVar(i, facts.DefaultFacts[i])
	}
	k.SetVar(facts.ClusterNameFactName, cluster.Name)
	ntps, _ := c.ntpServerRepo.GetAddressStr()
	k.SetVar(facts.NtpServerFactName, ntps)

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
	_ = c.msgService.SendMsg(constant.ClusterDelete, constant.System, cluster, true, map[string]string{"detailName": cluster.Name})
	if err := c.tasklogService.End(&cluster.TaskLog, true, ""); err != nil {
		logger.Log.Errorf("update tasklog error %s", err.Error())
	}
	logger.Log.Infof("start clearing cluster data %s", cluster.Name)
	if err := db.DB.Delete(&cluster).Error; err != nil {
		logger.Log.Errorf("delete cluster error %s", err.Error())
		return
	}
}

func (c *clusterService) destroyCluster(cluster *model.Cluster, force bool) {
	logger.Log.Infof("start to destroy cluster %s, isforce: %v", cluster.Name, force)

	_, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, cluster.TaskLog.ID)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%+v", err))
	}
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
	_ = c.msgService.SendMsg(constant.ClusterDelete, constant.System, cluster, true, map[string]string{"detailName": cluster.Name})
	logger.Log.Infof("start clearing cluster data %s", cluster.Name)
	if err := db.DB.Delete(&cluster).Error; err != nil {
		c.errClusterDelete(cluster, err)
		return
	}
}

func (c clusterService) GetApiServerEndpoint(name string) (string, error) {
	var result string
	cluster, err := c.clusterRepo.GetWithPreload(name, []string{"SpecConf"})
	if err != nil {
		return "", err
	}
	port := cluster.SpecConf.KubeApiServerPort
	if cluster.SpecConf.LbKubeApiserverIp != "" {
		result = fmt.Sprintf("%s:%d", cluster.SpecConf.LbKubeApiserverIp, port)
		return result, nil
	}
	master, err := c.clusterNodeRepo.FirstMaster(cluster.ID)
	if err != nil {
		return "", err
	}
	result = fmt.Sprintf("%s:%d", master.Host.Ip, port)
	return result, nil
}

func (c clusterService) GetApiServerEndpoints(name string) ([]string, error) {
	var result []string
	cluster, err := c.clusterRepo.GetWithPreload(name, []string{"SpecConf"})
	if err != nil {
		return nil, err
	}
	port := cluster.SpecConf.KubeApiServerPort
	if cluster.SpecConf.LbKubeApiserverIp != "" {
		result = append(result, fmt.Sprintf("%s:%d", cluster.SpecConf.LbKubeApiserverIp, port))
		return result, nil
	}
	masters, err := c.clusterNodeRepo.AllMaster(cluster.ID)
	if err != nil {
		return nil, err
	}
	for i := range masters {
		result = append(result, fmt.Sprintf("%s:%d", masters[i].Host.Ip, port))
	}
	return result, nil
}

func (c clusterService) GetRouterEndpoint(name string) (dto.Endpoint, error) {
	cluster, err := c.clusterRepo.GetWithPreload(name, []string{"SpecConf"})
	var endpoint dto.Endpoint
	if err != nil {
		return endpoint, err
	}
	endpoint.Address = cluster.SpecConf.KubeRouter
	return endpoint, nil
}

func (c clusterService) GetWebkubectlToken(name string) (dto.WebkubectlToken, error) {
	var token dto.WebkubectlToken
	endpoints, err := c.GetApiServerEndpoints(name)
	if err != nil {
		return token, err
	}
	aliveHost, err := clusterUtil.SelectAliveHost(endpoints)
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
	cluster, err := c.clusterRepo.GetWithPreload(name, []string{"SpecConf"})
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

	lbAddr := fmt.Sprintf("%s:%d", cluster.SpecConf.LbKubeApiserverIp, cluster.SpecConf.KubeApiServerPort)
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
