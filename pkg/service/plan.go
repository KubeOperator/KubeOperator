package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider/client"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/mitchellh/mapstructure"
	"sort"
)

type PlanService interface {
	Get(name string) (dto.Plan, error)
	List() ([]dto.Plan, error)
	Page(num, size int) (page.Page, error)
	Delete(name string) error
	Create(creation dto.PlanCreate) (dto.Plan, error)
	Batch(op dto.PlanOp) error
	GetConfigs(regionName string) ([]dto.PlanVmConfig, error)
}

type planService struct {
	planRepo repository.PlanRepository
}

func NewPlanService() PlanService {
	return &planService{
		planRepo: repository.NewPlanRepository(),
	}
}

func (p planService) Get(name string) (dto.Plan, error) {
	var planDTO dto.Plan
	mo, err := p.planRepo.Get(name)
	if err != nil {
		return planDTO, err
	}
	planDTO.Plan = mo
	return planDTO, err
}

func (p planService) List() ([]dto.Plan, error) {
	var planDTOs []dto.Plan
	mos, err := p.planRepo.List()
	if err != nil {
		return planDTOs, err
	}
	for _, mo := range mos {
		planDTOs = append(planDTOs, dto.Plan{Plan: mo})
	}
	return planDTOs, err
}

func (p planService) Page(num, size int) (page.Page, error) {
	var page page.Page
	var planDTOs []dto.Plan
	total, mos, err := p.planRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {

		planDTO := new(dto.Plan)
		r := make(map[string]interface{})
		json.Unmarshal([]byte(mo.Vars), &r)
		planDTO.PlanVars = r
		planDTO.Plan = mo

		planDTOs = append(planDTOs, *planDTO)
	}
	page.Total = total
	page.Items = planDTOs
	return page, err
}

func (p planService) Delete(name string) error {
	err := p.planRepo.Delete(name)
	if err != nil {
		return err
	}
	return nil
}

func (p planService) Create(creation dto.PlanCreate) (dto.Plan, error) {

	vars, _ := json.Marshal(creation.PlanVars)

	plan := model.Plan{
		BaseModel:      common.BaseModel{},
		Name:           creation.Name,
		Vars:           string(vars),
		RegionID:       creation.RegionId,
		DeployTemplate: creation.DeployTemplate,
	}

	err := p.planRepo.Save(&plan, creation.Zones)
	if err != nil {
		return dto.Plan{}, err
	}
	return dto.Plan{Plan: plan}, err
}

func (p planService) Batch(op dto.PlanOp) error {
	var deleteItems []model.Plan
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.Plan{
			BaseModel: common.BaseModel{},
			ID:        item.ID,
			Name:      item.Name,
		})
	}
	err := p.planRepo.Batch(op.Operation, deleteItems)
	if err != nil {
		return err
	}
	return nil
}

func (p planService) GetConfigs(regionName string) ([]dto.PlanVmConfig, error) {
	region, err := NewRegionService().Get(regionName)
	if err != nil {
		return nil, err
	}
	var configs []dto.PlanVmConfig
	if region.Provider == "OpenStack" {
		vars := region.RegionVars.(map[string]interface{})
		vars["datacenter"] = region.Datacenter
		cloudClient := client.NewCloudClient(vars)
		result, err := cloudClient.ListFlavors()
		if err != nil {
			return nil, err
		}
		err = mapstructure.Decode(result, &configs)
		if err != nil {
			return nil, err
		}
	} else {
		for k, v := range constant.VmConfigList {
			configs = append(configs, dto.PlanVmConfig{
				Name:   k,
				Config: v,
			})
		}
		sort.Slice(configs, func(i, j int) bool {
			return configs[i].Config.Cpu < configs[j].Config.Cpu
		})
	}
	return configs, nil
}
