package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type TemplateConfig struct {
	model.TemplateConfig
	ConfigVars map[string]interface{} `json:"config"`
}

type TemplateConfigCreate struct {
	Name   string      `json:"name" validate:"required"`
	Type   string      `json:"type" validate:"required"`
	Config interface{} `json:"config" validate:"required"`
}
