package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"time"
)

type Interface interface {
	TableName() string
}

var Models = []Interface{
	Cluster{},
	ClusterStatusCondition{},
	ClusterSpec{},
	ClusterStatus{},
	ClusterNode{},
	ClusterSecret{},
	ClusterTool{},
	ClusterStorageProvisioner{},
	Credential{},
	Host{},
	Volume{},
	User{},
	Demo{},
	Region{},
	Zone{},
	Plan{},
	PlanZones{},
	SystemSetting{},
	Project{},
	ProjectResource{},
	ProjectMember{},
}

var InitData = []Interface{
	User{
		BaseModel: common.BaseModel{
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
		ID:       "5e81095f-3c0c-4cb2-8033-bde03d60135c",
		Name:     "admin",
		Password: "47zHCOqO84rdzGgxw5XPfgDEapoOMXbgJnryG32xp6Y=",
		Email:    "admin@fit2cloud.com",
		Language: ZH,
		IsActive: true,
		IsAdmin:  true,
	},
	Credential{
		BaseModel: common.BaseModel{
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
		ID:       "f081498c-7c00-4955-8181-884f93088dc4",
		Name:     constant.ImageCredentialName,
		Password: "QK6fxpxyb/qf8Ssr2ShZeF//savV3zdtmcOS6FPd3yQ=",
		Username: constant.ImageUserName,
		Type:     constant.ImagePasswordType,
	},
}
