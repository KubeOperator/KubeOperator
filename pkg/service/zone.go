package service

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider/client"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	ZoneNameExist = "NAME_EXISTS"
)

type ZoneService interface {
	Get(name string) (dto.Zone, error)
	List() ([]dto.Zone, error)
	Page(num, size int) (page.Page, error)
	Delete(name string) error
	Create(creation dto.ZoneCreate) (*dto.Zone, error)
	Update(creation dto.ZoneUpdate) (dto.Zone, error)
	Batch(op dto.ZoneOp) error
	ListClusters(creation dto.CloudZoneRequest) ([]interface{}, error)
	ListTemplates(creation dto.CloudZoneRequest) ([]interface{}, error)
	ListByRegionId(regionId string) ([]dto.Zone, error)
	DownloadVMDKFile(vmdkUrl string) (string, error)
}

type zoneService struct {
	zoneRepo             repository.ZoneRepository
	systemSettingService SystemSettingService
}

func NewZoneService() ZoneService {
	return &zoneService{
		zoneRepo:             repository.NewZoneRepository(),
		systemSettingService: NewSystemSettingService(),
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
		zoneDTOs = append(zoneDTOs, dto.Zone{Zone: mo, RegionName: mo.Region.Name})
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

		zoneDTO.RegionName = mo.Region.Name
		zoneDTO.Provider = mo.Region.Provider

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

func (z zoneService) Create(creation dto.ZoneCreate) (*dto.Zone, error) {

	old, _ := z.Get(creation.Name)
	if old.ID != "" {
		return nil, errors.New(ZoneNameExist)
	}

	param := creation.CloudVars.(map[string]interface{})
	if param["subnet"] != nil {
		index := strings.Index(param["subnet"].(string), "/")
		networkCidr := param["subnet"].(string)
		param["netMask"] = networkCidr[index+1:]
	}

	if param["templateType"] != nil && param["templateType"].(string) == "default" {
		region, err := NewRegionService().Get(creation.RegionName)
		if err != nil {
			return nil, err
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
			return nil, err
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
		return nil, err
	}
	go z.uploadImage(creation)
	return &dto.Zone{Zone: zone}, err
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
	ip, err := z.systemSettingService.Get("ip")
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
		regionVars["ovfPath"] = fmt.Sprintf(constant.VSphereImageOvfPath, ip.Value)
	}
	if region.Provider == constant.OpenStack {
		regionVars["imagePath"] = fmt.Sprintf(constant.OpenStackImagePath, ip.Value)
	}
	cloudClient := client.NewCloudClient(regionVars)
	if cloudClient != nil {
		result, err := cloudClient.DefaultImageExist()
		zone, err := z.zoneRepo.Get(creation.Name)
		if err != nil {
			return err
		}
		if result {
			zone.Status = constant.Ready
			err = z.zoneRepo.Save(&zone)
			if err != nil {
				return err
			}
			return nil
		} else {
			if region.Provider == constant.VSphere {
				vmdkUrl := fmt.Sprintf(constant.VSphereImageVMDkPath, ip.Value)
				regionVars["vmdkPath"], err = z.DownloadVMDKFile(vmdkUrl)
				if err != nil {
					return err
				}
			}
		}
		err = cloudClient.UploadImage()
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

func (z zoneService) DownloadVMDKFile(vmdkUrl string) (string, error) {
	res, err := http.Get(vmdkUrl)
	if err != nil {
		return "", err
	}
	f, err := os.Create(constant.VMDKGZLocalPath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return "", err
	}
	fw, err := os.Open(constant.VMDKGZLocalPath)
	if err != nil {
		return "", err
	}
	defer fw.Close()
	gr, err := gzip.NewReader(fw)
	if err != nil {
		return "", err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		_, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		fw, err := os.OpenFile(constant.VMDKLocalPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}
		defer fw.Close()
		_, err = io.Copy(fw, tr)
		if err != nil {
			return "", err
		}
	}

	_, err = os.Open(constant.VMDKLocalPath)
	if err != nil {
		return "", err
	}
	return constant.VMDKLocalPath, err
}
