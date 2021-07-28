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
	ID                       string                   `json:"-"`
	Name                     string                   `json:"name" gorm:"not null;unique"`
	Source                   string                   `json:"source"`
	SpecID                   string                   `json:"-"`
	SecretID                 string                   `json:"-"`
	StatusID                 string                   `json:"-"`
	PlanID                   string                   `json:"-"`
	LogId                    string                   `json:"logId"`
	ProjectID                string                   `json:"projectID"`
	Dirty                    bool                     `json:"dirty"`
	Plan                     Plan                     `json:"-"`
	Spec                     ClusterSpec              `gorm:"save_associations:false" json:"spec"`
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
		Preload("Spec").
		Preload("Nodes").
		Preload("Nodes.Host").
		Preload("Tools").
		Preload("Istios").
		Preload("MultiClusterRepositories").
		First(&cluster).Error; err != nil {
		return err
	}
	if cluster.SpecID != "" {
		if err := tx.Delete(&ClusterSpec{ID: cluster.SpecID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if cluster.StatusID != "" {
		if err := tx.Delete(&ClusterStatus{ID: cluster.StatusID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if cluster.SecretID != "" {
		if err := tx.Delete(&ClusterSecret{ID: cluster.SecretID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	var (
		hostIDList []string
		hostIPList []string
	)
	for _, node := range cluster.Nodes {
		hostIDList = append(hostIDList, node.HostID)
		hostIPList = append(hostIPList, node.Host.Ip)
	}
	if cluster.Spec.Provider == constant.ClusterProviderPlan {
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
			Name:     "base",
			Version:  "v1.8.0",
			Describe: "",
			Status:   constant.ClusterWaiting,
		},
		{
			Name:     "pilot",
			Version:  "v1.8.0",
			Describe: "",
			Status:   constant.ClusterWaiting,
		},
		{
			Name:     "ingress",
			Version:  "v1.8.0",
			Describe: "",
			Status:   constant.ClusterWaiting,
		},
		{
			Name:     "egress",
			Version:  "v1.8.0",
			Describe: "",
			Status:   constant.ClusterWaiting,
		},
	}
}

func (c Cluster) PrepareTools() []ClusterTool {
	return []ClusterTool{
		{
			Name:         "dashboard",
			Version:      "v2.0.3",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "kubernetes.png",
			Frame:        true,
			Url:          "/proxy/dashboard/{cluster_name}/root",
			Architecture: supportedArchitectureAll,
		},
		{
			Name:         "kubeapps",
			Version:      "v1.10.2",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "kubeapps.png",
			Frame:        true,
			Url:          "/proxy/kubeapps/{cluster_name}/root",
			Architecture: supportedArchitectureAmd64,
		},
		{
			Name:         "prometheus",
			Version:      "v2.18.1",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "prometheus.png",
			Frame:        false,
			Url:          "",
			Architecture: supportedArchitectureAll,
		},
		{
			Name:         "logging",
			Version:      "v7.6.2",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "elasticsearch.png",
			Frame:        false,
			Url:          "/proxy/logging/{cluster_name}/root",
			Architecture: supportedArchitectureAmd64,
		},
		{
			Name:         "loki",
			Version:      "v2.0.0",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "loki.png",
			Frame:        false,
			Url:          "/proxy/loki/{cluster_name}/root",
			Architecture: supportedArchitectureAll,
		},
		{
			Name:         "grafana",
			Version:      "v7.3.3",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "grafana.png",
			Frame:        true,
			Url:          "/proxy/grafana/{cluster_name}",
			Architecture: supportedArchitectureAll,
		},
		{
			Name:         "chartmuseum",
			Version:      "v0.12.0",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "chartmuseum.png",
			Frame:        false,
			Url:          "",
			Architecture: supportedArchitectureAll,
		},
		{
			Name:         "registry",
			Version:      "v2.7.1",
			Describe:     "",
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
	if c.Spec.Version != "" {
		index := strings.Index(c.Spec.Version, "-")
		result[facts.KubeVersionFactName] = c.Spec.Version[:index]
	}
	if c.Spec.NetworkType != "" {
		result[facts.NetworkPluginFactName] = c.Spec.NetworkType
	}
	if c.Spec.CiliumNativeRoutingCidr != "" {
		result[facts.CiliumNativeRoutingCidrFactName] = c.Spec.CiliumNativeRoutingCidr
	}
	if c.Spec.CiliumTunnelMode != "" {
		result[facts.CiliumTunnelModeFactName] = c.Spec.CiliumTunnelMode
	}
	if c.Spec.CiliumVersion != "" {
		result[facts.CiliumVersionFactName] = c.Spec.CiliumVersion
	}
	if c.Spec.EnableDnsCache != "" {
		result[facts.EnableDnsCacheFactName] = c.Spec.EnableDnsCache
	}
	if c.Spec.DnsCacheVersion != "" {
		result[facts.DnsCacheVersionFactName] = c.Spec.DnsCacheVersion
	}
	if c.Spec.FlannelBackend != "" {
		result[facts.FlannelBackendFactName] = c.Spec.FlannelBackend
	}
	if c.Spec.CalicoIpv4poolIpip != "" {
		result[facts.CalicoIpv4poolIpIpFactName] = c.Spec.CalicoIpv4poolIpip
	}
	if c.Spec.RuntimeType != "" {
		result[facts.ContainerRuntimeFactName] = c.Spec.RuntimeType
	}
	if c.Spec.DockerStorageDir != "" {
		result[facts.DockerStorageDirFactName] = c.Spec.DockerStorageDir
	}
	if c.Spec.ContainerdStorageDir != "" {
		result[facts.ContainerdStorageDirFactName] = c.Spec.ContainerdStorageDir
	}
	if c.Spec.LbMode != "" {
		result[facts.LbModeFactName] = c.Spec.LbMode
	}
	if c.Spec.LbKubeApiserverIp != "" {
		result[facts.LbKubeApiserverIpFactName] = c.Spec.LbKubeApiserverIp
	}
	if c.Spec.LbKubeApiserverPort != "" {
		result[facts.LbKubeApiserverPortFactName] = c.Spec.LbKubeApiserverPort
	}
	if c.Spec.KubePodSubnet != "" {
		result[facts.KubePodSubnetFactName] = c.Spec.KubePodSubnet
	}
	if c.Spec.KubeServiceSubnet != "" {
		result[facts.KubeServiceSubnetFactName] = c.Spec.KubeServiceSubnet
	}
	if c.Spec.KubeMaxPods != 0 {
		result[facts.KubeMaxPodsFactName] = strconv.Itoa(c.Spec.KubeMaxPods)
	}
	if c.Spec.KubeProxyMode != "" {
		result[facts.KubeProxyModeFactName] = c.Spec.KubeProxyMode
	}
	if c.Spec.IngressControllerType != "" {
		result[facts.IngressControllerTypeFactName] = c.Spec.IngressControllerType
	}
	if c.Spec.KubernetesAudit != "" {
		result[facts.KubernetesAuditFactName] = c.Spec.KubernetesAudit
	}
	if c.Spec.DockerSubnet != "" {
		result[facts.DockerSubnetFactName] = c.Spec.DockerSubnet
	}
	if c.Spec.HelmVersion != "" {
		result[facts.HelmVersionFactName] = c.Spec.HelmVersion
	}
	if c.Spec.NetworkInterface != "" {
		result[facts.NetworkInterfaceFactName] = c.Spec.NetworkInterface
	}
	if c.Spec.NetworkCidr != "" {
		result[facts.NetworkCidrFactName] = c.Spec.NetworkCidr
	}
	if c.Spec.SupportGpu != "" {
		result[facts.SupportGpuName] = c.Spec.SupportGpu
	}
	if c.Spec.YumOperate != "" {
		result[facts.YumRepoFactName] = c.Spec.YumOperate
	}
	if c.Spec.KubeNetworkNodePrefix != 0 {
		result[facts.KubeNetworkNodePrefixFactName] = fmt.Sprint(c.Spec.KubeNetworkNodePrefix)
	}

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
		if c.Spec.LbMode == "external" {
			if node.Role == constant.NodeRoleNameMaster {
				lbhosts = append(lbhosts, node.Name)
			}
			if i == 0 {
				hosts = append(hosts, node.ToKobeHost("master"))
				i = 1
			} else {
				hosts = append(hosts, node.ToKobeHost("backup"))
			}
		} else {
			hosts = append(hosts, node.ToKobeHost("internal"))
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
