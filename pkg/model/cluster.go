package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
)

type Cluster struct {
	common.BaseModel
	ID             string `json:"-"`
	Name           string `json:"name" gorm:"not null;unique"`
	NodeNameRule   string `json:"nodeNameRule"`
	Source         string `json:"source"`
	Version        string `json:"version"`
	UpgradeVersion string `json:"upgradeVersion"`
	Provider       string `json:"provider"`
	Architectures  string `json:"architectures"`

	Status  string `json:"status"`
	Message string `json:"message" gorm:"type:text(65535)"`

	CurrentTaskID string `json:"currentTaskID"`
	SecretID      string `json:"-"`
	PlanID        string `json:"-"`
	ProjectID     string `json:"projectID"`
	Dirty         bool   `json:"dirty"`
	Plan          Plan   `json:"-"`

	SpecConf                 ClusterSpecConf          `gorm:"save_associations:false" json:"specConf"`
	SpecRuntime              ClusterSpecRuntime       `gorm:"save_associations:false" json:"specRuntime"`
	SpecNetwork              ClusterSpecNetwork       `gorm:"save_associations:false" json:"specNetwork"`
	SpecComponent            []ClusterSpecComponent   `gorm:"save_associations:false" json:"specComponent"`
	TaskLog                  TaskLog                  `gorm:"save_associations:false" json:"taskLog"`
	Secret                   ClusterSecret            `gorm:"save_associations:false" json:"-"`
	Nodes                    []ClusterNode            `gorm:"save_associations:false" json:"-"`
	Tools                    []ClusterTool            `gorm:"save_associations:false" json:"-"`
	MultiClusterRepositories []MultiClusterRepository `gorm:"many2many:cluster_multi_cluster_repository"`
}

func (c *Cluster) BeforeCreate() error {
	c.ID = uuid.NewV4().String()
	return nil
}

