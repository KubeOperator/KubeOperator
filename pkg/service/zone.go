package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"

	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_storage"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type ZoneService interface {
	Get(name string) (*dto.Zone, error)
	List(conditions condition.Conditions) ([]dto.Zone, error)
	Page(num, size int, conditions condition.Conditions) (*page.Page, error)
	Delete(name string) error
	Create(creation dto.ZoneCreate) (*dto.Zone, error)
	Update(name string, creation dto.ZoneUpdate) (*dto.Zone, error)
	Batch(op dto.ZoneOp) error
	ListClusters(creation dto.CloudZoneRequest) ([]interface{}, error)
	ListTemplates(creation dto.CloudZoneRequest) ([]interface{}, error)
	ListByRegionName(regionName string) ([]dto.Zone, error)
	ListDatastores(creation dto.CloudZoneRequest) ([]dto.CloudDatastore, error)
}

type zoneService struct {
	zoneRepo             repository.ZoneRepository
	regionRepo           repository.RegionRepository
	systemSettingService SystemSettingService
	ipPoolService        IpPoolService
}

func NewZoneService() ZoneService {
	return &zoneService{
		zoneRepo:             repository.NewZoneRepository(),
		systemSettingService: NewSystemSettingService(),
		regionRepo:           repository.NewRegionRepository(),
		ipPoolService:        NewIpPoolService(),
	}
}

func (z zoneService) Get(name string) (*dto.Zone, error) {
	var (
		zoneDTO dto.Zone
		zone    model.Zone
	)
	if err := db.DB.Model(model.Zone{}).Where("name = ?", name).Preload("Region").Preload("IpPool").Preload("IpPool.Ips").Find(&zone).Error; err != nil {
		return nil, err
	}
	zoneDTO.Zone = zone
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(zone.Vars), &m); err != nil {
		return nil, err
	}
	zoneDTO.CloudVars = m
	zoneDTO.RegionName = zone.Region.Name
	zoneDTO.Provider = zone.Region.Provider
	ipUsed := 0
	for _, ip := range zone.IpPool.Ips {
		if ip.Status != constant.IpAvailable {
			ipUsed++
		}
	}
	zoneDTO.IpPool = dto.IpPool{
		IpUsed: ipUsed,
		IpPool: zone.IpPool,
	}
	zoneDTO.IpPoolName = zone.IpPool.Name
	return &zoneDTO, nil
}

func (z zoneService) List(conditions condition.Conditions) ([]dto.Zone, error) {

	var (
		zoneDTOs []dto.Zone
		zones    []model.Zone
	)
	d := db.DB.Model(model.Zone{})
	if err := dbUtil.WithConditions(&d, model.Zone{}, conditions); err != nil {
		return nil, err
	}
	err := d.Preload("IpPool").Find(&zones).Error
	if err != nil {
		return zoneDTOs, err
	}
	for _, mo := range zones {
		zoneDTOs = append(zoneDTOs, dto.Zone{Zone: mo, RegionName: mo.Region.Name})
	}
	return zoneDTOs, err
}

func (z zoneService) Page(num, size int, conditions condition.Conditions) (*page.Page, error) {

	var (
		p        page.Page
		zoneDTOs []dto.Zone
		zones    []model.Zone
	)

	d := db.DB.Model(model.Zone{})
	if err := dbUtil.WithConditions(&d, model.Zone{}, conditions); err != nil {
		return nil, err
	}
	err := d.
		Count(&p.Total).
		Offset((num - 1) * size).
		Limit(size).
		Preload("Region").
		Preload("IpPool").
		Preload("IpPool.Ips").
		Find(&zones).
		Error

	for _, mo := range zones {
		zoneDTO := new(dto.Zone)
		m := make(map[string]interface{})
		zoneDTO.Zone = mo
		if err := json.Unmarshal([]byte(mo.Vars), &m); err != nil {
			return nil, err
		}
		zoneDTO.CloudVars = m
		zoneDTO.RegionName = mo.Region.Name
		zoneDTO.Provider = mo.Region.Provider
		ipUsed := 0
		for _, ip := range mo.IpPool.Ips {
			if ip.Status != constant.IpAvailable {
				ipUsed++
			}
		}
		zoneDTO.IpPool = dto.IpPool{
			IpUsed: ipUsed,
			IpPool: mo.IpPool,
		}
		zoneDTO.IpPoolName = mo.IpPool.Name
		zoneDTOs = append(zoneDTOs, *zoneDTO)
	}
	p.Items = zoneDTOs
	return &p, err
}

