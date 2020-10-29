package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type F5Repository interface {
	Save(f5 *model.F5Setting) error
	Get(name string) ([]model.F5Setting, error)
}

func NewF5Repository() F5Repository {
	return &f5Repository{}
}

type f5Repository struct {
}

func (f f5Repository) Save(f5 *model.F5Setting) error {
	if db.DB.NewRecord(f5) {
		return db.DB.Create(&f5).Error
	} else {
		return db.DB.Model(&f5).Update(&f5).Error
	}
}

func (f f5Repository) Get(clusterName string) ([]model.F5Setting, error) {
	var cluster model.Cluster
	if err := db.DB.Where(model.Cluster{Name: clusterName}).First(&cluster).Error; err != nil {
		return nil, err
	}
	var f5 []model.F5Setting
	if err := db.DB.
		Where(model.F5Setting{ClusterID: cluster.ID}).
		Find(&f5).Error; err != nil {
		return nil, err
	}
	return f5, nil
}
