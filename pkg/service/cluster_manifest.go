package service

import (
	"encoding/json"
	"fmt"
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
		if err := json.Unmarshal([]byte(mo.CoreVars), &core); err != nil {
			fmt.Printf("func (c clusterManifestService) List(mo.CoreVars) json.Unmarshal err: %v\n", err)
		}
		clusterManifest.CoreVars = core
		var network []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.NetworkVars), &network); err != nil {
			fmt.Printf("func (c clusterManifestService) List(mo.NetworkVars) json.Unmarshal err: %v\n", err)
		}
		clusterManifest.NetworkVars = network
		var other []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.OtherVars), &other); err != nil {
			fmt.Printf("func (c clusterManifestService) List(mo.OtherVars) json.Unmarshal err: %v\n", err)
		}
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
		if err := json.Unmarshal([]byte(mo.CoreVars), &core); err != nil {
			fmt.Printf("func (c clusterManifestService) ListActive(mo.CoreVars) json.Unmarshal err: %v\n", err)
		}
		clusterManifest.CoreVars = core
		var network []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.NetworkVars), &network); err != nil {
			fmt.Printf("func (c clusterManifestService) ListActive(mo.NetworkVars) json.Unmarshal err: %v\n", err)
		}
		clusterManifest.NetworkVars = network
		var other []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.OtherVars), &other); err != nil {
			fmt.Printf("func (c clusterManifestService) ListActive(mo.OtherVars) json.Unmarshal err: %v\n", err)
		}
		clusterManifest.OtherVars = other
		clusterManifests = append(clusterManifests, clusterManifest)
	}
	return clusterManifests, err
}

func (c clusterManifestService) Update(update dto.ClusterManifestUpdate) (model.ClusterManifest, error) {
	var manifest model.ClusterManifest
	manifest, err := c.clusterManifestRepo.Get(update.Version)
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
