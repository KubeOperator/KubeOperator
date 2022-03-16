package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterGpuRepository interface {
	Add(clusterID string) error
	GetByClusterName(clusterName string) (model.ClusterGpu, error)
	GetByClusterID(clusterID string) (model.ClusterGpu, error)
	Save(gpu *model.ClusterGpu) error
	Delete(clusterID string) error
}

func NewClusterGpuRepository() ClusterGpuRepository {
	return &clusterGpuRepository{}
}

type clusterGpuRepository struct{}

func (c clusterGpuRepository) Add(clusterID string) error {
	gpu := &model.ClusterGpu{
		ClusterID: clusterID,
		Status:    constant.StatusEnabled,
	}
	return db.DB.Create(&gpu).Error
}

func (c clusterGpuRepository) GetByClusterID(clusterID string) (model.ClusterGpu, error) {
	var gpu model.ClusterGpu
	if err := db.DB.Where("cluster_id = ?", clusterID).First(&gpu).Error; err != nil {
		return gpu, err
	}
	return gpu, nil
}

func (c clusterGpuRepository) GetByClusterName(clusterName string) (model.ClusterGpu, error) {
	var cluster model.Cluster
	var gpu model.ClusterGpu
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return gpu, err
	}
	if err := db.DB.Where("cluster_id = ?", cluster.ID).First(&gpu).Error; err != nil {
		return gpu, err
	}
	return gpu, nil
}

func (c clusterGpuRepository) Delete(clusterID string) error {
	return db.DB.Where("cluster_id = ?", clusterID).Delete(&model.ClusterGpu{}).Error
}

func (c clusterGpuRepository) Save(gpu *model.ClusterGpu) error {
	if db.DB.NewRecord(gpu) {
		if err := db.DB.Create(gpu).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(gpu).Error; err != nil {
			return nil
		}
	}
	return nil
}
