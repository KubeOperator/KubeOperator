package service

import (
	"encoding/json"
	"sort"
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
		var storage []dto.NameVersion
		if err := json.Unmarshal([]byte(mo.StorageVars), &storage); err != nil {
			return clusterManifests, err
		}
		clusterManifest.StorageVars = storage
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
	sortManifest(largeVersions)
	for _, largeVersion := range largeVersions {
		var clusterManifestGroup dto.ClusterManifestGroup
		clusterManifestGroup.LargeVersion = largeVersion.Version
		var manifests []model.ClusterManifest
		db.DB.Where("version LIKE ?", "%"+largeVersion.Version+"%").Find(&manifests)
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

func sortManifest(arr []dto.ClusterManifest) []dto.ClusterManifest {
	if len(arr) < 1 {
		return arr
	}
	sort.SliceStable(arr, func(i, j int) bool {
		return compareVersion(arr[i], arr[j])
	})
	return arr
}

func getKoVersion(manifest dto.ClusterManifest) float64 {
	koIndex := strings.Index(manifest.Name, "ko")
	koVersionString := manifest.Name[koIndex+2:]
	version, err := strconv.ParseFloat(koVersionString, 64)
	if err != nil {
		log.Errorf("ko version %s parse float failed: %v", koVersionString, err)
	}
	return version
}

func compareVersion(version1 dto.ClusterManifest, version2 dto.ClusterManifest) bool {
	v1slice := getVersionSlice(version1)
	v2slice := getVersionSlice(version2)
	if len(v1slice) < 3 || len(v2slice) < 3 {
		return false
	}

	if getVersionNumber(v1slice[0]) > getVersionNumber(v2slice[0]) {
		return true
	} else if getVersionNumber(v1slice[0]) == getVersionNumber(v2slice[0]) {
		if getVersionNumber(v1slice[1]) > getVersionNumber(v2slice[1]) {
			return true
		} else if getVersionNumber(v1slice[1]) == getVersionNumber(v2slice[1]) {
			if getVersionNumber(v1slice[2]) > getVersionNumber(v2slice[2]) {
				return true
			} else if getVersionNumber(v1slice[2]) == getVersionNumber(v2slice[2]) {
				if getKoVersion(version1) > getKoVersion(version2) {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}

func getVersionNumber(versionStr string) float64 {
	version, err := strconv.ParseFloat(versionStr, 64)
	if err != nil {
		log.Errorf("ko version %s parse float failed: %v", versionStr, err)
	}
	return version
}

func getVersionSlice(version dto.ClusterManifest) []string {
	versionNumStr := strings.Replace(version.Version, "v", "", -1)
	slice := strings.Split(versionNumStr, ".")
	return slice
}
