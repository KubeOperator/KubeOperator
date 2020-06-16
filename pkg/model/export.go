package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/user"
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
	Credential{},
	Host{},
	Volume{},
	user.User{},
	Demo{},
}
