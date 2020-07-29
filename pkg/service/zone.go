package service

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider/client"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/jinzhu/gorm"
	"strings"
)

type ZoneService interface {
	Get(name string) (dto.Zone, error)
	List() ([]dto.Zone, error)
	Page(num, size int) (page.Page, error)
	Delete(name string) error
	Create(creation dto.ZoneCreate) (dto.Zone, error)
	Update(creation dto.ZoneUpdate) (dto.Zone, error)
	Batch(op dto.ZoneOp) error
	ListClusters(creation dto.CloudZoneRequest) ([]interface{}, error)
	ListTemplates(creation dto.CloudZoneRequest) ([]interface{}, error)
	ListByRegionId(regionId string) ([]dto.Zone, error)
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService() ZoneService {
	return &zoneService{
		zoneRepo: repository.NewZoneRepository(),
	}
}

func (z zoneService) Get(name string) (dto.Zone, error) {
	var zoneDTO dto.Zone
	mo, err := z.zoneRepo.Get(name)
	if err != nil {
		return zoneDTO, err
	}
	zoneDTO.Zone = mo
	return zoneDTO, err
}

func (z zoneService) List() ([]dto.Zone, error) {
	var zoneDTOs []dto.Zone
	mos, err := z.zoneRepo.List()
	if err != nil {
		return zoneDTOs, err
	}
	for _, mo := range mos {
		zoneDTOs = append(zoneDTOs, dto.Zone{Zone: mo})
	}
	return zoneDTOs, err
}

func (z zoneService) Page(num, size int) (page.Page, error) {
	var page page.Page
	var zoneDTOs []dto.Zone
	total, mos, err := z.zoneRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		zoneDTO := new(dto.Zone)
		m := make(map[string]interface{})
		zoneDTO.Zone = mo
		json.Unmarshal([]byte(mo.Vars), &m)
		zoneDTO.CloudVars = m

		regionDTO := new(dto.Region)
		r := make(map[string]interface{})
		json.Unmarshal([]byte(mo.Region.Vars), &r)
		regionDTO.RegionVars = r
		regionDTO.Region = mo.Region
		zoneDTO.Region = *regionDTO

		zoneDTOs = append(zoneDTOs, *zoneDTO)
	}
	page.Total = total
	page.Items = zoneDTOs
	return page, err
}

func (z zoneService) Delete(name string) error {

	err := z.zoneRepo.Delete(name)
	if err != nil {
		return err
	}
	return nil
}

func (z zoneService) Create(creation dto.ZoneCreate) (dto.Zone, error) {

	param := creation.CloudVars.(map[string]interface{})
	if param["subnet"] != nil {
		index := strings.Index(param["subnet"].(string), "/")
		networkCidr := param["subnet"].(string)
		param["netMask"] = networkCidr[index+1:]
	}

	if param["templateType"] != nil && param["templateType"].(string) == "default" {
		region, err := NewRegionService().Get(creation.RegionName)
		if err != nil {
			return dto.Zone{}, err
		}
		switch region.Provider {
		case constant.OpenStack:
			param["imageName"] = constant.OpenStackImageName
			break
		case constant.VSphere:
			param["imageName"] = constant.VSphereImageName
			break
		default:
			param["imageName"] = constant.VSphereImageName
			break
		}
		credentialService := NewCredentialService()
		credential, err := credentialService.Get(constant.ImageCredentialName)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				credential, err = credentialService.Create(dto.CredentialCreate{
					Name:     constant.ImageCredentialName,
					Password: constant.ImageDefaultPassword,
					Username: constant.ImageUserName,
					Type:     constant.ImagePasswordType,
				})
				if err != nil {
					return dto.Zone{}, err
				}
			} else {
				return dto.Zone{}, err
			}
		}
		creation.CredentialId = credential.ID
	}

	vars, _ := json.Marshal(creation.CloudVars)

	zone := model.Zone{
		BaseModel:    common.BaseModel{},
		Name:         creation.Name,
		Vars:         string(vars),
		RegionID:     creation.RegionID,
		CredentialID: creation.CredentialId,
		Status:       constant.Initializing,
	}

	err := z.zoneRepo.Save(&zone)
	if err != nil {
		return dto.Zone{}, err
	}
	go z.uploadImage(creation)
	return dto.Zone{Zone: zone}, err
}