func (z zoneService) Delete(name string) error {
	zone, err := z.zoneRepo.Get(name)
	if err != nil {
		return err
	}
	if err := db.DB.Delete(&zone).Error; err != nil {
		return err
	}
	return nil
}

func (z zoneService) Create(creation dto.ZoneCreate) (*dto.Zone, error) {
	var (
		credential dto.Credential
		repo       model.SystemRegistry
		region     dto.Region
		old        model.Zone
	)
	if err := db.DB.Where("architecture = ?", constant.ArchitectureOfAMD64).First(&repo).Error; err != nil {
		return nil, errors.New("IP_NOT_EXISTS")
	}

	db.DB.Model(model.Zone{}).Where("name = ?", creation.Name).Find(&old)
	if old.ID != "" {
		return nil, errors.New("NAME_EXISTS")
	}

	param := creation.CloudVars.(map[string]interface{})
	region, err := NewRegionService().Get(creation.RegionName)
	if err != nil {
		return nil, err
	}
	if param["templateType"] != nil && param["templateType"].(string) == "default" {
		switch region.Provider {
		case constant.OpenStack:
			param["imageName"] = constant.OpenStackImageName
		case constant.VSphere:
			param["imageName"] = constant.VSphereImageName
		case constant.FusionCompute:
			param["template"] = constant.FusionComputeImageName
		default:
			param["imageName"] = constant.VSphereImageName
		}
		credentialService := NewCredentialService()
		credential, err = credentialService.Get(constant.ImageCredentialName)
		if err != nil {
			return nil, err
		}
	} else {
		credentialService := NewCredentialService()
		credential, err = credentialService.Get(creation.CredentialName)
		if err != nil {
			return nil, err
		}
	}

	if region.Provider == constant.VSphere {
		regionVars := region.RegionVars.(map[string]interface{})
		regionVars["datacenter"] = region.Datacenter
		cloudClient := cloud_provider.NewCloudClient(regionVars)
		err = cloudClient.CreateDefaultFolder()
		if err != nil {
			return nil, err
		}
	}

	ipPool, err := z.ipPoolService.Get(creation.IpPoolName)
	if err != nil {
		return nil, err
	}
	if len(ipPool.Ips) == 0 {
		return nil, errors.New("IP_SHORT")
	}
	index := strings.Index(ipPool.Subnet, "/")
	networkCidr := ipPool.Subnet
	param["netMask"] = networkCidr[index+1:]
	param["gateway"] = ipPool.Ips[0].Gateway
	param["dns1"] = ipPool.Ips[0].DNS1
	param["dns2"] = ipPool.Ips[0].DNS2

	vars, _ := json.Marshal(creation.CloudVars)
	zone := model.Zone{
		BaseModel:    common.BaseModel{},
		Name:         creation.Name,
		Vars:         string(vars),
		RegionID:     region.ID,
		CredentialID: credential.ID,
		IpPoolID:     ipPool.ID,
		Status:       constant.Ready,
	}

	if param["templateType"] != nil && param["templateType"].(string) == "default" {
		zone.Status = constant.Initializing
	}
	err = z.zoneRepo.Save(&zone)
	if err != nil {
		return nil, err
	}
	if param["templateType"] != nil && param["templateType"].(string) == "default" {
		go z.uploadZoneImage(creation)
	}
	return &dto.Zone{Zone: zone}, err
}

