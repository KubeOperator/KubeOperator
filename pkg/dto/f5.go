package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type F5Setting struct {
	model.F5Setting
	Password    string `json:"password" gorm:"type:varchar(64)"`
	ClusterName string `json:"clusterName"`
}

type F5SettingResult struct {
	model.F5Setting
	ClusterName string `json:"clusterName"`
}
