package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type ClusterManifestService interface {
	List() ([]dto.ClusterManifest, error)
	ListActive() ([]dto.ClusterManifest, error)
	Update(update dto.ClusterManifestUpdate) (model.ClusterManifest, error)
}

type clusterManifestService struct {
	clusterManifestRepo repository.ClusterManifestRepository
}

func NewClusterManifestService() ClusterManifestService {
	return &clusterManifestService{
		clusterManifestRepo: repository.NewClusterManifestRepository(),
	}
}

func (c clusterManifestService) List() ([]dto.ClusterManifest, error) {
	var clusterManifests []dto.ClusterManifest
	mos, err := c.clusterManifestRepo.List()
	if err != nil {
		return clusterManifests, err
	}
	for _, mo := range mos {
		var clusterManifest dto.ClusterManifest
		clusterManifest.Name = mo.Name
		clusterManifest.Version = mo.Version
		clusterManifest.IsActive = mo.IsActive
		var core []dto.NameVersion
		json.Unmarshal([]byte(mo.CoreVars), &core)
		clusterManifest.CoreVars = core
		var network []dto.NameVersion
		json.Unmarshal([]byte(mo.NetworkVars), &network)
		clusterManifest.NetworkVars = network
		var other []dto.NameVersion
		json.Unmarshal([]byte(mo.OtherVars), &other)
		clusterManifest.OtherVars = other
		clusterManifests = append(clusterManifests, clusterManifest)
	}
	return clusterManifests, err
}

func (c clusterManifestService) ListActive() ([]dto.ClusterManifest, error) {
	var clusterManifests []dto.ClusterManifest
	mos, err := c.clusterManifestRepo.ListByStatus()
	if err != nil {
		return clusterManifests, err
	}
	for _, mo := range mos {
		var clusterManifest dto.ClusterManifest
		clusterManifest.Name = mo.Name
		clusterManifest.Version = mo.Version
		clusterManifest.IsActive = mo.IsActive
		var core []dto.NameVersion
		json.Unmarshal([]byte(mo.CoreVars), &core)
		clusterManifest.CoreVars = core
		var network []dto.NameVersion
		json.Unmarshal([]byte(mo.NetworkVars), &network)
		clusterManifest.NetworkVars = network
		var other []dto.NameVersion
		json.Unmarshal([]byte(mo.OtherVars), &other)
		clusterManifest.OtherVars = other
		clusterManifests = append(clusterManifests, clusterManifest)
	}
	return clusterManifests, err
}

func (c clusterManifestService) Update(update dto.ClusterManifestUpdate) (model.ClusterManifest, error) {
	var manifest model.ClusterManifest
	manifest, err := c.clusterManifestRepo.Get(update.Name)
	if err != nil {
		return manifest, err
	}
	manifest.IsActive = update.IsActive
	err = c.clusterManifestRepo.Save(manifest)
	if err != nil {
		return manifest, err
	}
	return manifest, err
}
