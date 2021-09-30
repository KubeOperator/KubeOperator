package service

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"sort"

	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
	"github.com/mitchellh/mapstructure"
)

var (
	PlanNameExist = "NAME_EXISTS"
)

type PlanService interface {
	Get(name string) (*dto.Plan, error)
	List(projectName string) ([]dto.Plan, error)
	Page(num, size int, projectName string, conditions condition.Conditions) (*page.Page, error)
	Delete(name string) error
	Create(creation dto.PlanCreate) (*dto.Plan, error)
	Batch(op dto.PlanOp) error
	GetConfigs(regionName string) ([]dto.PlanVmConfig, error)
	PatchBy(name string, update dto.PlanUpdate) (*dto.Plan, error)
}

type planService struct {
	planRepo            repository.PlanRepository
	regionRepo          repository.RegionRepository
	projectResourceRepo repository.ProjectResourceRepository
	projectRepo         repository.ProjectRepository
	vmConfigRepo        repository.VmConfigRepository
}

func NewPlanService() PlanService {
	return &planService{
		planRepo:            repository.NewPlanRepository(),
		regionRepo:          repository.NewRegionRepository(),
		projectResourceRepo: repository.NewProjectResourceRepository(),
		projectRepo:         repository.NewProjectRepository(),
		vmConfigRepo:        repository.NewVmConfigRepository(),
	}
}

func (p planService) Get(name string) (*dto.Plan, error) {
	var (
		planDTO          dto.Plan
		plan             model.Plan
		projectResources []model.ProjectResource
	)
	if err := db.DB.Where("name = ?", name).
		Preload("Zones").
		Preload("Region").First(&plan).Error; err != nil {
		return nil, err
	}
	r := make(map[string]interface{})
	if err := json.Unmarshal([]byte(plan.Vars), &r); err != nil {
		return nil, err
	}
	for _, zone := range plan.Zones {
		planDTO.Zones = append(planDTO.Zones, zone.Name)
	}
	planDTO.PlanVars = r
	planDTO.Plan = plan
	planDTO.Region = plan.Region.Name
	planDTO.Provider = plan.Region.Provider
	if err := db.DB.Where("resource_id = ?", plan.ID).Preload("Project").Find(&projectResources).Error; err != nil {
		return nil, err
	}
	for _, pr := range projectResources {
		planDTO.Projects = append(planDTO.Projects, pr.Project.Name)
	}
	return &planDTO, nil
}

func (p planService) List(projectName string) ([]dto.Plan, error) {
	var planDTOs []dto.Plan
	mos, err := p.planRepo.List(projectName)
	if err != nil {
		return planDTOs, err
	}
	for _, mo := range mos {
		planDTOs = append(planDTOs, dto.Plan{Plan: mo})
	}
	return planDTOs, err
}

func (p planService) Page(num, size int, projectName string, conditions condition.Conditions) (*page.Page, error) {

	var (
		pa       page.Page
		planDTOs []dto.Plan
		plans    []model.Plan
	)

	d := db.DB.Model(model.Plan{})
	if err := dbUtil.WithConditions(&d, model.Plan{}, conditions); err != nil {
		return nil, err
	}

	if len(projectName) != 0 {
		if _, err := dbUtil.WithProjectResource(&d, projectName, constant.ResourcePlan); err != nil {
			return nil, err
		}
	}

	if err := d.Preload("Region").Preload("Zones").Order("name asc").Count(&pa.Total).Offset((num - 1) * size).Limit(size).Find(&plans).Error; err != nil {
		return nil, err
	}
	for _, p := range plans {

		planDTO := new(dto.Plan)
		r := make(map[string]interface{})
		if err := json.Unmarshal([]byte(p.Vars), &r); err != nil {
			return nil, err
		}
		var zoneNames []string
		for _, zone := range p.Zones {
			zoneNames = append(zoneNames, zone.Name)
		}
		planDTO.PlanVars = r
		planDTO.Plan = p
		planDTO.Region = p.Region.Name
		planDTO.Zones = zoneNames
		var projectResources []model.ProjectResource
		if err := db.DB.Where("resource_id = ?", p.ID).Preload("Project").Find(&projectResources).Error; err != nil {
			return nil, err
		}
		for _, pr := range projectResources {
			planDTO.Projects = append(planDTO.Projects, pr.Project.Name)
		}
		planDTOs = append(planDTOs, *planDTO)
	}
	pa.Items = planDTOs
	return &pa, nil
}

