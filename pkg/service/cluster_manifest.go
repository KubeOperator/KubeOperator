package service

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/db"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type ClusterManifestService interface {
	List() ([]dto.ClusterManifest, error)
	ListActive() ([]dto.ClusterManifest, error)
	Update(update dto.ClusterManifestUpdate) (model.ClusterManifest, error)
	ListByLargeVersion() ([]dto.ClusterManifestGroup, error)
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

		var tool []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.ToolVars), &tool); err != nil {
			return clusterManifests, err
		}
		clusterManifest.ToolVars = tool

		var storage []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.StorageVars), &storage); err != nil {
			return clusterManifests, err
		}
		clusterManifest.StorageVars = storage

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
			log.Errorf("clusterManifestService ListActive(mo.CoreVars) json.Unmarshal failed, error: %s", err.Error())
		}
		clusterManifest.CoreVars = core
		var network []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.NetworkVars), &network); err != nil {
			log.Errorf("clusterManifestService ListActive(mo.NetworkVars) json.Unmarshal failed, error: %s", err.Error())
		}
		clusterManifest.NetworkVars = network
		var other []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.OtherVars), &other); err != nil {
			log.Errorf("clusterManifestService ListActive(mo.OtherVars) json.Unmarshal failed, error: %s", err.Error())
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

func (c clusterManifestService) ListByLargeVersion() ([]dto.ClusterManifestGroup, error) {

	var clusterManifestGroups []dto.ClusterManifestGroup
	var largeVersions []dto.ClusterManifest
	db.DB.Raw("select distinct substring_index(version,'.',2) version from ko_cluster_manifest").Scan(&largeVersions)
	if len(largeVersions) == 0 {
		return []dto.ClusterManifestGroup{}, nil
	}
	for _, largeVersion := range largeVersions {
		var clusterManifestGroup dto.ClusterManifestGroup
		clusterManifestGroup.LargeVersion = largeVersion.Version
		var manifests []model.ClusterManifest
		db.DB.Model(model.ClusterManifest{}).Where("version LIKE ?", "%"+largeVersion.Version+"%").Find(&manifests)
		if len(manifests) == 0 {
			continue
		}
		var clusterManifests []dto.ClusterManifest
		for _, mo := range manifests {
			var clusterManifest dto.ClusterManifest
			clusterManifest.Name = mo.Name
			clusterManifest.Version = mo.Version
			clusterManifest.IsActive = mo.IsActive
			var core []dto.NameVersion
			if err := json.Unmarshal([]byte(mo.CoreVars), &core); err != nil {
				log.Errorf("clusterManifestService ListActive(mo.CoreVars) json.Unmarshal failed, error: %s", err.Error())
			}

			clusterManifest.CoreVars = core
			var network []dto.NameVersion
			if err := json.Unmarshal([]byte(mo.NetworkVars), &network); err != nil {
				log.Errorf("clusterManifestService ListActive(mo.NetworkVars) json.Unmarshal failed, error: %s", err.Error())
			}
			clusterManifest.NetworkVars = network

			var tool []dto.NameVersion
			if err := json.Unmarshal([]byte(mo.ToolVars), &tool); err != nil {
				log.Errorf("clusterManifestService ListActive(mo.ToolVars) json.Unmarshal failed, error: %s", err.Error())
			}
			clusterManifest.ToolVars = tool

			var storage []dto.NameVersion
			if err := json.Unmarshal([]byte(mo.StorageVars), &storage); err != nil {
				log.Errorf("clusterManifestService ListActive(mo.StorageVars) json.Unmarshal failed, error: %s", err.Error())
			}
			clusterManifest.StorageVars = storage

			var other []dto.NameVersion
			if err := json.Unmarshal([]byte(mo.OtherVars), &other); err != nil {
				log.Errorf("clusterManifestService ListActive(mo.OtherVars) json.Unmarshal failed, error: %s", err.Error())
			}
			clusterManifest.OtherVars = other
			clusterManifests = append(clusterManifests, clusterManifest)
		}
		clusterManifestGroup.ClusterManifests = sortManifest(clusterManifests)
		clusterManifestGroups = append(clusterManifestGroups, clusterManifestGroup)
	}

	return clusterManifestGroups, nil
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
		sortKoVersion(v)
		result = append(result, v...)
	}
	return result
}

func isExist(version string, versions map[string][]dto.ClusterManifest) bool {
	for k := range versions {
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

func sortKoVersion(arr []dto.ClusterManifest) {
	var value dto.ClusterManifest
	for index := range arr {
		if index > 0 {
			value = arr[index-1]
		}
		if arr[index].Version == value.Version {
			if getKoVersion(value) < getKoVersion(arr[index]) {
				arr[index-1] = arr[index]
				arr[index] = value
			}
		}
	}
}

func getKoVersion(manifest dto.ClusterManifest) float64 {
	koIndex := strings.Index(manifest.Name, "ko")
	koVersionString := manifest.Name[koIndex+2:]
	version, _ := strconv.ParseFloat(koVersionString, 64)
	return version
}
