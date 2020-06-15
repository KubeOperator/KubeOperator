package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ClusterSecretRepository interface {
	Get(id string) (model.ClusterSecret, error)
	Save(status *model.ClusterSecret) error
	Delete(id string) error
}

func NewClusterSecretRepository() ClusterSecretRepository {
	return &clusterSecretRepository{}
}

type clusterSecretRepository struct {
}

func (c clusterSecretRepository) Get(id string) (model.ClusterSecret, error) {
	status := model.ClusterSecret{
		ID: id,
	}
	if err := db.DB.First(&status).Error; err != nil {
		return status, err
	}
	return status, nil
}

func (c clusterSecretRepository) Save(status *model.ClusterSecret) error {
	if db.DB.NewRecord(status) {
		if err := db.DB.Create(&status).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(&status).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c clusterSecretRepository) Delete(id string) error {
	secret := model.ClusterSecret{ID: id}
	if err := db.DB.First(&secret).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&secret).Error; err != nil {
		return err
	}
	return nil
}
