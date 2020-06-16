package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
	"time"
)

type ClusterNode struct {
	common.BaseModel
	ID        string     `json:"_"`
	Name      string     `gorm:"not null;unique" json:"name"`
	Host      Host `gorm:"save_associations:false"`
	ClusterID string     `json:"clusterId"`
	Role      string     `json:"role"`
}

func (n *ClusterNode) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}

func (n ClusterNode) AfterDelete() error {
	var host hostModel.Host
	if err := db.DB.Where(hostModel.Host{
		NodeID: n.ID,
	}).First(&host).Error; err != nil {
		return err
	}
	host.NodeID = ""
	if err := db.DB.Save(&host).Error; err != nil {
		return err
	}
	return nil
}

func (n ClusterNode) ToKobeHost() *api.Host {
	password, _, _ := hostService.GetHostPasswordAndPrivateKey(&n.Host)
	return &api.Host{
		Ip:       n.Host.Ip,
		Name:     n.Name,
		Port:     int32(n.Host.Port),
		User:     n.Host.Credential.Username,
		Password: password,
	}
}

func (n ClusterNode) ToSSHConfig() ssh.Config {
	password, _, _ := hostService.GetHostPasswordAndPrivateKey(&n.Host)
	return ssh.Config{
		User:        n.Host.Credential.Username,
		Host:        n.Host.Ip,
		Port:        n.Host.Port,
		Password:    password,
		DialTimeOut: 5 * time.Second,
		Retry:       3,
	}
}
func (n ClusterNode) TableName() string {
	return "ko_cluster_node"
}
