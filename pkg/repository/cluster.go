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
		Where(&model.Cluster{Name: name}).
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
	db.DB.Model(&model.Cluster{})
	if err := db.DB.Model(&model.Cluster{}).
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
	var total int
	var clusters []model.Cluster
	var project model.Project
	err := db.DB.Model(&model.Project{}).Where(&model.Project{Name: projectName}).First(&project).Error
	if err != nil {
		return 0, nil, err
	}
	var projectResources []model.ProjectResource
	err = db.DB.Model(&model.ProjectResource{}).Where(&model.ProjectResource{ProjectID: project.ID, ResourceType: constant.ResourceCluster}).Find(&projectResources).Error
	if err != nil {
		return 0, nil, err
	}
	var resourceIds []string
	for _, pr := range projectResources {
		resourceIds = append(resourceIds, pr.ResourceID)
	}

	if err := db.DB.Model(&model.Cluster{}).
		Offset((num-1)*size).
		Limit(size).
		Where("id in (?)", resourceIds).
		Preload("Status").
		Preload("Spec").
		Preload("Nodes").
		Preload("MultiClusterRepositories").
		Count(&total).
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
	if err := db.DB.Where(&model.Cluster{Name: name}).First(&cluster).Error; err != nil {
		return err
	}
	var prometheus model.ClusterTool
	err := db.DB.Where(&model.ClusterTool{Name: "prometheus", ClusterID: cluster.ID}).First(&prometheus).Error
	if err != nil {
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