func (z zoneService) Update(name string, update dto.ZoneUpdate) (*dto.Zone, error) {

	param := update.CloudVars.(map[string]interface{})
	ipPool, err := z.ipPoolService.Get(update.IpPoolName)
	if err != nil {
		return nil, err
	}
	if len(ipPool.Ips) == 0 {
		return nil, errors.New("IP_SHORT")
	}

	index := strings.Index(ipPool.Subnet, "/")
	networkCidr := ipPool.Subnet
	param["netMask"] = networkCidr[index+1:]
	param["gateway"] = ipPool.Ips[0].Gateway
	param["dns1"] = ipPool.Ips[0].DNS1
	param["dns2"] = ipPool.Ips[0].DNS2

	vars, _ := json.Marshal(update.CloudVars)
	zone, err := z.zoneRepo.Get(name)
	if err != nil {
		return nil, err
	}
	zone.Vars = string(vars)
	zone.RegionID = update.RegionID
	zone.IpPoolID = ipPool.ID
	if err := db.DB.Save(&zone).Error; err != nil {
		return nil, err
	}
	return &dto.Zone{Zone: zone}, err
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
	cloudClient := cloud_provider.NewCloudClient(creation.CloudVars.(map[string]interface{}))
	var result []interface{}
	if cloudClient != nil {
		result, err := cloudClient.ListClusters()
		if err != nil {
			return result, err
		}
		if result == nil {
			return result, errors.New("CLUSTER_IS_NULL")
		}
		return result, err
	}
	return result, nil
}

func (z zoneService) ListTemplates(creation dto.CloudZoneRequest) ([]interface{}, error) {
	var result []interface{}
	var clientVars map[string]interface{}
	if creation.RegionName != "" {
		region, err := z.regionRepo.Get(creation.RegionName)
		if err != nil {
			return result, err
		}
		creation.Datacenter = region.Datacenter
		m := make(map[string]interface{})
		if err := json.Unmarshal([]byte(region.Vars), &m); err != nil {
			return result, err
		}
		vars := creation.CloudVars.(map[string]interface{})
		m["cluster"] = vars["cluster"].(string)
		m["datacenter"] = region.Datacenter
		clientVars = m
	} else {
		clientVars = creation.CloudVars.(map[string]interface{})
	}
	cloudClient := cloud_provider.NewCloudClient(clientVars)

	if cloudClient != nil {
		result, err := cloudClient.ListTemplates()
		if err != nil {
			return result, err
		}
		if result == nil {
			return result, errors.New("IMAGE_IS_NULL")
		}
		return result, err
	}
	return result, nil
}

func (z zoneService) uploadZoneImage(creation dto.ZoneCreate) {
	zone, err := z.zoneRepo.Get(creation.Name)
	if err != nil {
		logger.Log.Error(err)
	}
	err = z.uploadImage(creation)
	if err != nil {
		logger.Log.Error(err)
		zone.Status = constant.UploadImageError
	} else {
		zone.Status = constant.Ready
	}
	err = z.zoneRepo.Save(&zone)
	if err != nil {
		logger.Log.Error(err)
	}
}

