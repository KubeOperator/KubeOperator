package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Spec struct {
	Version     string
	NetworkType string
	ClusterCIDR string
	ServiceCIDR string
	Nodes       []Node
}

type Condition struct {
	Name          string
	Status        string
	Message       string
	LastProbeTime time.Time
}

type Status struct {
	Version    string
	Message    string
	Phase      string
	Conditions []Condition
}

type Cluster struct {
	common.BaseModel
	Spec   Spec
	Status Status
}

func (c *Cluster) BeforeCreate() error {
	c.ID = uuid.NewV4().String()
	c.CreatedDate = time.Now()
	c.UpdatedDate = time.Now()
	return nil
}

func (c *Cluster) BeforeUpdate() error {
	c.UpdatedDate = time.Now()
	return nil
}

func (c Cluster) TableName() string {
	return "ko_cluster"
}

func (c Cluster) ParseInventory() api.Inventory {
	var masters []string
	var workers []string
	var hosts []*api.Host
	for _, node := range c.Spec.Nodes {
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
