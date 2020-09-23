package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type ClusterEventService interface {
	List(clusterName string) ([]dto.ClusterEventDTO, error)
	ListLimitOneDay(clusterName string) ([]dto.ClusterEventDTO, error)
	ExistEventUid(uid, clusterId string) (bool, error)
	Save(event model.ClusterEvent) error
}

type clusterEventService struct {
	clusterEventRepo repository.ClusterEventRepository
	clusterRepo      repository.ClusterRepository
}

func NewClusterEventService() ClusterEventService {
	return &clusterEventService{
		clusterEventRepo: repository.NewClusterEventRepository(),
		clusterRepo:      repository.NewClusterRepository(),
	}
}

func (c clusterEventService) List(clusterName string) ([]dto.ClusterEventDTO, error) {
	cluster, err := c.clusterRepo.Get(clusterName)
	if err != nil {
		return nil, err
	}
	var eventDTOs []dto.ClusterEventDTO
	events, err := c.clusterEventRepo.List(cluster.ID)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		eventDTOs = append(eventDTOs, dto.ClusterEventDTO{
			ClusterEvent: event,
		})
	}
	return eventDTOs, nil
}

func (c clusterEventService) ListLimitOneDay(clusterName string) ([]dto.ClusterEventDTO, error) {
	cluster, err := c.clusterRepo.Get(clusterName)
	if err != nil {
		return nil, err
	}
	var eventDTOs []dto.ClusterEventDTO
	events, err := c.clusterEventRepo.ListLimitOneDay(cluster.ID)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		eventDTOs = append(eventDTOs, dto.ClusterEventDTO{
			ClusterEvent: event,
		})
	}
	return eventDTOs, nil
}

func (c clusterEventService) ExistEventUid(uid, clusterId string) (bool, error) {
	events, err := c.clusterEventRepo.ListByUidAndClusterId(uid, clusterId)
	if err != nil {
		return false, err
	}
	if len(events) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (c clusterEventService) Save(event model.ClusterEvent) error {
	return c.clusterEventRepo.Save(&event)
}
