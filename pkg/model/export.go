package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/model/credential"
	"github.com/KubeOperator/KubeOperator/pkg/model/host"
)

type Interface interface {
	TableName() string
}

var Models = []Interface{
	cluster.Cluster{},
	cluster.Condition{},
	cluster.Label{},
	cluster.Spec{},
	cluster.Status{},
	credential.Credential{},
	host.Host{},
}
