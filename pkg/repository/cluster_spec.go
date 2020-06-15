package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterSpecRepository interface {
	Get(id string) (model.ClusterSpec, error)
	Save(spec *model.ClusterSpec) error
	Delete(id string) error
}

func NewClusterSpecRepository() ClusterSpecRepository {
	return &clusterSpecRepository{}
}

type clusterSpecRepository struct{}

func (c clusterSpecRepository) Get(id string) (model.ClusterSpec, error) {
	spec := model.ClusterSpec{
		ID: id,
	}
	if err := db.DB.First(&spec).Error; err != nil {
		return spec, err
	}
	return spec, nil
}

func (c clusterSpecRepository) Save(spec *model.ClusterSpec) error {
	if db.DB.NewRecord(spec) {
		if err := db.DB.Create(&spec).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(&spec).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c clusterSpecRepository) Delete(id string) error {
	spec := model.ClusterSpec{ID: id}
	if err := db.DB.First(&spec).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&spec).Error; err != nil {
		return err
	}
	return nil
}
