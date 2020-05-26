package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
)

type Cluster struct {
	common.BaseModel
	ID         string
	Name       string      `gorm:"not null;unique"`
	Spec       Spec        `gorm:"save_associations:false"`
	Status     Status      `gorm:"save_associations:false"`
	Nodes      []Node      `gorm:"save_associations:false"`
	Conditions []Condition `gorm:"save_associations:false"`
}

func (c *Cluster) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return nil
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
		switch node.Role {
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