func (z zoneService) Update(creation dto.ZoneUpdate) (dto.Zone, error) {

	vars, _ := json.Marshal(creation.CloudVars)

	zone := model.Zone{
		BaseModel: common.BaseModel{},
		Name:      creation.Name,
		Vars:      string(vars),
		RegionID:  creation.RegionID,
		ID:        creation.ID,
	}

	err := z.zoneRepo.Save(&zone)
	if err != nil {
		return dto.Zone{}, err
	}
	return dto.Zone{Zone: zone}, err
}

func (z zoneService) Batch(op dto.ZoneOp) error {
	var deleteItems []model.Zone
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.Zone{
			BaseModel: common.BaseModel{},
			ID:        item.ID,
			Name:      item.Name,
		})
	}
	err := z.zoneRepo.Batch(op.Operation, deleteItems)
	if err != nil {
		return err
	}
	return nil
}

func (z zoneService) ListClusters(creation dto.CloudZoneRequest) ([]interface{}, error) {
	cloudClient := client.NewCloudClient(creation.CloudVars.(map[string]interface{}))
	var result []interface{}
	if cloudClient != nil {
		result, err := cloudClient.ListClusters()
		if err != nil {
			return result, err
		}
		if result == nil {
			return result, errors.New("cluster is null")
		}
		return result, err
	}
	return result, nil
}

func (z zoneService) ListTemplates(creation dto.CloudZoneRequest) ([]interface{}, error) {
	cloudClient := client.NewCloudClient(creation.CloudVars.(map[string]interface{}))
	var result []interface{}
	if cloudClient != nil {
		result, err := cloudClient.ListTemplates()
		if err != nil {
			return result, err
		}
		if result == nil {
			return result, errors.New("cluster is null")
		}
		return result, err
	}
	return result, nil
}
func (z zoneService) uploadImage(creation dto.ZoneCreate) error {
	region, err := NewRegionService().Get(creation.RegionName)
	if err != nil {
		return err
	}
	regionVars := region.RegionVars.(map[string]interface{})
	regionVars["datacenter"] = region.Datacenter
	if region.Provider == constant.VSphere {
		zoneVars := creation.CloudVars.(map[string]interface{})
		if zoneVars["cluster"] != nil {
			regionVars["cluster"] = zoneVars["cluster"]
		}
		if zoneVars["resourcePool"] != nil {
			regionVars["resourcePool"] = zoneVars["resourcePool"]
		}
		if zoneVars["datastore"] != nil {
			regionVars["datastore"] = zoneVars["datastore"]
		}
		if zoneVars["network"] != nil {
			regionVars["network"] = zoneVars["network"]
		}
	}
	cloudClient := client.NewCloudClient(regionVars)
	if cloudClient != nil {
		err := cloudClient.UploadImage()
		if err != nil {
			zone, err := z.zoneRepo.Get(creation.Name)
			if err != nil {
				return err
			}
			zone.Status = constant.UploadImageError
			err = z.zoneRepo.Save(&zone)
			if err != nil {
				return err
			}
			return err
		}
		zone, err := z.zoneRepo.Get(creation.Name)
		if err != nil {
			return err
		}
		zone.Status = constant.Ready
		err = z.zoneRepo.Save(&zone)
		if err != nil {
			return err
		}
	}
	return nil
}

func (z zoneService) ListByRegionId(regionId string) ([]dto.Zone, error) {
	var zoneDTOs []dto.Zone
	mos, err := z.zoneRepo.ListByRegionId(regionId)
	if err != nil {
		return zoneDTOs, err
	}
	for _, mo := range mos {
		zoneDTOs = append(zoneDTOs, dto.Zone{Zone: mo})
	}
	return zoneDTOs, err
}
