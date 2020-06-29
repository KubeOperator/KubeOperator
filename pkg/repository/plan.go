package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type PlanRepository interface {
	Get(name string) (model.Plan, error)
	List() ([]model.Plan, error)
	Page(num, size int) (int, []model.Plan, error)
	Save(plan *model.Plan) error
	Delete(name string) error
	Batch(operation string, items []model.Plan) error
}

func NewPlanRepository() PlanRepository {
	return &planRepository{}
}

type planRepository struct {
}

func (p planRepository) Get(name string) (model.Plan, error) {
	var plan model.Plan
	plan.Name = name
	if err := db.DB.Where(plan).First(&plan).Error; err != nil {
		return plan, err
	}
	//if err := db.DB.First(&host).Related(&host.Volumes).Error; err != nil {
	//	return host, err
	//}
	//if err := db.DB.First(&host).Related(&host.Credential).Error; err != nil {
	//	return host, err
	//}
	return plan, nil
}

func (p planRepository) List() ([]model.Plan, error) {
	var plans []model.Plan
	err := db.DB.Model(model.Zone{}).Find(&plans).Error
	return plans, err
}

func (p planRepository) Page(num, size int) (int, []model.Plan, error) {
	var total int
	var plans []model.Plan
	err := db.DB.Model(model.Plan{}).
		Count(&total).
		Find(&plans).
		Offset((num - 1) * size).
		Limit(size).
		Error
	return total, plans, err
}

func (p planRepository) Save(plan *model.Plan) error {
	if db.DB.NewRecord(plan) {
		return db.DB.Create(&plan).Error
	} else {
		return db.DB.Save(&plan).Error
	}
}

func (p planRepository) Delete(name string) error {
	plan, err := p.Get(name)
	if err != nil {
		return err
	}
	return db.DB.Delete(&plan).Error
}

func (p planRepository) Batch(operation string, items []model.Plan) error {
	switch operation {
	case constant.BatchOperationDelete:
		//TODO 关联校验
		//var clusterIds []string
		//for _, item := range items {
		//	clusterIds = append(clusterIds, item.ClusterID)
		//}
		//var clusters []model.Cluster
		//err := db.DB.Where("id in (?)", clusterIds).Find(&clusters).Error
		//if err != nil {
		//	return err
		//}
		//if len(clusters) > 0 {
		//	return errors.New(DeleteFailedError)
		//}
		var ids []string
		for _, item := range items {
			ids = append(ids, item.ID)
		}
		err := db.DB.Where("id in (?)", ids).Delete(&items).Error
		if err != nil {
			return err
		}
	default:
		return constant.NotSupportedBatchOperation
	}
	return nil
}
