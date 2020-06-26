package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type CloudProviderRepository interface {
	List() ([]model.CloudProvider, error)
}

func NewCloudProviderRepository() CloudProviderRepository {
	return &cloudProviderRepository{}
}

type cloudProviderRepository struct {
}

func (c cloudProviderRepository) List() ([]model.CloudProvider, error) {
	var cloudProviders []model.CloudProvider
	err := db.DB.Model(model.CloudProvider{}).Find(&cloudProviders).Error
	return cloudProviders, err
}
