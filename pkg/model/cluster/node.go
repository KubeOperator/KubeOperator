package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/model/host"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
)

type Node struct {
	common.BaseModel
	ID        string
	Name      string    `gorm:"not null;unique"`
	Host      host.Host `gorm:"save_associations:false"`
	HostID    string
	ClusterID string
	Role      string
}

func (n *Node) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}
func (n Node) ToKobeHost() *api.Host {
	return &api.Host{
		Ip:       n.Host.Ip,
		Name:     n.Name,
		Port:     int32(n.Host.Port),
		User:     n.Host.Credential.Username,
		Password: n.Host.Credential.Password,
	}
}

func (n Node) TableName() string {
	return "ko_cluster_node"
}
