package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
	hostService "github.com/KubeOperator/KubeOperator/pkg/service/host"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
	"log"
)

type Node struct {
	common.BaseModel
	ID        string
	Name      string         `gorm:"not null;unique"`
	Host      hostModel.Host `gorm:"save_associations:false"`
	HostID    string
	ClusterID string
	Role      string
}

func (n *Node) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}
func (n Node) ToKobeHost() *api.Host {
	password, _, err := hostService.GetHostPasswordAndPrivateKey(&n.Host)
	if err != nil {
		log.Println(err)
	}
	return &api.Host{
		Ip:       n.Host.Ip,
		Name:     n.Name,
		Port:     int32(n.Host.Port),
		User:     n.Host.Credential.Username,
		Password: password,
	}
}

func (n Node) TableName() string {
	return "ko_cluster_node"
}
