package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider/client"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
)

type RegionService interface {
	Get(name string) (dto.Region, error)
	List() ([]dto.Region, error)
	Page(num, size int) (page.Page, error)
	Delete(name string) error
	Create(creation dto.RegionCreate) (dto.Region, error)
	Batch(op dto.RegionOp) error
	CheckValid(creation dto.RegionCreate) error
}

type regionService struct {
	regionRepo repository.RegionRepository
}

func NewRegionService() RegionService {
	return &regionService{
		regionRepo: repository.NewRegionRepository(),
	}
}

func (r regionService) Get(name string) (dto.Region, error) {
	var regionDTO dto.Region
	mo, err := r.regionRepo.Get(name)
	if err != nil {
		return regionDTO, err
	}
	regionDTO.Region = mo
	return regionDTO, err
}

func (r regionService) List() ([]dto.Region, error) {
	var regionDTOs []dto.Region
	mos, err := r.regionRepo.List()
	if err != nil {
		return regionDTOs, err
	}
	for _, mo := range mos {
		regionDTOs = append(regionDTOs, dto.Region{Region: mo})
	}
	return regionDTOs, err
}
func (r regionService) Page(num, size int) (page.Page, error) {
	var page page.Page
	var regionDTOs []dto.Region
	total, mos, err := r.regionRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		regionDTOs = append(regionDTOs, dto.Region{Region: mo})
	}
	page.Total = total
	page.Items = regionDTOs
	return page, err
}

func (r regionService) Delete(name string) error {
	err := r.regionRepo.Delete(name)
	if err != nil {
		return err
	}
	return nil
}

func (r regionService) Create(creation dto.RegionCreate) (dto.Region, error) {

	region := model.Region{
		BaseModel: common.BaseModel{},
		Name:      creation.Name,
		//Vars:       creation.Vars,
		Datacenter: creation.Datacenter,
	}

	err := r.regionRepo.Save(&region)
	if err != nil {
		return dto.Region{}, err
	}
	return dto.Region{Region: region}, err
}

func (r regionService) Batch(op dto.RegionOp) error {
	var deleteItems []model.Region
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.Region{
			BaseModel: common.BaseModel{},
			ID:        item.ID,
			Name:      item.Name,
		})
	}
	err := r.regionRepo.Batch(op.Operation, deleteItems)
	if err != nil {
		return err
	}
	return nil
}

func (r regionService) CheckValid(creation dto.RegionCreate) error {
	cloudClient := client.NewCloudClient(creation.RegionVars.(map[string]interface{}))
	if cloudClient != nil {
		_, err := cloudClient.ListDatacenter()
		if err != nil {
			return err
		}
	}
	return nil
}
