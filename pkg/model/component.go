package model

import "github.com/KubeOperator/KubeOperator/pkg/model/common"

type ComponentDic struct {
	common.BaseModel
	ID       string `json:"-"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Version  string `json:"version"`
	Describe string `json:"describe"`
}