func (p planService) Delete(name string) error {
	err := p.planRepo.Delete(name)
	if err != nil {
		return err
	}
	return nil
}

func (p planService) Create(creation dto.PlanCreate) (*dto.Plan, error) {

	var old model.Plan
	if err := db.DB.Where("name = ?", creation.Name).Find(&old).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if old.ID != "" {
		return nil, errors.New(PlanNameExist)
	}
	vars, _ := json.Marshal(creation.PlanVars)
	var region model.Region
	if err := db.DB.Where("name = ?", creation.Region).First(&region).Error; err != nil {
		return nil, err
	}
	tx := db.DB.Begin()
	plan := model.Plan{
		BaseModel:      common.BaseModel{},
		Name:           creation.Name,
		Vars:           string(vars),
		RegionID:       region.ID,
		DeployTemplate: creation.DeployTemplate,
	}
	err := tx.Create(&plan).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var zones []model.Zone
	err = tx.Where("name in (?)", creation.Zones).Find(&zones).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, z := range zones {
		err = tx.Create(&model.PlanZones{
			PlanID: plan.ID,
			ZoneID: z.ID,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	var projects []model.Project
	err = tx.Where("name in (?)", creation.Projects).Find(&projects).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, project := range projects {
		err = tx.Create(&model.ProjectResource{
			ResourceType: constant.ResourcePlan,
			ResourceID:   plan.ID,
			ProjectID:    project.ID,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &dto.Plan{Plan: plan}, err
}

func (p planService) PatchBy(name string, update dto.PlanUpdate) (*dto.Plan, error) {
	var plan model.Plan
	if err := db.DB.Where("name = ?", name).Find(&plan).Error; err != nil {
		return nil, err
	}
	vars, _ := json.Marshal(update.PlanVars)
	plan.Vars = string(vars)
	var projects []model.Project
	tx := db.DB.Begin()
	if err := tx.Where("name in (?)", update.Projects).Find(&projects).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Model(model.ProjectResource{}).Where("resource_id = ?", plan.ID).Delete(&model.ProjectResource{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, project := range projects {
		err := tx.Create(&model.ProjectResource{
			ResourceType: constant.ResourcePlan,
			ResourceID:   plan.ID,
			ProjectID:    project.ID,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Save(&plan)
	tx.Commit()
	return &dto.Plan{Plan: plan}, nil
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
	if region.Provider == constant.OpenStack {
		vars := region.RegionVars.(map[string]interface{})
		vars["datacenter"] = region.Datacenter
		cloudClient := cloud_provider.NewCloudClient(vars)
		result, err := cloudClient.ListFlavors()
		if err != nil {
			return nil, err
		}
		err = mapstructure.Decode(result, &configs)
		if err != nil {
			return nil, err
		}
	} else {
		vmConfigs, err := p.vmConfigRepo.List()
		if err != nil {
			return nil, err
		}
		for _, config := range vmConfigs {
			configs = append(configs, dto.PlanVmConfig{
				Name: config.Name,
				Config: constant.VmConfig{
					Cpu:    config.Cpu,
					Memory: config.Memory,
					Disk:   config.Disk,
				},
			})
		}
		sort.Slice(configs, func(i, j int) bool {
			return configs[i].Config.Cpu < configs[j].Config.Cpu
		})
	}
	return configs, nil
}
