package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
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

	SecretID  string `json:"-"`
	StatusID  string `json:"-"`
	PlanID    string `json:"-"`
	LogId     string `json:"logId"`
	ProjectID string `json:"projectID"`
	Dirty     bool   `json:"dirty"`
	Plan      Plan   `json:"-"`

	SpecConf                 ClusterSpecConf          `gorm:"save_associations:false" json:"-"`
	SpecRuntime              ClusterSpecRuntime       `gorm:"save_associations:false" json:"-"`
	SpecNetwork              ClusterSpecNetwork       `gorm:"save_associations:false" json:"-"`
	SpecComponent            []ClusterSpecComponent   `gorm:"save_associations:false" json:"-"`
	Secret                   ClusterSecret            `gorm:"save_associations:false" json:"-"`
	Status                   ClusterStatus            `gorm:"save_associations:false" json:"-"`
	Nodes                    []ClusterNode            `gorm:"save_associations:false" json:"-"`
	Tools                    []ClusterTool            `gorm:"save_associations:false" json:"-"`
	Istios                   []ClusterIstio           `gorm:"save_associations:false" json:"-"`
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
		Preload("Status").
		Preload("SpecConf").
		Preload("SpecRuntime").
		Preload("SpecNetwork").
		Preload("Nodes").
		Preload("Nodes.Host").
		Preload("Tools").
		Preload("Istios").
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
	if err := tx.Where("cluster_id = ?", cluster.ID).Delete(&ClusterGpu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if cluster.StatusID != "" {
		if err := tx.Delete(&ClusterStatus{ID: cluster.StatusID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where("node_cluster_id = ?", cluster.ID).Delete(&ClusterStatus{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if cluster.SecretID != "" {
		if err := tx.Delete(&ClusterSecret{ID: cluster.SecretID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where("cluster = ?", cluster.Name).Delete(&KubepiBind{}).Error; err != nil {
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

	if len(cluster.Tools) > 0 {
		if err := tx.Where("cluster_id = ?", c.ID).Delete(&ClusterTool{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(cluster.Istios) > 0 {
		if err := tx.Where("cluster_id = ?", c.ID).Delete(&ClusterIstio{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	var cisTasks []CisTask
	if err := tx.Where("cluster_id = ?", c.ID).Find(&cisTasks).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(cisTasks) > 0 {
		for _, task := range cisTasks {
			if err := tx.Delete(&task).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Where("cluster_id = ?", c.ID).Delete(&ClusterStorageProvisioner{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("resource_id = ?", c.ID).Delete(&ProjectResource{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("cluster_id = ?", c.ID).Delete(&ClusterBackupStrategy{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("cluster_id = ?", c.ID).Delete(&ClusterBackupFile{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	var (
		messages   []Message
		messageIDs []string
	)
	if err := tx.Where("cluster_id = ? AND type != ?", c.ID, constant.System).Find(&messages).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, m := range messages {
		messageIDs = append(messageIDs, m.ID)
	}
	if len(messageIDs) > 0 {
		if err := tx.Where("message_id in (?)", messageIDs).Delete(&UserMessage{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where("cluster_id = ? AND type != ?", c.ID, constant.System).Delete(&Message{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(cluster.MultiClusterRepositories) > 0 {
		for _, repo := range cluster.MultiClusterRepositories {
			var clusterMultiClusterRepository ClusterMultiClusterRepository
			if err := tx.Where("cluster_id = ? AND multi_cluster_repository_id  = ?", c.ID, repo.ID).First(&clusterMultiClusterRepository).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Delete(&clusterMultiClusterRepository).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		var clusterSyncLogs []MultiClusterSyncClusterLog
		if err := tx.Where("cluster_id = ?", c.ID).Find(&clusterSyncLogs).Error; err != nil {
			tx.Rollback()
			return err
		}
		for _, clusterLog := range clusterSyncLogs {
			if clusterLog.ID != "" {
				if err := tx.Delete(&clusterLog).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
			var clusterResourceSyncLogs []MultiClusterSyncClusterResourceLog
			if err := tx.Where("multi_cluster_sync_cluster_log_id = ?", clusterLog.ID).Find(&clusterResourceSyncLogs).Error; err != nil {
				return err
			}
			for _, resourceLog := range clusterResourceSyncLogs {
				if err := tx.Delete(&resourceLog).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	if err := tx.Where("cluster_id = ?", cluster.ID).
		Delete(&ClusterResource{}).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (c Cluster) PrepareIstios() []ClusterIstio {
	return []ClusterIstio{
		{
			ClusterID: c.ID,
			Name:      "base",
			Version:   "v1.11.8",
			Describe:  "",
			Status:    constant.ClusterWaiting,
		},
		{
			ClusterID: c.ID,
			Name:      "pilot",
			Version:   "v1.11.8",
			Describe:  "",
			Status:    constant.ClusterWaiting,
		},
		{
			ClusterID: c.ID,
			Name:      "ingress",
			Version:   "v1.11.8",
			Describe:  "",
			Status:    constant.ClusterWaiting,
		},
		{
			ClusterID: c.ID,
			Name:      "egress",
			Version:   "v1.11.8",
			Describe:  "",
			Status:    constant.ClusterWaiting,
		},
	}
}

func (c Cluster) PrepareComponent(enableGpu, enableDnsCache, ingressController, ingressVersion string) []ClusterSpecComponent {
	var components []ClusterSpecComponent
	if enableGpu == constant.StatusEnabled {
		components = append(components, ClusterSpecComponent{
			ClusterID: c.ID,
			Name:      "gpu",
			Type:      "GPU",
			Version:   "v1.7.0",
			Status:    constant.StatusEnabled,
		})
	}
	if enableGpu == constant.StatusEnabled {
		components = append(components, ClusterSpecComponent{
			ClusterID: c.ID,
			Name:      "dns cache",
			Type:      "DNS_CACHE",
			Version:   "1.17.0",
			Status:    constant.StatusEnabled,
		})
	}
	if len(ingressController) != 0 {
		components = append(components, ClusterSpecComponent{
			ClusterID: c.ID,
			Name:      ingressController,
			Type:      "INGRESS_CONTROLLER",
			Version:   ingressVersion,
			Status:    constant.StatusEnabled,
		})
	}
	return components
}

func (c Cluster) PrepareTools() []ClusterTool {
	return []ClusterTool{
		{
			ClusterID:    c.ID,
			Name:         "gatekeeper",
			Version:      "v3.7.0",
			Describe:     "OPA GateKeeper|OPA GateKeeper",
			Status:       constant.ClusterWaiting,
			Logo:         "gatekeeper.jpg",
			Frame:        false,
			Url:          "",
			ProxyType:    "",
			Architecture: supportedArchitectureAll,
		},
		{
			ClusterID:    c.ID,
			Name:         "kubeapps",
			Version:      "v1.10.2",
			Describe:     "应用商店|App store",
			Status:       constant.ClusterWaiting,
			Logo:         "kubeapps.png",
			Frame:        true,
			Url:          "/proxy/kubeapps/{cluster_name}/root",
			Architecture: supportedArchitectureAmd64,
		},
		{
			ClusterID:    c.ID,
			Name:         "prometheus",
			Version:      "v2.18.1",
			Describe:     "监控|Monitor",
			Status:       constant.ClusterWaiting,
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
			Status:       constant.ClusterWaiting,
			Logo:         "elasticsearch.png",
			Frame:        false,
			Url:          "/proxy/logging/{cluster_name}/root",
			Architecture: supportedArchitectureAmd64,
		},
		{
			ClusterID:    c.ID,
			Name:         "loki",
			Version:      "v2.0.0",
			Describe:     "日志|Logs",
			Status:       constant.ClusterWaiting,
			Logo:         "loki.png",
			Frame:        false,
			Url:          "/proxy/loki/{cluster_name}/root",
			Architecture: supportedArchitectureAll,
		},
		{
			ClusterID:    c.ID,
			Name:         "grafana",
			Version:      "v7.3.3",
			Describe:     "监控|Monitor",
			Status:       constant.ClusterWaiting,
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
			Status:       constant.ClusterWaiting,
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
			Status:       constant.ClusterWaiting,
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
			if node.Status == "" || node.Status == constant.ClusterRunning {
				masters = append(masters, node.Name)
			}
		case constant.NodeRoleNameWorker:
			if node.Status == "" || node.Status == constant.ClusterRunning {
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
	if c.SpecConf.KubernetesAudit != "" {
		result[facts.KubernetesAuditFactName] = c.SpecConf.KubernetesAudit
	}
	if c.SpecConf.NodeportAddress != "" {
		result[facts.NodeportAddressFactName] = c.SpecConf.NodeportAddress
	}
	if c.SpecConf.KubeServiceNodePortRange != "" {
		result[facts.KubeServiceNodePortRangeFactName] = c.SpecConf.KubeServiceNodePortRange
	}
	if c.SpecConf.DnsCacheVersion != "" {
		result[facts.DnsCacheVersionFactName] = c.SpecConf.DnsCacheVersion
	}
	if c.SpecConf.EnableDnsCache != "" {
		result[facts.EnableDnsCacheFactName] = c.SpecConf.EnableDnsCache
	}
	if c.SpecConf.IngressControllerType != "" {
		result[facts.IngressControllerTypeFactName] = c.SpecConf.IngressControllerType
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
}
