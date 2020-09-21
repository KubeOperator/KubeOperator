package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type UserNotificationConfig struct {
	model.UserNotificationConfig
}

type UserNotificationConfigDTO struct {
	ID     string            `json:"id"`
	UserID string            `json:"userId"`
	Vars   map[string]string `json:"vars"`
	Type   string            `json:"type"`
}
