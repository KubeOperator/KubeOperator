package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type TemplateConfig struct {
	model.TemplateConfig
	ConfigVars map[string]interface{}
}
