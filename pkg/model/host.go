package model

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

const (
	Disconnect string = "DisConnect"
	SshError   string = "SshError"
)

type Host struct {
	common.BaseModel
	ID           string     `json:"-"`
	Name         string     `json:"name" gorm:"type:varchar(256);not null;unique"`
	Memory       int        `json:"memory" gorm:"type:int(64)"`
	CpuCore      int        `json:"cpuCore" gorm:"type:int(64)"`
	Os           string     `json:"os" gorm:"type:varchar(64)"`
	OsVersion    string     `json:"osVersion" gorm:"type:varchar(64)"`
	GpuNum       int        `json:"gpuNum" gorm:"type:int(64)"`
	GpuInfo      string     `json:"gpuInfo" gorm:"type:varchar(128)"`
	Ip           string     `json:"ip" gorm:"type:varchar(128);not null;unique"`
	HasGpu       bool       `json:"hasGpu" gorm:"type:boolean;default:false"`
	Port         int        `json:"port" gorm:"type:varchar(64)"`
	CredentialID string     `json:"credentialId" gorm:"type:varchar(64)"`
	ClusterID    string     `json:"clusterId" gorm:"type:varchar(64)"`
	ZoneID       string     `json:"zoneId" gorm:"type:varchar(64)"`
	Zone         Zone       `json:"-"  gorm:"save_associations:false" `
	Volumes      []Volume   `json:"volumes" gorm:"save_associations:false"`
	Credential   Credential `json:"-" gorm:"save_associations:false" `
	Cluster      Cluster    `json:"-" gorm:"save_associations:false" `
	Status       string     `json:"status" gorm:"type:varchar(64)"`
	Message      string     `json:"message" gorm:"type:text(65535)"`
	Datastore    string     `json:"datastore" gorm:"type:varchar(64)"`
}

func (h Host) GetHostPasswordAndPrivateKey() (string, []byte, error) {
	password := ""
	privateKey := []byte("")
	switch h.Credential.Type {
	case "password":
		p, err := encrypt.StringDecrypt(h.Credential.Password)
		if err != nil {
			return "", nil, err
		}
		password = p
	case "privateKey":
		privateKey = []byte(h.Credential.PrivateKey)
	}
	return password, privateKey, nil
}

func (h *Host) BeforeCreate() error {
	h.ID = uuid.NewV4().String()
	return nil
}

func (h *Host) BeforeDelete(tx *gorm.DB) error {
	if h.ID != "" {
		if h.ClusterID != "" {
			var cluster Cluster
			cluster.ID = h.ClusterID
			notFound := tx.First(&cluster).RecordNotFound()
			if !notFound {
				return errors.New("DELETE_HOST_FAILED")
			}
		}
		var projectResources []ProjectResource
		err := tx.Where(ProjectResource{ResourceID: h.ID}).Find(&projectResources).Error
		if err != nil {
			return err
		}
		if len(projectResources) > 0 {
			return errors.New("DELETE_HOST_FAILED_BY_PROJECT")
		}
		var ip Ip
		tx.Where(Ip{Address: h.Ip}).First(&ip)
		if ip.ID != "" {
			ip.Status = constant.IpAvailable
			if err := tx.Save(&ip).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}
