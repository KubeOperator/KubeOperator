package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/model/host"
	"github.com/KubeOperator/kobe/api"
)

type Node struct {
	common.BaseModel
	ID     string
	Name   string
	Host   host.Host
	HostID string
	Labels []Label
}

func (n Node) LabelValue(name string) string {
	result := ""
	for _, label := range n.Labels {
		if n.Name == name {
			result = label.Value
		}
	}
	return result
}

func (n Node) ToKobeHost() *api.Host {
	return &api.Host{
		Ip:       n.Host.Ip,
		Name:     n.Host.Name,
		Port:     int32(n.Host.Port),
		User:     n.Host.Credential.Username,
		Password: n.Host.Credential.Password,
	}
}

func (n Node) TableName() string {
	return "ko_cluster_node"
}
