package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type SystemLog struct {
	model.SystemLog
}

type SystemLogCreate struct {
	Name          string `json:"name" gorm:"type:varchar(256);not null;"`
	Operation     string `json:"operation" gorm:"type:varchar(256);not null;"`
	OperationInfo string `json:"operationInfo" gorm:"type:varchar(256);"`
	IP            string `json:"ip" gorm:"type:varchar(20);"`
}
