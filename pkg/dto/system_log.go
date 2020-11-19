package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type SystemLog struct {
	model.SystemLog
}

type SystemLogCreate struct {
	Name          string `json:"name" gorm:"type:varchar(256);not null;"`
	OperationUnit string `json:"operationUnit" gorm:"type:varchar(256);not null;"`
	Operation     string `json:"operation" gorm:"type:varchar(256);not null;"`
	RequestPath   string `json:"requestPath" gorm:"type:varchar(256);not null;"`
}
