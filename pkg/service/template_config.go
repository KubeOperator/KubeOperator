package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
)

type TemplateConfigService interface {
	List() ([]dto.TemplateConfig, error)
}

type templateConfigService struct {
}

func NewTemplateConfigService() TemplateConfigService {
	return &templateConfigService{}
}

func (t *templateConfigService) List() ([]dto.TemplateConfig, error) {
	var configs []dto.TemplateConfig
	err := db.DB.Find(&configs).Error
	return configs, err
}