func (c Cluster) BeforeDelete() error {
	var cluster Cluster
	cluster.ID = c.ID
	tx := db.DB.Begin()
	if err := tx.
		Preload("SpecConf").
		Preload("SpecRuntime").
		Preload("SpecNetwork").
		Preload("Nodes").
		Preload("Nodes.Host").
		Preload("Tools").
		Preload("MultiClusterRepositories").
		First(&cluster).Error; err != nil {
		return err
	}
	if err := tx.Where("cluster_id = ?", cluster.ID).Delete(&ClusterSpecConf{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("cluster_id = ?", cluster.ID).Delete(&ClusterSpecRuntime{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("cluster_id = ?", cluster.ID).Delete(&ClusterSpecNetwork{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("cluster_id = ?", cluster.ID).Delete(&ClusterSpecComponent{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	var (
		hostIDList []string
		hostIPList []string
	)
	for _, node := range cluster.Nodes {
		hostIDList = append(hostIDList, node.HostID)
		hostIPList = append(hostIPList, node.Host.Ip)
	}
	if cluster.Provider == constant.ClusterProviderPlan {
		if len(hostIDList) > 0 {
			if err := tx.Where("resource_id in (?) AND resource_type = ?", hostIDList, constant.ResourceHost).
				Delete(&ProjectResource{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("cluster_id = ?", c.ID).Delete(&Host{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		if len(hostIPList) > 0 {
			if err := tx.Model(&Ip{}).Where("address in (?)", hostIPList).Update("status", constant.IpAvailable).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		if err := tx.Model(&Host{}).Where("cluster_id = ?", c.ID).Updates(map[string]interface{}{"ClusterID": ""}).Error; err != nil {
			return err
		}
	}
	if len(hostIDList) > 0 {
		if err := tx.Where("cluster_id = ?", c.ID).Delete(&ClusterNode{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where("resource_id = ?", c.ID).Delete(&ProjectResource{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("cluster_id = ?", cluster.ID).Delete(&ClusterResource{}).Error; err != nil {
		return err
	}

	tx.Commit()
	go cluster.DeleteClusterAbout()
	return nil
}

func (c Cluster) DeleteClusterAbout() {
	if err := db.DB.Delete(&ClusterSecret{ID: c.SecretID}).Error; err != nil {
		logger.Log.Infof("delete secret failed, err: %v", err)
	}
	if err := db.DB.Where("cluster = ?", c.Name).Delete(&KubepiBind{}).Error; err != nil {
		logger.Log.Infof("delete kubepi bind failed, err: %v", err)
	}
	if err := db.DB.Where("cluster_id = ?", c.ID).Delete(&ClusterTool{}).Error; err != nil {
		logger.Log.Infof("delete tools failed, err: %v", err)
	}
	if err := db.DB.Where("cluster_id = ?", c.ID).Delete(&TaskLog{}).Error; err != nil {
		logger.Log.Infof("delete kubepi bind failed, err: %v", err)
	}
	if err := db.DB.Where("cluster_id = ?", c.ID).Delete(&TaskRetryLog{}).Error; err != nil {
		logger.Log.Infof("delete kubepi bind failed, err: %v", err)
	}

	var cisTasks []CisTask
	if err := db.DB.Where("cluster_id = ?", c.ID).Find(&cisTasks).Error; err != nil {
		logger.Log.Infof("delete cis tasks failed, err: %v", err)
	}
	if err := db.DB.Where("cluster_id = ?", c.ID).Delete(&ClusterStorageProvisioner{}).Error; err != nil {
		logger.Log.Infof("delete provisioner failed, err: %v", err)
	}

	if err := db.DB.Where("cluster_id = ?", c.ID).Delete(&ClusterBackupStrategy{}).Error; err != nil {
		logger.Log.Infof("delete backup strategy failed, err: %v", err)
	}

	if err := db.DB.Where("cluster_id = ?", c.ID).Delete(&ClusterBackupFile{}).Error; err != nil {
		logger.Log.Infof("delete backup file failed, err: %v", err)
	}

	if err := db.DB.Where("resource_id = ?", c.ID).Delete(&MsgSubscribe{}).Error; err != nil {
		logger.Log.Infof("delete backup file failed, err: %v", err)
	}

	var (
		messages   []Msg
		messageIDs []string
	)
	if err := db.DB.Where("resource_id = ? AND type != ?", c.ID, constant.System).Find(&messages).Error; err != nil {
		logger.Log.Infof("select message failed, err: %v", err)
	}
	for _, m := range messages {
		messageIDs = append(messageIDs, m.ID)
	}
	if len(messageIDs) > 0 {
		if err := db.DB.Where("msg_id in (?)", messageIDs).Delete(&UserMsg{}).Error; err != nil {
			logger.Log.Infof("delete user message failed, err: %v", err)
		}
	}
	if err := db.DB.Where("resource_id = ? AND type != ?", c.ID, constant.System).Delete(&Msg{}).Error; err != nil {
		logger.Log.Infof("delete message failed, err: %v", err)
	}

	if len(c.MultiClusterRepositories) > 0 {
		for _, repo := range c.MultiClusterRepositories {
			var clusterMultiClusterRepository ClusterMultiClusterRepository
			if err := db.DB.Where("cluster_id = ? AND multi_cluster_repository_id  = ?", c.ID, repo.ID).First(&clusterMultiClusterRepository).Error; err != nil {
				logger.Log.Infof("select multi cluster failed, err: %v", err)
			}
			if err := db.DB.Delete(&clusterMultiClusterRepository).Error; err != nil {
				logger.Log.Infof("delete multi cluster failed, err: %v", err)
			}
		}

		var clusterSyncLogs []MultiClusterSyncClusterLog
		if err := db.DB.Where("cluster_id = ?", c.ID).Find(&clusterSyncLogs).Error; err != nil {
			logger.Log.Infof("delete multi cluster sync logs failed, err: %v", err)
		}
		for _, clusterLog := range clusterSyncLogs {
			if clusterLog.ID != "" {
				if err := db.DB.Delete(&clusterLog).Error; err != nil {
					logger.Log.Infof("delete multi cluster logs failed, err: %v", err)
				}
			}
			var clusterResourceSyncLogs []MultiClusterSyncClusterResourceLog
			if err := db.DB.Where("multi_cluster_sync_cluster_log_id = ?", clusterLog.ID).Find(&clusterResourceSyncLogs).Error; err != nil {
				logger.Log.Infof("select multi cluster resource sync logs failed, err: %v", err)
			}
			for _, resourceLog := range clusterResourceSyncLogs {
				if err := db.DB.Delete(&resourceLog).Error; err != nil {
					logger.Log.Infof("select multi cluster resource sync logs failed, err: %v", err)
				}
			}
		}
	}
}

func (c Cluster) PrepareComponent(ingressType, ingressVersion, dnsCache, supportGpu string) []ClusterSpecComponent {
	var components []ClusterSpecComponent
	if ingressType == "traefik" {
		components = append(components, ClusterSpecComponent{
			ClusterID: c.ID,
			Name:      "traefik",
			Type:      "Ingress Controller",
			Version:   ingressVersion,
			Status:    constant.StatusEnabled,
		})
	}
	if ingressType == "nginx" {
		components = append(components, ClusterSpecComponent{
			ClusterID: c.ID,
			Name:      "ingress-nginx",
			Type:      "Ingress Controller",
			Version:   ingressVersion,
			Status:    constant.StatusEnabled,
		})
	}
	if dnsCache == constant.StatusEnabled {
		components = append(components, ClusterSpecComponent{
			ClusterID: c.ID,
			Name:      "dns-cache",
			Type:      "Dns Cache",
			Version:   "1.17.0",
			Status:    constant.StatusEnabled,
		})
	}
	if supportGpu == constant.StatusEnabled {
		components = append(components, ClusterSpecComponent{
			ClusterID: c.ID,
			Name:      "gpu",
			Type:      "GPU",
			Version:   "v1.7.0",
			Status:    constant.StatusEnabled,
		})
	}
	components = append(components, ClusterSpecComponent{
		ClusterID: c.ID,
		Name:      "metrics-server",
		Type:      "Metrics Server",
		Version:   "v0.5.0",
		Status:    constant.StatusEnabled,
	})
	return components
}

func (c Cluster) PrepareTools() []ClusterTool {
	return []ClusterTool{
		{
			ClusterID:    c.ID,
			Name:         "gatekeeper",
			Version:      "v3.7.0",
			Describe:     "OPA GateKeeper|OPA GateKeeper",
			Status:       constant.StatusWaiting,
			Logo:         "gatekeeper.jpg",
			Frame:        false,
			Url:          "",
			ProxyType:    "",
			Architecture: supportedArchitectureAll,
		},
		{
			ClusterID:    c.ID,
			Name:         "kubeapps",
			Version:      "2.4.2",
			Describe:     "应用商店|App store",
			Status:       constant.StatusWaiting,
			Logo:         "kubeapps.png",
			Frame:        true,
			Url:          "/proxy/kubeapps/{cluster_name}/root",
			Architecture: supportedArchitectureAmd64,
		},
		{
			ClusterID:    c.ID,
			Name:         "prometheus",
			Version:      "2.31.1",
			Describe:     "监控|Monitor",
			Status:       constant.StatusWaiting,
			Logo:         "prometheus.png",
			Frame:        true,
			Url:          "",
			ProxyType:    "nodeport",
			Architecture: supportedArchitectureAll,
		},
		{
			ClusterID:    c.ID,
			Name:         "logging",
			Version:      "v7.6.2",
			Describe:     "日志|Logs",
			Status:       constant.StatusWaiting,
			Logo:         "elasticsearch.png",
			Frame:        false,
			Url:          "/proxy/logging/{cluster_name}/root",
			Architecture: supportedArchitectureAmd64,
		},
		{
			ClusterID:    c.ID,
			Name:         "loki",
			Version:      "v2.1.0",
			Describe:     "日志|Logs",
			Status:       constant.StatusWaiting,
			Logo:         "loki.png",
			Frame:        false,
			Url:          "/proxy/loki/{cluster_name}/root",
			Architecture: supportedArchitectureAll,
		},
		{
			ClusterID:    c.ID,
			Name:         "grafana",
			Version:      "8.3.1",
			Describe:     "监控|Monitor",
			Status:       constant.StatusWaiting,
			Logo:         "grafana.png",
			Frame:        true,
			Url:          "/proxy/grafana/{cluster_name}",
			Architecture: supportedArchitectureAll,
		},
		{
			ClusterID:    c.ID,
			Name:         "chartmuseum",
			Version:      "v0.12.0",
			Describe:     "Chart 仓库|Chart warehouse",
			Status:       constant.StatusWaiting,
			Logo:         "chartmuseum.png",
			Frame:        false,
			Url:          "",
			Architecture: supportedArchitectureAll,
		},
		{
			ClusterID:    c.ID,
			Name:         "registry",
			Version:      "v2.7.1",
			Describe:     "镜像仓库|Image warehouse",
			Status:       constant.StatusWaiting,
			Logo:         "registry.png",
			Frame:        false,
			Url:          "",
			Architecture: supportedArchitectureAll,
		},
	}
}

func (c Cluster) GetKobeVars() map[string]string {
	result := map[string]string{}
	if c.Version != "" {
		index := strings.Index(c.Version, "-")
		result[facts.KubeVersionFactName] = c.Version[:index]
	}
	if c.NodeNameRule != "" {
		if c.NodeNameRule == constant.NodeNameRuleIP {
			result[facts.NodeNameRuleFactName] = c.NodeNameRule
		} else {
			result[facts.NodeNameRuleFactName] = constant.NodeNameRuleHostName
		}
	}
	c.loadRuntimeVars(result)
	c.loadConfVars(result)
	c.loadNetworkVars(result)
	c.loadComponentVars(result)

	return result
}

func (c Cluster) ParseInventory() *api.Inventory {
	var masters []string
	var workers []string
	var chrony []string
	var hosts []*api.Host
	var lbhosts []string

	i := 0
	for _, node := range c.Nodes {
		switch node.Role {
		case constant.NodeRoleNameMaster:
			if node.Status == "" || node.Status == constant.StatusRunning {
				masters = append(masters, node.Name)
			}
		case constant.NodeRoleNameWorker:
			if node.Status == "" || node.Status == constant.StatusRunning {
				workers = append(workers, node.Name)
			}
		}
		if c.SpecConf.LbMode == "external" {
			if node.Role == constant.NodeRoleNameMaster {
				lbhosts = append(lbhosts, node.Name)
			}
			if i == 0 {
				hosts = append(hosts, node.ToKobeHost(c.NodeNameRule, "master"))
				i = 1
			} else {
				hosts = append(hosts, node.ToKobeHost(c.NodeNameRule, "backup"))
			}
		} else {
			hosts = append(hosts, node.ToKobeHost(c.NodeNameRule, "internal"))
		}
	}
	if len(masters) > 0 {
		chrony = append(chrony, masters[0])
	}
	return &api.Inventory{
		Hosts: hosts,
		Groups: []*api.Group{
			{
				Name:     "kube-master",
				Hosts:    masters,
				Children: []string{},
				Vars:     map[string]string{},
			},
			{
				Name:  "kube-worker",
				Hosts: workers,
				Children: []string{
					"kube-master",
				},
				Vars: map[string]string{},
			},
			{
				Name:  "new-worker",
				Hosts: []string{},
				Vars:  map[string]string{},
			}, {
				Name:  "ex_lb",
				Hosts: lbhosts,
				Vars:  map[string]string{},
			},
			{
				Name:     "etcd",
				Hosts:    masters,
				Children: []string{"kube-master"},
				Vars:     map[string]string{},
			}, {
				Name:     "chrony",
				Hosts:    chrony,
				Children: []string{},
				Vars:     map[string]string{},
			},
			{
				Name:     "del-worker",
				Hosts:    []string{},
				Children: []string{},
				Vars:     map[string]string{},
			},
		},
	}
}

func (c Cluster) loadNetworkVars(result map[string]string) {
	if c.SpecNetwork.NetworkType != "" {
		result[facts.NetworkPluginFactName] = c.SpecNetwork.NetworkType
	}
	if c.SpecNetwork.CiliumNativeRoutingCidr != "" {
		result[facts.CiliumNativeRoutingCidrFactName] = c.SpecNetwork.CiliumNativeRoutingCidr
	}
	if c.SpecNetwork.CiliumTunnelMode != "" {
		result[facts.CiliumTunnelModeFactName] = c.SpecNetwork.CiliumTunnelMode
	}
	if c.SpecNetwork.CiliumVersion != "" {
		result[facts.CiliumVersionFactName] = c.SpecNetwork.CiliumVersion
	}
	if c.SpecNetwork.FlannelBackend != "" {
		result[facts.FlannelBackendFactName] = c.SpecNetwork.FlannelBackend
	}
	if c.SpecNetwork.CalicoIpv4PoolIpip != "" {
		result[facts.CalicoIpv4poolIpIpFactName] = c.SpecNetwork.CalicoIpv4PoolIpip
	}
}

func (c Cluster) loadRuntimeVars(result map[string]string) {
	if c.SpecRuntime.RuntimeType != "" {
		result[facts.ContainerRuntimeFactName] = c.SpecRuntime.RuntimeType
	}
	if c.SpecRuntime.DockerMirrorRegistry != "" {
		result[facts.DockerMirrorRegistryFactName] = c.SpecRuntime.DockerMirrorRegistry
	}
	if c.SpecRuntime.DockerRemoteApi != "" {
		result[facts.DockerRemoteApiFactName] = c.SpecRuntime.DockerRemoteApi
	}
	if c.SpecRuntime.DockerStorageDir != "" {
		result[facts.DockerStorageDirFactName] = c.SpecRuntime.DockerStorageDir
	}
	if c.SpecRuntime.ContainerdStorageDir != "" {
		result[facts.ContainerdStorageDirFactName] = c.SpecRuntime.ContainerdStorageDir
	}
	if c.SpecRuntime.DockerSubnet != "" {
		result[facts.DockerSubnetFactName] = c.SpecRuntime.DockerSubnet
	}
	if c.SpecRuntime.HelmVersion != "" {
		result[facts.HelmVersionFactName] = c.SpecRuntime.HelmVersion
	}
}

func (c Cluster) loadConfVars(result map[string]string) {
	if c.SpecConf.YumOperate != "" {
		result[facts.YumRepoFactName] = c.SpecConf.YumOperate
	}

	if c.SpecConf.KubeMaxPods != 0 {
		result[facts.KubeMaxPodsFactName] = strconv.Itoa(c.SpecConf.KubeMaxPods)
	}
	if c.SpecConf.KubeNetworkNodePrefix != 0 {
		result[facts.KubeNetworkNodePrefixFactName] = fmt.Sprint(c.SpecConf.KubeNetworkNodePrefix)
	}
	if c.SpecConf.KubePodSubnet != "" {
		result[facts.KubePodSubnetFactName] = c.SpecConf.KubePodSubnet
	}
	if c.SpecConf.KubeServiceSubnet != "" {
		result[facts.KubeServiceSubnetFactName] = c.SpecConf.KubeServiceSubnet
	}
	if c.SpecConf.CgroupDriver != "" {
		result[facts.CgroupDriverFactName] = c.SpecConf.CgroupDriver
	}
	if c.SpecConf.KubernetesAudit != "" {
		result[facts.KubernetesAuditFactName] = c.SpecConf.KubernetesAudit
	}
	if c.SpecConf.NodeportAddress != "" {
		result[facts.NodeportAddressFactName] = c.SpecConf.NodeportAddress
	}
	if c.SpecConf.KubeServiceNodePortRange != "" {
		result[facts.KubeServiceNodePortRangeFactName] = c.SpecConf.KubeServiceNodePortRange
	}

	if c.SpecConf.KubeProxyMode != "" {
		result[facts.KubeProxyModeFactName] = c.SpecConf.KubeProxyMode
	}
	if c.SpecConf.KubeDnsDomain != "" {
		result[facts.KubeDnsDomainFactName] = c.SpecConf.KubeDnsDomain
	}

	if c.SpecConf.MasterScheduleType != "" {
		result[facts.MasterScheduleTypeFactName] = c.SpecConf.MasterScheduleType
	}
	if c.SpecConf.LbMode != "" {
		result[facts.LbModeFactName] = c.SpecConf.LbMode
	}
	if c.SpecConf.LbKubeApiserverIp != "" {
		result[facts.LbKubeApiserverIpFactName] = c.SpecConf.LbKubeApiserverIp
	}
	if c.SpecConf.KubeApiServerPort != 0 {
		result[facts.KubeApiserverPortFactName] = fmt.Sprint(c.SpecConf.KubeApiServerPort)
	}
	result[facts.EtcdDataDirFactName] = c.SpecConf.EtcdDataDir
	result[facts.EtcdSnapshotCountFactName] = strconv.Itoa(c.SpecConf.EtcdSnapshotCount)
	result[facts.EtcdCompactionRetentionFactName] = strconv.Itoa(c.SpecConf.EtcdCompactionRetention)
	result[facts.EtcdMaxRequestFactName] = strconv.Itoa(c.SpecConf.EtcdMaxRequest * 1048576)
	result[facts.EtcdQuotaBackendFactName] = strconv.Itoa(c.SpecConf.EtcdQuotaBackend * 1073741824)
}

func (c Cluster) loadComponentVars(result map[string]string) {
	for _, c := range c.SpecComponent {
		switch c.Name {
		case "gpu":
			if c.Status == constant.StatusEnabled {
				result[facts.SupportGpuFactName] = constant.StatusEnabled
			}
		case "ingress-nginx":
			if c.Status == constant.StatusEnabled {
				result[facts.IngressControllerTypeFactName] = "nginx"
			}
		case "traefik":
			if c.Status == constant.StatusEnabled {
				result[facts.IngressControllerTypeFactName] = "traefik"
			}
		case "dns-cache":
			if c.Status == constant.StatusEnabled {
				result[facts.EnableDnsCacheFactName] = constant.StatusEnabled
			}
		}
	}
}
