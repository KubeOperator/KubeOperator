package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/model/credential"
	"github.com/KubeOperator/KubeOperator/pkg/model/host"
	"github.com/KubeOperator/KubeOperator/pkg/model/user"
)

type Interface interface {
	TableName() string
}

var Models = []Interface{
	cluster.Cluster{},
	cluster.Condition{},
	cluster.Spec{},
	cluster.Status{},
	cluster.Node{},
	cluster.Secret{},
	cluster.Task{},
	credential.Credential{},
	host.Host{},
	host.Volume{},
	user.User{},
}
