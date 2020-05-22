package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/kobe/api"
)

type Cluster struct {
	common.BaseModel
	ID     string
	Name   string
	Spec   Spec
	Status Status
	Nodes  []Node
}

func (c Cluster) TableName() string {
	return "ko_cluster"
}

func (c Cluster) ParseInventory() api.Inventory {
	var masters []string
	var workers []string
	var hosts []*api.Host
	for _, node := range c.Nodes {
		hosts = append(hosts, node.ToKobeHost())
		switch node.LabelValue(constant.NodeRoleLabelKey) {
		case constant.NodeRoleNameMaster:
			masters = append(masters, node.Name)
		case constant.NodeRoleNameWorker:
			workers = append(workers, node.Name)
		}
	}
	return api.Inventory{
		Hosts: hosts,
		Groups: []*api.Group{
			{
				Name:     constant.NodeRoleNameMaster,
				Children: masters,
			}, {
				Name:     constant.NodeRoleNameWorker,
				Children: workers,
			},
		},
	}
}
