package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type ClusterStorageProvisionerService interface {
	ListStorageProvisioner(clusterName string) ([]dto.ClusterStorageProvisioner, error)
	CreateStorageProvisioner(clusterName string, creation dto.ClusterStorageProvisionerCreation) (dto.ClusterStorageProvisioner, error)
}

type clusterStorageProvisionerService struct {
	provisionerRepo repository.ClusterStorageProvisionerRepository
}

func NewClusterStorageProvisionerService() ClusterStorageProvisionerService {
	return &clusterStorageProvisionerService{
		provisionerRepo: repository.NewClusterStorageProvisionerRepository(),
	}
}

func (c clusterStorageProvisionerService) ListStorageProvisioner(clusterName string) ([]dto.ClusterStorageProvisioner, error) {
	var clusterStorageProvisionerDTOS []dto.ClusterStorageProvisioner
	ps, err := c.provisionerRepo.List(clusterName)
	if err != nil {
		return clusterStorageProvisionerDTOS, err
	}
	for _, p := range ps {
		var vars map[string]interface{}
		_ = json.Unmarshal([]byte(p.Vars), &vars)
		clusterStorageProvisionerDTOS = append(clusterStorageProvisionerDTOS, dto.ClusterStorageProvisioner{
			ClusterStorageProvisioner: p,
			Vars:                      vars,
		})
	}
	return clusterStorageProvisionerDTOS, nil
}

func (c clusterStorageProvisionerService) CreateStorageProvisioner(clusterName string, creation dto.ClusterStorageProvisionerCreation) (dto.ClusterStorageProvisioner, error) {
	vars, _ := json.Marshal(creation.Vars)
	var dp dto.ClusterStorageProvisioner
	p := model.ClusterStorageProvisioner{
		Name:   creation.Name,
		Type:   creation.Type,
		Vars:   string(vars),
		Status: constant.ClusterWaiting,
	}
	err := c.provisionerRepo.Save(clusterName, &p)
	if err != nil {
		return dp, err
	}
	dp.ClusterStorageProvisioner = p
	_ = json.Unmarshal([]byte(p.Vars), &dp.Vars)
	return dp, nil
}