func (z zoneService) uploadImage(creation dto.ZoneCreate) error {
	region, err := NewRegionService().Get(creation.RegionName)
	if err != nil {
		return err
	}
	var repo model.SystemRegistry
	if err := db.DB.Where("architecture = ?", constant.ArchitectureOfAMD64).First(&repo).Error; err != nil {
		return fmt.Errorf("can't find local ip from system setting, err %s", err.Error())
	}
	ip := repo.Hostname

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
		regionVars["ovfPath"] = fmt.Sprintf(constant.VSphereImageOvfPath, ip)
		regionVars["vmdkPath"] = fmt.Sprintf(constant.VSphereImageVMDkPath, ip)
	}
	if region.Provider == constant.OpenStack {
		regionVars["imagePath"] = fmt.Sprintf(constant.OpenStackImagePath, ip)
	}
	if region.Provider == constant.FusionCompute {
		zoneVars := creation.CloudVars.(map[string]interface{})
		if zoneVars["cluster"] != nil {
			regionVars["cluster"] = zoneVars["cluster"]
		}
		if zoneVars["datastore"] != nil {
			regionVars["datastore"] = zoneVars["datastore"]
		}
		if zoneVars["portgroup"] != nil {
			regionVars["portgroup"] = zoneVars["portgroup"]
		}
	}

	cloudClient := cloud_provider.NewCloudClient(regionVars)
	if cloudClient != nil {
		result, err := cloudClient.DefaultImageExist()
		if err != nil {
			return err
		}
		if result {
			return nil
		}
		if region.Provider == constant.FusionCompute {
			zoneVars := creation.CloudVars.(map[string]interface{})
			nfsVars := make(map[string]interface{})
			nfsVars["type"] = "SFTP"
			nfsVars["address"] = zoneVars["nfsAddress"]
			nfsVars["port"] = zoneVars["nfsPort"]
			nfsVars["username"] = zoneVars["nfsUsername"]
			nfsVars["password"] = zoneVars["nfsPassword"]
			nfsVars["bucket"] = zoneVars["nfsFolder"]
			client, err := cloud_storage.NewCloudStorageClient(nfsVars)
			if err != nil {
				return err
			}
			ovfResp, err := http.Get(fmt.Sprintf(constant.FusionComputeOvfPath, ip))
			if err != nil {
				return err
			}
			if ovfResp.StatusCode == 404 {
				return errors.New(constant.FusionComputeOvfName + "not found")
			}
			defer ovfResp.Body.Close()
			ovfOut, err := os.Create(constant.FusionComputeOvfLocal)
			if err != nil {
				return err
			}
			defer ovfOut.Close()
			_, err = io.Copy(ovfOut, ovfResp.Body)
			if err != nil {
				return err
			}
			vhdResp, err := http.Get(fmt.Sprintf(constant.FusionComputeVhdPath, ip))
			if err != nil {
				return err
			}
			if vhdResp.StatusCode == 404 {
				return errors.New(constant.FusionComputeVhdName + "not found")
			}
			defer vhdResp.Body.Close()
			vhdOut, err := os.Create(constant.FusionComputeVhdLocal)
			if err != nil {
				return err
			}
			defer vhdOut.Close()
			_, err = io.Copy(vhdOut, vhdResp.Body)
			if err != nil {
				return err
			}
			_, err = client.Upload(constant.FusionComputeOvfLocal, constant.FusionComputeOvfName)
			if err != nil {
				return err
			}
			_, err = client.Upload(constant.FusionComputeVhdLocal, constant.FusionComputeVhdName)
			if err != nil {
				return err
			}
			regionVars["ovfPath"] = zoneVars["nfsAddress"].(string) + ":" + zoneVars["nfsFolder"].(string) + "/" + constant.FusionComputeOvfName
			cloudClient = cloud_provider.NewCloudClient(regionVars)
		}
		err = cloudClient.UploadImage()
		if err != nil {
			return err
		}
	}
	return nil
}

func (z zoneService) ListByRegionName(regionName string) ([]dto.Zone, error) {
	var zoneDTOs []dto.Zone
	region, err := z.regionRepo.Get(regionName)
	if err != nil {
		return nil, err
	}
	mos, err := z.zoneRepo.ListByRegionId(region.ID)
	if err != nil {
		return zoneDTOs, err
	}
	for _, mo := range mos {
		zoneDTOs = append(zoneDTOs, dto.Zone{Zone: mo})
	}
	return zoneDTOs, err
}

func (z zoneService) ListDatastores(creation dto.CloudZoneRequest) ([]dto.CloudDatastore, error) {
	var result []dto.CloudDatastore
	var clientVars map[string]interface{}
	if creation.RegionName != "" {
		region, err := z.regionRepo.Get(creation.RegionName)
		if err != nil {
			return result, err
		}
		creation.Datacenter = region.Datacenter
		m := make(map[string]interface{})
		if err := json.Unmarshal([]byte(region.Vars), &m); err != nil {
			return result, err
		}
		vars := creation.CloudVars.(map[string]interface{})
		m["cluster"] = vars["cluster"].(string)
		m["datacenter"] = region.Datacenter
		clientVars = m
	} else {
		clientVars = creation.CloudVars.(map[string]interface{})
	}
	cloudClient := cloud_provider.NewCloudClient(clientVars)
	datastores, err := cloudClient.ListDatastores()

	for i := range datastores {
		result = append(result, dto.CloudDatastore{
			Name:      datastores[i].Name,
			Capacity:  datastores[i].Capacity,
			FreeSpace: datastores[i].FreeSpace,
		})
	}

	if err != nil {
		return result, err
	}
	return result, err
}
