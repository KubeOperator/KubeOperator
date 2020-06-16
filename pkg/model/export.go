package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/credential"
	"github.com/KubeOperator/KubeOperator/pkg/model/host"
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
	credential.Credential{},
	host.Host{},
	host.Volume{},
	user.User{},
	Demo{},
}
