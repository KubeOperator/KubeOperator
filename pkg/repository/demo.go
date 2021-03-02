package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type DemoRepository interface {
	List() ([]model.Demo, error)
	Get(name string) (model.Demo, error)
	Save(demo model.Demo) error
}

type demoRepository struct{}

func NewDemoRepository() *demoRepository {
	return &demoRepository{}
}

func (d *demoRepository) List() ([]model.Demo, error) {
	var results []model.Demo
	err := db.DB.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (d *demoRepository) Get(name string) (model.Demo, error) {
	var demo model.Demo
	if err := db.DB.Where("name = ?", name).First(&demo).Error; err != nil {
		return demo, err
	}
	return demo, nil
}

func (d *demoRepository) Save(demo model.Demo) error {
	if db.DB.NewRecord(demo) {
		if err := db.DB.Create(&demo).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(&demo).Error; err != nil {
			return err
		}
	}
	return nil
}
