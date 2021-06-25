package model

import (
	"fmt"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
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
	Architecture       string
	Protocol           string
	Hostname           string
	RepoPort           int
	RegistryPort       int
	RegistryHostedPort int
}

func (n *ClusterNode) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}

func (n ClusterNode) GetRegistry(arch string) (*Registry, error) {
	var systemRegistry SystemRegistry
	var registry Registry
	err := db.DB.Where("architecture = ?", arch).First(&systemRegistry).Error
	if err != nil {
		return &registry, err
	}
	registry.Hostname = systemRegistry.Hostname
	registry.Protocol = systemRegistry.Protocol
	registry.RepoPort = systemRegistry.RepoPort
	registry.RegistryPort = systemRegistry.RegistryPort
	registry.RegistryHostedPort = systemRegistry.RegistryHostedPort
	switch n.Host.Architecture {
	case "x86_64":
		registry.Architecture = "amd64"
	case "aarch64":
		registry.Architecture = "arm64"
	}
	return &registry, nil
}

func (n ClusterNode) ToKobeHost() *api.Host {
	if err := n.Host.GetHostConfig(); err != nil {
		logger.Log.Errorf("get host config err, err: %s", err.Error())
	}
	if err := db.DB.Model(&Host{}).Where("id = ?", n.Host.ID).Updates(map[string]interface{}{
		"architecture": n.Host.Architecture,
		"os":           n.Host.Os,
		"os_version":   n.Host.OsVersion}).Error; err != nil {
		logger.Log.Errorf("get host config err, err: %s", err.Error())
	}

	r, _ := n.GetRegistry(n.Host.Architecture)
	return &api.Host{
		Ip:         n.Host.Ip,
		Name:       n.Name,
		Port:       int32(n.Host.Port),
		User:       n.Host.Credential.Username,
		Password:   n.Host.Credential.Password,
		PrivateKey: n.Host.Credential.PrivateKey,
		Vars: map[string]string{
			"has_gpu":              fmt.Sprintf("%v", n.Host.HasGpu),
			"architectures":        r.Architecture,
			"registry_protocol":    r.Protocol,
			"registry_hostname":    r.Hostname,
			"repo_port":            fmt.Sprintf("%v", r.RepoPort),
			"registry_port":        fmt.Sprintf("%v", r.RegistryPort),
			"registry_hosted_port": fmt.Sprintf("%v", r.RegistryHostedPort),
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
