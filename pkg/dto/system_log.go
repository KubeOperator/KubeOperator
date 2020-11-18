package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type SystemLog struct {
	model.SystemLog
}

type SystemLogCreate struct {
	Name          string `json:"name" gorm:"type:varchar(256);not null;"`
	OperationUnit string `json:"operation_unit" gorm:"type:varchar(256);not null;"`
	Operation     string `json:"operation" gorm:"type:varchar(256);not null;"`
	RequestPath   string `json:"request_path" gorm:"type:varchar(256);not null;"`
}
