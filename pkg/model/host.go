package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/kobe/api"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"k8s.io/apimachinery/pkg/util/wait"
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
	FlexIp       string     `json:"flexIp" gorm:"type:varchar(128);unique"`
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
	Architecture string     `json:"architecture" gorm:"type:varchar(64)"`
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

func (h *Host) GetHostConfig() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("gather fact error!")
		}
	}()

	password, privateKey, err := h.GetHostPasswordAndPrivateKey()
	h.Credential.Password = password
	h.Credential.PrivateKey = string(privateKey)
	if err != nil {
		return err
	}
	ansible := kobe.NewAnsible(&kobe.Config{
		Inventory: &api.Inventory{
			Hosts: []*api.Host{
				{
					Ip:         h.Ip,
					Name:       h.Name,
					Port:       int32(h.Port),
					User:       h.Credential.Username,
					Password:   password,
					PrivateKey: string(privateKey),
					Vars:       map[string]string{},
				},
			},
			Groups: []*api.Group{
				{
					Name:     "master",
					Children: []string{},
					Vars:     map[string]string{},
					Hosts:    []string{h.Name},
				},
			},
		},
	})
	resultId, err := ansible.RunAdhoc("master", "setup", "")
	if err != nil {
		return err
	}
	var result kobe.Result
	err = wait.Poll(5*time.Second, 5*time.Minute, func() (done bool, err error) {
		res, err := ansible.GetResult(resultId)
		if err != nil {
			return true, err
		}
		if res.Finished {
			if res.Success {
				result, err = kobe.ParseResult(res.Content)
				if err != nil {
					return true, err
				}
			} else {
				if res.Content != "" {
					result, err = kobe.ParseResult(res.Content)
					if err != nil {
						return true, err
					}
					result.GatherFailedInfo()
					if result.HostFailedInfo != nil && len(result.HostFailedInfo) > 0 {
						by, _ := json.Marshal(&result.HostFailedInfo)
						return true, errors.New(string(by))
					}
				}
			}
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	var facts interface{}
	if len(result.Plays) > 0 && len(result.Plays[0].Tasks) > 0 {
		facts = result.Plays[0].Tasks[0].Hosts[h.Name]["ansible_facts"]
	} else {
		return errors.New("no result return")
	}

	if facts == nil {
		return err
	} else {
		result, ok := facts.(map[string]interface{})
		if !ok {
			return err
		}
		h.Os = result["ansible_distribution"].(string)
		h.OsVersion = result["ansible_distribution_version"].(string)
		h.Architecture = result["ansible_architecture"].(string)
		if result["ansible_processor_vcpus"] != nil {
			h.CpuCore = int(result["ansible_processor_vcpus"].(float64))
		}
	}
	return nil
}
