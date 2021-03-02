package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/grafana"
)

var log = logger.Default

type ClusterRepository interface {
	Get(name string) (model.Cluster, error)
	List() ([]model.Cluster, error)
	Save(cluster *model.Cluster) error
	Delete(name string) error
	Page(num, size int, projectName string) (int, []model.Cluster, error)
}

func NewClusterRepository() ClusterRepository {
	return &clusterRepository{}
}

type clusterRepository struct {
}

func (c clusterRepository) Get(name string) (model.Cluster, error) {
	var cluster model.Cluster
	if err := db.DB.
		Where("name = ?", name).
		Preload("Status").
		Preload("Spec").
		Preload("Nodes").
		Preload("Nodes.Host").
		Preload("Nodes.Host.Credential").
		Preload("Nodes.Host.Zone").
		Preload("MultiClusterRepositories").
		Find(&cluster).Error; err != nil {
		return cluster, err
	}
	return cluster, nil
}

func (c clusterRepository) List() ([]model.Cluster, error) {
	var clusters []model.Cluster
	if err := db.DB.
		Preload("Status").
		Preload("Spec").
		Preload("Nodes").
		Preload("Nodes.Host").
		Preload("Nodes.Host.Credential").
		Preload("Nodes.Host.Zone").
		Preload("MultiClusterRepositories").
		Find(&clusters).Error; err != nil {
		return clusters, err
	}
	return clusters, nil
}

func (c clusterRepository) Page(num, size int, projectName string) (int, []model.Cluster, error) {
	var (
		total    int
		clusters []model.Cluster
		project  model.Project
	)
	if err := db.DB.Where("name = ?", projectName).First(&project).Error; err != nil {
		return 0, nil, err
	}
	var projectResources []model.ProjectResource
	if err := db.DB.Where("project_id = ? AND resource_type = ?", project.ID, constant.ResourceCluster).Find(&projectResources).Error; err != nil {
		return 0, nil, err
	}
	var resourceIds []string
	for _, pr := range projectResources {
		resourceIds = append(resourceIds, pr.ResourceID)
	}

	if err := db.DB.Model(&model.Cluster{}).
		Where("id in (?)", resourceIds).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Preload("Status").
		Preload("Spec").
		Preload("Nodes").
		Preload("MultiClusterRepositories").
		Find(&clusters).Error; err != nil {
		return total, clusters, err
	}
	return total, clusters, nil
}

func (c clusterRepository) Save(cluster *model.Cluster) error {
	if db.DB.NewRecord(cluster) {
		if err := db.DB.Create(cluster).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(cluster).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c clusterRepository) Delete(name string) error {
	var cluster model.Cluster
	if err := db.DB.Where("name = ?", name).First(&cluster).Error; err != nil {
		return err
	}
	var prometheus model.ClusterTool
	if err := db.DB.Where("name = ? AND cluster_id = ?", "prometheus", cluster.ID).First(&prometheus).Error; err != nil {
		log.Error(err)
	}
	if prometheus.Status == constant.ClusterRunning {
		// 尝试删除 grafana
		gClient := grafana.NewClient()
		if err := gClient.DeleteDashboard(cluster.Name); err != nil {
			log.Error(err)
		}
		if err := gClient.DeleteDataSource(cluster.Name); err != nil {
			log.Error(err)
		}
	}
	if err := db.DB.Delete(&cluster).Error; err != nil {
		return err
	}

	return nil
}
