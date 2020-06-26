package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
)

type CloudProviderService interface {
	List() ([]dto.CloudProvider, error)
}

type cloudProviderService struct {
	cloudProviderRepo repository.CloudProviderRepository
}

func NewCloudProviderService() CloudProviderService {
	return &cloudProviderService{
		cloudProviderRepo: repository.NewCloudProviderRepository(),
	}
}

func (c cloudProviderService) List() ([]dto.CloudProvider, error) {
	var cloudProviderDTOs []dto.CloudProvider
	mos, err := c.cloudProviderRepo.List()
	if err != nil {
		return cloudProviderDTOs, err
	}
	for _, mo := range mos {
		cloudProviderDTOs = append(cloudProviderDTOs, dto.CloudProvider{CloudProvider: mo})
	}
	return cloudProviderDTOs, err
}
