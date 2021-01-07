package service

import (
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"strconv"
	"strings"
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
			return clusterManifests, err
		}
		clusterManifest.CoreVars = core
		var network []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.NetworkVars), &network); err != nil {
			return clusterManifests, err
		}
		clusterManifest.NetworkVars = network
		var other []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.OtherVars), &other); err != nil {
			return clusterManifests, err
		}
		clusterManifest.OtherVars = other
		clusterManifests = append(clusterManifests, clusterManifest)
	}

	return sortManifest(clusterManifests), err
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
	return sortManifest(clusterManifests), err
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

func sortManifest(mos []dto.ClusterManifest) []dto.ClusterManifest {

	version1s := make(map[string][]dto.ClusterManifest)
	for _, manifest := range mos {
		versionStr := strings.Replace(manifest.Version, "v", "", -1)
		version1Index := strings.Index(versionStr, ".")
		if version1Index == -1 {
			continue
		}
		version1 := versionStr[0:version1Index]
		if isExist(version1, version1s) {
			version1s[version1] = append(version1s[version1], manifest)
		} else {
			version1s[version1] = []dto.ClusterManifest{manifest}
		}
	}
	var result []dto.ClusterManifest
	for _, v := range version1s {
		quickSortVersion(v, 0, len(v)-1)
		result = append(result, v...)
	}
	return result
}

func isExist(version string, versions map[string][]dto.ClusterManifest) bool {
	for k, _ := range versions {
		if k == version {
			return true
		}
	}
	return false
}
func quickSortVersion(arr []dto.ClusterManifest, start, end int) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for getVersion(arr[i]) > getVersion(key) {
				i++
			}
			for getVersion(arr[j]) < getVersion(key) {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}

		if end > i {
			quickSortVersion(arr, i, end)
		}

		if start < j {
			quickSortVersion(arr, start, j)
		}
	}
}

func getVersion(manifest dto.ClusterManifest) float64 {
	versionStr := strings.Replace(manifest.Version, "v", "", -1)
	version1Index := strings.Index(versionStr, ".")
	version2 := strings.Replace(versionStr[version1Index+1:], ".", "", -1)
	version, _ := strconv.ParseFloat(version2, 64)
	return version
}
