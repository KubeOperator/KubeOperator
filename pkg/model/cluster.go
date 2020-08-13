package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

type Cluster struct {
	common.BaseModel
	ID       string        `json:"-"`
	Name     string        `json:"name" gorm:"not null;unique"`
	Source   string        `json:"source"`
	SpecID   string        `json:"-"`
	SecretID string        `json:"-"`
	StatusID string        `json:"-"`
	PlanID   string        `json:"-"`
	Plan     Plan          `json:"-"`
	Spec     ClusterSpec   `gorm:"save_associations:false" json:"spec"`
	Secret   ClusterSecret `gorm:"save_associations:false" json:"-"`
	Status   ClusterStatus `gorm:"save_associations:false" json:"-"`
	Nodes    []ClusterNode `gorm:"save_associations:false" json:"-"`
	Tools    []ClusterTool `gorm:"save_associations:false" json:"-"`
}

func (c Cluster) TableName() string {
	return "ko_cluster"
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
	for i, _ := range c.Nodes {
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
					if err := db.DB.Where(ProjectResource{ResourceId: host.ID, ResourceType: constant.ResourceHost}).Find(&projectResources).Error; err != nil {
						return err
					}
					if len(projectResources) > 0 {
						for _, p := range projectResources {
							db.DB.Delete(&p)
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

	var projectResource ProjectResource
	if err := tx.Where(ProjectResource{ResourceId: c.ID, ResourceType: constant.ResourceCluster}).Delete(&projectResource).Error; err != nil {
		tx.Rollback()
		return err
	}

	var clusterBackupStrategy ClusterBackupStrategy
	if err := tx.Where(ClusterBackupStrategy{ClusterID: c.ID}).Delete(&clusterBackupStrategy).Error; err != nil {
		tx.Rollback()
		return err
	}

	var clusterBackupFiles []ClusterBackupFile
	if err := tx.Where(ClusterBackupFile{ClusterID: c.ID}).Delete(&clusterBackupFiles).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c Cluster) PrepareTools() []ClusterTool {
	return []ClusterTool{
		{
			Name:         "dashboard",
			Version:      "v1.0.0",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "kubernetes.png",
			Frame:        true,
			Url:          "/proxy/dashboard/{cluster_name}/root",
			Architecture: supportedArchitectureAll,
		},
		{
			Name:         "kubeapps",
			Version:      "v1.0.0",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "kubeapps.png",
			Frame:        true,
			Url:          "/proxy/kubeapps/{cluster_name}/root",
			Architecture: supportedArchitectureAmd64,
		},
		{
			Name:         "prometheus",
			Version:      "v1.0.0",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "prometheus.png",
			Frame:        false,
			Url:          "",
			Architecture: supportedArchitectureAll,
		},
		{
			Name:         "chartmuseum",
			Version:      "v1.0.0",
			Describe:     "",
			Status:       constant.ClusterWaiting,
			Logo:         "chartmuseum.png",
			Frame:        false,
			Url:          "",
			Architecture: supportedArchitectureAmd64,
		},
		{
			Name:         "registry",
			Version:      "v1.0.0",
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

	return result
}

func (c Cluster) ParseInventory() api.Inventory {
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
	return api.Inventory{
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
