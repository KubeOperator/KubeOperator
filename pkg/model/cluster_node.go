package model

import (
	"fmt"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	_ "github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
)

type ClusterNode struct {
	common.BaseModel
	ID        string `json:"-"`
	Name      string `json:"name"`
	HostID    string `json:"-"`
	Host      Host   `json:"-" gorm:"save_associations:false"`
	ClusterID string `json:"clusterId"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	Dirty     bool   `json:"dirty"`
	Message   string `json:"message"`
}

type Registry struct {
	Architecture string
	Protocol     string
	Hostname     string
}

func (n *ClusterNode) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}

func getSetting(key string) (string, error) {
	var systemSetting SystemSetting
	if err := db.DB.Where(map[string]interface{}{"key": key}).First(&systemSetting).Error; err != nil {
		return systemSetting.Value, err
	}
	return systemSetting.Value, nil
}

func (n ClusterNode) GetRegistry(arch string) (*Registry, error) {
	var systemRegistry SystemRegistry
	var registry Registry

	archType, err := getSetting("arch_type")
	if err != nil {
		return nil, err
	}
	if archType == "single" {
		registry.Hostname, err = getSetting("ip")
		if err != nil {
			return nil, err
		}
		registry.Protocol, err = getSetting("REGISTRY_PROTOCOL")
		if err != nil {
			return nil, err
		}
		switch n.Host.Architecture {
		case "x86_64":
			registry.Architecture = "amd64"
		case "aarch64":
			registry.Architecture = "arm64"
		default:
			registry.Architecture = "amd64"
		}
	} else if archType == "mixed" {
		err := db.DB.Where("architecture = ?", arch).First(&systemRegistry).Error
		if err != nil {
			return nil, err
		}
		registry.Hostname = systemRegistry.Hostname
		registry.Protocol = systemRegistry.Protocol
		switch n.Host.Architecture {
		case "x86_64":
			registry.Architecture = "amd64"
		case "aarch64":
			registry.Architecture = "arm64"
		}
	}
	return &registry, nil
}

func (n ClusterNode) ToKobeHost() *api.Host {
	password, privateKey, _ := n.Host.GetHostPasswordAndPrivateKey()
	r, _ := n.GetRegistry(n.Host.Architecture)
	return &api.Host{
		Ip:         n.Host.Ip,
		Name:       n.Name,
		Port:       int32(n.Host.Port),
		User:       n.Host.Credential.Username,
		Password:   password,
		PrivateKey: string(privateKey),
		Vars: map[string]string{
			"has_gpu":           fmt.Sprintf("%v", n.Host.HasGpu),
			"architecture":      r.Architecture,
			"registry_protocol": r.Protocol,
			"registry_hostname": r.Hostname,
		},
	}
}

func (n ClusterNode) ToSSHConfig() ssh.Config {
	password, privateKey, _ := n.Host.GetHostPasswordAndPrivateKey()
	return ssh.Config{
		User:        n.Host.Credential.Username,
		Host:        n.Host.Ip,
		Port:        n.Host.Port,
		PrivateKey:  privateKey,
		Password:    password,
		DialTimeOut: 5 * time.Second,
		Retry:       3,
	}
}
