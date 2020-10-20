package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type F5Repository interface {
	Save(plan *model.F5Setting) error
	//Get(name string) (model.F5Setting,error)
}

func NewF5Repository() F5Repository {
	return &f5Repository{}
}

type f5Repository struct {
}

func (f f5Repository) Save(f5 *model.F5Setting) error {
	if db.DB.NewRecord(f5) {
		tx := db.DB.Begin()
		err := tx.Create(&f5).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return err
	} else {
		tx := db.DB.Begin()
		err := tx.Create(&f5).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Where("id = ?", f5.ID).Delete(&model.F5Setting{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return err
	}
}

// 未完成
//func (f f5Repository) Get(name *model.F5Setting) (model.F5Setting,error)  {
//	var f5 model.F5Setting
//	if err  := db.DB.Where(f5).First(&f5).Error; err != nil {
//		return f5,err
//	}
//	return f5,nil
//}
