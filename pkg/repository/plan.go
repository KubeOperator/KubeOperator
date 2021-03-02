package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type PlanRepository interface {
	Get(name string) (model.Plan, error)
	List(projectName string) ([]model.Plan, error)
	Page(num, size int) (int, []model.Plan, error)
	Save(plan *model.Plan, zones []string) error
	Delete(name string) error
	Batch(operation string, items []model.Plan) error

	GetById(id string) (model.Plan, error)
}

func NewPlanRepository() PlanRepository {
	return &planRepository{}
}

type planRepository struct {
}

func (p planRepository) Get(name string) (model.Plan, error) {
	var plan model.Plan
	if err := db.DB.Where("name = ?", name).
		Preload("Zones").
		Preload("Region").First(&plan).Error; err != nil {
		return plan, err
	}
	return plan, nil
}

func (p planRepository) GetById(id string) (model.Plan, error) {
	var plan model.Plan
	plan.ID = id
	err := db.DB.First(&plan).
		Preload("Zones").
		Preload("Region").
		Find(&plan).Error
	if err != nil {
		return plan, err
	}
	return plan, nil
}

func (p planRepository) List(projectName string) ([]model.Plan, error) {
	var plans []model.Plan
	if projectName == "" {
		err := db.DB.Find(&plans).Error
		return plans, err
	} else {
		var project model.Project
		err := db.DB.Where("name = ?", projectName).First(&project).Error
		if err != nil {
			return nil, err
		}
		var projectResources []model.ProjectResource
		err = db.DB.Where("project_id = ? AND resource_type = ?", project.ID, constant.ResourcePlan).Find(&projectResources).Error
		if err != nil {
			return nil, err
		}
		var resourceIds []string
		for _, pr := range projectResources {
			resourceIds = append(resourceIds, pr.ResourceID)
		}
		err = db.DB.Where("id in (?)", resourceIds).Find(&plans).Error
		return plans, err
	}
}

func (p planRepository) Page(num, size int) (int, []model.Plan, error) {
	var total int
	var plans []model.Plan
	err := db.DB.Model(&model.Plan{}).Count(&total).Offset((num - 1) * size).Limit(size).Find(&plans).Error

	for i, p := range plans {
		var zoneIds []string
		var planZones []model.PlanZones
		db.DB.Where("plan_id = ?", p.ID).Find(&planZones)
		for _, p := range planZones {
			zoneIds = append(zoneIds, p.ZoneID)
		}
		var zones []model.Zone
		db.DB.Where("id in (?)", zoneIds).Find(&zones)
		plans[i].Zones = zones

		var region model.Region
		db.DB.Where("id = ?", p.RegionID).First(&region)
		plans[i].Region = region
	}

	return total, plans, err
}

func (p planRepository) Save(plan *model.Plan, zones []string) error {
	if db.DB.NewRecord(plan) {
		tx := db.DB.Begin()
		err := tx.Create(&plan).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, z := range zones {
			err = tx.Create(&model.PlanZones{
				PlanID: plan.ID,
				ZoneID: z,
			}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
		return err
	} else {
		tx := db.DB.Begin()
		err := db.DB.Save(&plan).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Where("plan_id = ?", plan.ID).Delete(&model.PlanZones{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, z := range zones {
			err = tx.Create(model.PlanZones{
				PlanID: plan.ID,
				ZoneID: z,
			}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
		return err
	}
}

func (p planRepository) Delete(name string) error {
	plan, err := p.Get(name)
	if err != nil {
		return err
	}
	tx := db.DB.Begin()
	err = tx.Delete(&plan).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Where("plan_id = ?", plan.ID).Delete(&model.PlanZones{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (p planRepository) Batch(operation string, items []model.Plan) error {
	switch operation {
	case constant.BatchOperationDelete:
		var ids []string
		for _, item := range items {
			ids = append(ids, item.ID)
		}

		tx := db.DB.Begin()
		err := tx.Where("id in (?)", ids).Delete(&items).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		var planZones []model.PlanZones
		err = tx.Where("plan_id in (?)", ids).Delete(&planZones).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()

	default:
		return constant.NotSupportedBatchOperation
	}
	return nil
}
