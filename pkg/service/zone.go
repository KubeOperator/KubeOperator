package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider/client"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
)

type ZoneService interface {
	Get(name string) (dto.Zone, error)
	List() ([]dto.Zone, error)
	Page(num, size int) (page.Page, error)
	Delete(name string) error
	Create(creation dto.ZoneCreate) (dto.Zone, error)
	Batch(op dto.ZoneOp) error
	ListClusters(creation dto.CloudZoneRequest) ([]string, error)
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
		zoneDTOs = append(zoneDTOs, dto.Zone{Zone: mo})
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

	zone := model.Zone{
		BaseModel: common.BaseModel{},
		Name:      creation.Name,
		//Vars:       creation.Vars,
		//Datacenter: creation.Datacenter,
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

func (z zoneService) ListClusters(creation dto.CloudZoneRequest) ([]string, error) {
	cloudClient := client.NewCloudClient(creation.CloudVars.(map[string]interface{}))
	var result []string
	if cloudClient != nil {
		result, err := cloudClient.ListClusters(creation.Datacenter)
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
