package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type ClusterToolService interface {
	List(clusterName string) ([]dto.ClusterTool, error)
	Save(clusterName string, item dto.ClusterTool) (dto.ClusterTool, error)
}

func NewClusterToolService() *clusterToolService {
	return &clusterToolService{
		toolRepo: repository.NewClusterToolRepository(),
	}
}

type clusterToolService struct {
	toolRepo repository.ClusterToolRepository
}

func (c clusterToolService) List(clusterName string) ([]dto.ClusterTool, error) {
	var items []dto.ClusterTool
	ms, err := c.toolRepo.List(clusterName)
	if err != nil {
		return items, err
	}
	for _, m := range ms {
		items = append(items, dto.ClusterTool{ClusterTool: m})
	}
	return items, nil
}

func (c clusterToolService) Save(clusterName string, item dto.ClusterTool) (dto.ClusterTool, error) {
	if err := c.toolRepo.Save(clusterName, &item.ClusterTool); err != nil {
		return item, err
	}
	return item, nil
}
