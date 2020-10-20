package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type F5Setting struct {
	model.F5Setting
	ClusterName string `json:"clusterId"`
}

type F5SettingCreate struct {
	Vars map[string]string `json:"vars" validate:"required"`
}

type F5SettingUpdate struct {
	Vars map[string]string `json:"vars" validate:"required"`
}

type F5SettingResult struct {
	Vars map[string]string `json:"vars" validate:"required"`
}
