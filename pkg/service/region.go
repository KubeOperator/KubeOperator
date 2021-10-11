package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"

	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

var (
	DeleteRegionError = "DELETE_REGION_FAILED_RESOURCE"
	RegionNameExist   = "NAME_EXISTS"
)

type RegionService interface {
	Get(name string) (dto.Region, error)
	List(conditions condition.Conditions) ([]dto.Region, error)
	Page(num, size int, conditions condition.Conditions) (*page.Page, error)
	Delete(name string) error
	Create(creation dto.RegionCreate) (*dto.Region, error)
	Batch(op dto.RegionOp) error
	ListDatacenter(creation dto.RegionDatacenterRequest) ([]string, error)
	Update(name string, update dto.RegionUpdate) (*dto.Region, error)
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

	m := make(map[string]interface{})
	regionDTO.Region = mo
	if err := json.Unmarshal([]byte(mo.Vars), &m); err != nil {
		logger.Log.Errorf("regionService Get json.Unmarshal failed, error: %s", err.Error())
	}
	regionDTO.RegionVars = m

	return regionDTO, err
}

func (r regionService) List(conditions condition.Conditions) ([]dto.Region, error) {
	var (
		regionDTOs []dto.Region
		regions    []model.Region
	)
	d := db.DB.Model(model.Region{})
	if err := dbUtil.WithConditions(&d, model.Region{}, conditions); err != nil {
		return nil, err
	}
	err := d.Find(&regions).Error
	if err != nil {
		return nil, err
	}
	for _, mo := range regions {
		regionDTO := new(dto.Region)
		m := make(map[string]interface{})
		regionDTO.Region = mo
		if err := json.Unmarshal([]byte(mo.Vars), &m); err != nil {
			logger.Log.Errorf("regionService Page json.Unmarshal failed, error: %s", err.Error())
		}
		regionDTO.RegionVars = m
		regionDTOs = append(regionDTOs, *regionDTO)
	}
	return regionDTOs, err
}
func (r regionService) Page(num, size int, conditions condition.Conditions) (*page.Page, error) {

	var (
		p          page.Page
		regionDTOs []dto.Region
		regions    []model.Region
	)

	d := db.DB.Model(model.Region{})
	if err := dbUtil.WithConditions(&d, model.Region{}, conditions); err != nil {
		return nil, err
	}

	if err := d.Order("CONVERT(name using gbk) asc").Count(&p.Total).Offset((num - 1) * size).Limit(size).Find(&regions).Error; err != nil {
		return nil, err
	}
	for _, mo := range regions {
		regionDTO := new(dto.Region)
		m := make(map[string]interface{})
		regionDTO.Region = mo
		if err := json.Unmarshal([]byte(mo.Vars), &m); err != nil {
			logger.Log.Errorf("regionService Page json.Unmarshal failed, error: %s", err.Error())
		}
		regionDTO.RegionVars = m
		regionDTOs = append(regionDTOs, *regionDTO)
	}
	p.Items = regionDTOs
	return &p, nil
}

func (r regionService) Delete(name string) error {
	region, err := r.regionRepo.Get(name)
	if err != nil {
		return err
	}

	var zones []model.Zone
	if err := db.DB.Where("region_id = ?", region.ID).Find(&zones).Error; err != nil {
		return err
	}
	if len(zones) > 0 {
		return fmt.Errorf(DeleteRegionError)
	}
	if err := db.DB.Delete(&region).Error; err != nil {
		return err
	}
	return nil
}

func (r regionService) Create(creation dto.RegionCreate) (*dto.Region, error) {

	old, _ := r.Get(creation.Name)
	if old.ID != "" {
		return nil, errors.New(RegionNameExist)
	}

	vars, _ := json.Marshal(creation.RegionVars)
	region := model.Region{
		BaseModel:  common.BaseModel{},
		Name:       creation.Name,
		Vars:       string(vars),
		Datacenter: creation.Datacenter,
		Provider:   creation.Provider,
	}

	err := r.regionRepo.Save(&region)
	if err != nil {
		return nil, err
	}
	return &dto.Region{Region: region}, err
}

func (r regionService) Update(name string, update dto.RegionUpdate) (*dto.Region, error) {
	var region model.Region
	if err := db.DB.Where("name = ?", name).First(&region).Error; err != nil {
		return nil, err
	}
	vars, _ := json.Marshal(update.RegionVars)
	region.Vars = string(vars)
	region.Datacenter = update.Datacenter
	db.DB.Save(&region)
	return &dto.Region{Region: region}, nil
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

func (r regionService) ListDatacenter(creation dto.RegionDatacenterRequest) ([]string, error) {
	cloudClient := cloud_provider.NewCloudClient(creation.RegionVars.(map[string]interface{}))
	var result []string
	if cloudClient != nil {
		result, err := cloudClient.ListDatacenter()
		if err != nil {
			return result, err
		}
		return result, err
	}
	return result, nil
}
