package model

import (
	"strconv"

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
	Plan                     Plan                     `json:"-"`
	Spec                     ClusterSpec              `gorm:"save_associations:false" json:"spec"`
	Secret                   ClusterSecret            `gorm:"save_associations:false" json:"-"`
	Status                   ClusterStatus            `gorm:"save_associations:false" json:"-"`
	Nodes                    []ClusterNode            `gorm:"save_associations:false" json:"-"`
	Tools                    []ClusterTool            `gorm:"save_associations:false" json:"-"`
	MultiClusterRepositories []MultiClusterRepository `gorm:"many2many:cluster_multi_cluster_repository"`
}

func (c *Cluster) BeforeCreate() error {
	c.ID = uuid.NewV4().String()
	tx := db.DB.Begin()
	if err := tx.Create(&c.Spec).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&c.Status).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&c.Secret).Error; err != nil {
		tx.Rollback()
		return err
	}
	c.SpecID = c.Spec.ID
	c.StatusID = c.Status.ID
	c.SecretID = c.Secret.ID
	for i := range c.Nodes {
		c.Nodes[i].ClusterID = c.ID
		if err := tx.Create(&c.Nodes[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
		if c.Nodes[i].Host.ID != "" {
			c.Nodes[i].Host.ClusterID = c.ID
			err := tx.Save(&c.Nodes[i].Host).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	for _, tool := range c.PrepareTools() {
		tool.ClusterID = c.ID
		err := tx.Create(&tool).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
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
		Preload("Tools").
		Preload("MultiClusterRepositories").
		First(&cluster).Error; err != nil {
		return err
	}
	if cluster.SpecID != "" {
		if err := tx.Delete(ClusterSpec{ID: cluster.SpecID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if cluster.StatusID != "" {
		if err := tx.Delete(ClusterStatus{ID: cluster.StatusID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if cluster.SecretID != "" {
		if err := tx.Delete(ClusterSecret{ID: cluster.SecretID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(cluster.Nodes) > 0 {
		for _, node := range cluster.Nodes {
			if err := tx.Where(ClusterNode{ID: node.ID}).
				Delete(ClusterNode{}).Error; err != nil {
				tx.Rollback()
				return err
			}
			if node.HostID != "" {
				host := Host{ID: node.HostID}
				if err := tx.First(&host).Error; err != nil {
					tx.Rollback()
					return err
				}
				if cluster.Spec.Provider == constant.ClusterProviderBareMetal {
					host.ClusterID = ""
					if err := tx.Save(&host).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
				if cluster.Spec.Provider == constant.ClusterProviderPlan {
					host.ClusterID = ""
					if err := tx.Save(&host).Error; err != nil {
						tx.Rollback()
						return err
					}
					var projectResources []ProjectResource
					if err := tx.Where(ProjectResource{ResourceId: host.ID, ResourceType: constant.ResourceHost}).Find(&projectResources).Error; err != nil {
						return err
					}
					if len(projectResources) > 0 {
						for _, p := range projectResources {
							tx.Delete(&p)
						}
					}
					if err := tx.Delete(&host).Error; err != nil {
						tx.Rollback()
						return err
					}

				}
			}
		}
	}
	if cluster.Spec.Provider == constant.ClusterProviderBareMetal {
		var hosts []Host
		tx.Where(Host{ClusterID: c.ID}).Find(&hosts)
		if len(hosts) > 0 {
			for i := range hosts {
				hosts[i].ClusterID = ""
				tx.Save(&hosts[i])
			}
		}

	}
	if len(cluster.Tools) > 0 {
		for _, tool := range cluster.Tools {
			if tool.ID != "" {
				if err := tx.Delete(&tool).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	var cisTasks []CisTask
	db.DB.Where(CisTask{ClusterID: c.ID}).Find(&cisTasks)
	if len(cisTasks) > 0 {
		for _, task := range cisTasks {
			if err := tx.Delete(&task).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	var storageProvisioners []ClusterStorageProvisioner
	db.DB.Where(ClusterStorageProvisioner{ClusterID: c.ID}).Find(&storageProvisioners)
	if len(storageProvisioners) > 0 {
		for _, p := range storageProvisioners {
			if err := tx.Delete(&p).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	var projectResource ProjectResource
	tx.Where(ProjectResource{ResourceId: c.ID, ResourceType: constant.ResourceCluster}).First(&projectResource)
	if projectResource.ID != "" {
		if err := tx.Delete(&projectResource).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	var clusterBackupStrategy ClusterBackupStrategy
	tx.Where(ClusterBackupStrategy{ClusterID: c.ID}).First(&clusterBackupStrategy)
	if clusterBackupStrategy.ID != "" {
		if err := tx.Delete(&clusterBackupStrategy).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	var clusterBackupFiles []ClusterBackupFile
	tx.Where(ClusterBackupFile{ClusterID: c.ID}).Find(&clusterBackupFiles)
	if len(clusterBackupFiles) > 0 {
		for _, c := range clusterBackupFiles {
			if err := tx.Delete(&c).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	var messages []Message
	var messageIds []string
	tx.Where(Message{ClusterID: c.ID}).Find(&messages)
	if len(messages) > 0 {
		for _, m := range messages {
			messageIds = append(messageIds, m.ID)
			if err := tx.Delete(&m).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if len(messageIds) > 0 {
		var userMessages []UserMessage
		if err := tx.Where("message_id in (?)", messageIds).Delete(&userMessages).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(cluster.MultiClusterRepositories) > 0 {
		for _, repo := range cluster.MultiClusterRepositories {
			var clusterMultiClusterRepository ClusterMultiClusterRepository
			if err := tx.Where(ClusterMultiClusterRepository{ClusterID: c.ID, MultiClusterRepositoryID: repo.ID}).First(&clusterMultiClusterRepository).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Delete(&clusterMultiClusterRepository).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		var clusterSyncLogs []MultiClusterSyncClusterLog
		if err := tx.Where(MultiClusterSyncClusterLog{ClusterID: c.ID}).Find(&clusterSyncLogs).Error; err != nil {
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
			if err := tx.Where(MultiClusterSyncClusterResourceLog{MultiClusterSyncClusterLogID: clusterLog.ID}).Find(&clusterResourceSyncLogs).Error; err != nil {
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
	tx.Commit()
	return nil
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
		result[facts.KubeVersionFactName] = c.Spec.Version
	}
	if c.Spec.NetworkType != "" {
		result[facts.NetworkPluginFactName] = c.Spec.NetworkType
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
	if c.Spec.LbKubeApiserverIp != "" {
		result[facts.LbKubeApiserverPortFactName] = c.Spec.LbKubeApiserverIp
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
	if c.Spec.Architectures != "" {
		result[facts.ArchitecturesFactName] = c.Spec.Architectures
	}
	if c.Spec.KubernetesAudit != "" {
		result[facts.KubernetesAuditFactName] = c.Spec.KubernetesAudit
	}
	if c.Spec.Architectures != "" {
		result[facts.DockerSubnetFactName] = c.Spec.DockerSubnet
	}
	if c.Spec.HelmVersion != "" {
		result[facts.HelmVersionFactName] = c.Spec.HelmVersion
	}
	if c.Spec.NetworkInterface != "" {
		result[facts.NetworkInterfaceFactName] = c.Spec.NetworkInterface
	}
	if c.Spec.SupportGpu != "" {
		result[facts.SupportGpuName] = c.Spec.SupportGpu
	}
	return result
}

func (c Cluster) ParseInventory() *api.Inventory {
	var masters []string
	var workers []string
	var chrony []string
	var hosts []*api.Host
	for _, node := range c.Nodes {
		hosts = append(hosts, node.ToKobeHost())
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

				Name:     "lb",
				Hosts:    []string{},
				Children: []string{},
				Vars:     map[string]string{},
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
