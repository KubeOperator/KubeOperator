package service

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
)

type TemplateConfigService interface {
	List() ([]dto.TemplateConfig, error)
	Page(num, size int, conditions condition.Conditions) (*page.Page, error)
	Create(creation dto.TemplateConfigCreate) (*dto.TemplateConfig, error)
	Get(name string) (*dto.TemplateConfig, error)
	Delete(name string) error
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

func (t *templateConfigService) Get(name string) (*dto.TemplateConfig, error) {
	var (
		mo     model.TemplateConfig
		config dto.TemplateConfig
	)
	if err := db.DB.Where("name = ?", name).First(&mo).Error; err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	config.TemplateConfig = mo
	if err := json.Unmarshal([]byte(mo.Config), &m); err != nil {
		logger.Log.Errorf("templateConfigService Get json.Unmarshal failed, error: %s", err.Error())
	}
	config.ConfigVars = m
	return &config, nil
}

func (t *templateConfigService) Create(creation dto.TemplateConfigCreate) (*dto.TemplateConfig, error) {

	old, _ := t.Get(creation.Name)
	if old != nil && old.ID != "" {
		return nil, errors.New("NAME_EXISTS")
	}
	config, _ := json.Marshal(creation.Config)
	mo := model.TemplateConfig{
		Config:    string(config),
		BaseModel: common.BaseModel{},
		Name:      creation.Name,
		Type:      creation.Type,
	}

	return &dto.TemplateConfig{TemplateConfig: mo}, db.DB.Create(&mo).Error
}

func (t *templateConfigService) Delete(name string) error {
	return db.DB.Where("name = ?", name).Delete(&model.TemplateConfig{}).Error
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
