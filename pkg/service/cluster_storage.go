package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service/storage"
)

type ClusterStorageService interface {
	CreateStorageClass(storageClass dto.StorageClass) error
}

type clusterStorageService struct{}

func (c clusterStorageService) CreateStorageClass(name string, storageClass dto.StorageClass) error {
	_, err := storage.NewStorageClassCreation(dto.ClusterWithEndpoint{}, storageClass)
	if err != nil {
		return err
	}
	return nil
}
