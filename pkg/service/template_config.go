package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
)

type TemplateConfigService interface {
	List() ([]dto.TemplateConfig, error)
	Page(num, size int, conditions condition.Conditions) (*page.Page, error)
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

func (t *templateConfigService) Create(config dto.TemplateConfig) (dto.TemplateConfig, error) {
	return config, db.DB.Create(config).Error
}

func (t *templateConfigService) Page(num, size int, conditions condition.Conditions) (*page.Page, error) {

	var (
		p            page.Page
		templateDTOs []dto.TemplateConfig
		templates    []model.TemplateConfig
	)

	d := db.DB.Model(model.TemplateConfig{})
	if err := dbUtil.WithConditions(&d, model.TemplateConfig{}, conditions); err != nil {
		return nil, err
	}
	if err := d.Order("created_at asc").Count(&p.Total).Offset((num - 1) * size).Limit(size).Find(&templates).Error; err != nil {
		return nil, err
	}
	for _, mo := range templates {
		templateDTO := new(dto.TemplateConfig)
		templateDTO.TemplateConfig = mo
		m := make(map[string]interface{})
		if err := json.Unmarshal([]byte(mo.Config), &m); err != nil {
			logger.Log.Errorf("regionService Page json.Unmarshal failed, error: %s", err.Error())
		}
		templateDTO.ConfigVars = m
		templateDTOs = append(templateDTOs, *templateDTO)
	}
	p.Items = templateDTOs
	return &p, nil
}
