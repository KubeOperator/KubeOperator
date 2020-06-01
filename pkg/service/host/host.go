package host

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"time"
)

var (
	getHostConfigError = "get host config error,%s"
)

func Page(num, size int) (host []hostModel.Host, total int, err error) {
	err = db.DB.Model(hostModel.Host{}).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Preload("Volumes").
		Find(&host).
		Error
	return
}

func List() (host []hostModel.Host, err error) {
	err = db.DB.Model(hostModel.Host{}).Preload("Volumes").Find(&host).Error
	return
}

func Get(name string) (hostModel.Host, error) {
	var result hostModel.Host
	result.Name = name
	if err := db.DB.Where(result).First(&result).Error; err != nil {
		return result, err
	}
	if err := db.DB.First(&result).Related(&result.Volumes).Error; err != nil {
		return result, err
	}
	return result, nil
}

func Save(item *hostModel.Host) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}

func Delete(name string) error {
	tx := db.DB.Begin()
	var h hostModel.Host
	h.Name = name
	if err := db.DB.Delete(&h).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := db.DB.Where(hostModel.Volume{HostID: h.ID}).Delete(hostModel.Volume{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func Batch(operation string, items []hostModel.Host) ([]hostModel.Host, error) {
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, item := range items {
			err := db.DB.Model(hostModel.Host{}).Delete(&item).Error
			if err != nil {
				tx.Rollback()
			}
		}
		tx.Commit()
	default:
		return nil, constant.NotSupportedBatchOperation
	}
	return items, nil
}

func GetHostGpu(host *hostModel.Host) error {

	password, privateKey, err := GetHostPasswordAndPrivateKey(host)
	if err != nil {
		return err
	}

	client, err := ssh.New(&ssh.Config{
		User:        host.Credential.Username,
		Host:        host.Ip,
		Port:        host.Port,
		Password:    password,
		PrivateKey:  privateKey,
		PassPhrase:  nil,
		DialTimeOut: 5 * time.Second,
		Retry:       3,
	})
	if err != nil {
		host.Status = hostModel.SshError
		return err
	}
	if err := client.Ping(); err != nil {
		host.Status = hostModel.Disconnect
		return err
	}
	return err
}

func RunGetHostConfig(host *hostModel.Host) {
	err := GetHostConfig(host)
	if err != nil {
		if sErr := Save(host); sErr != nil {
		}
		log.Fatalf("get host [%s] config failed reason: %s", host.Name, err.Error())
	}
	fmt.Println(host)
	if sErr := Save(host); sErr != nil {
	}
}

func GetHostConfig(host *hostModel.Host) error {

	//TODO
	password, _, err := GetHostPasswordAndPrivateKey(host)
	if err != nil {
		return err
	}

	ansible := kobe.NewAnsible(&kobe.Config{
		Host: "localhost",
		Port: 8088,
		Inventory: api.Inventory{
			Hosts: []*api.Host{
				{
					Ip:       host.Ip,
					Name:     host.Name,
					Port:     int32(host.Port),
					User:     host.Credential.Username,
					Password: password,
					Vars:     map[string]string{},
				},
			},
			Groups: []*api.Group{
				{
					Name:     "master",
					Children: []string{},
					Vars:     map[string]string{},
					Hosts:    []string{host.Name},
				},
			},
		},
	})
	resultId, err := ansible.RunAdhoc("master", "setup", "")
	if err != nil {
		host.Status = hostModel.AnsibleError
		return err
	}
	if err = ansible.Watch(os.Stdout, resultId); err != nil {
		host.Status = hostModel.AnsibleError
		return err
	}
	res, err := ansible.GetResult(resultId)
	if err != nil {
		host.Status = hostModel.AnsibleError
		return err
	}
	result, err := kobe.ParseResult(res.Content)
	if err != nil {
		host.Status = hostModel.AnsibleError
		return err
	}
	facts := result.Plays[0].Tasks[0].Hosts[host.Name]["ansible_facts"]
	if facts == nil {
		host.Status = hostModel.AnsibleError
		return err
	} else {
		result, ok := facts.(map[string]interface{})
		if !ok {
			return err
		}
		host.Os = result["ansible_distribution"].(string)
		host.OsVersion = result["ansible_distribution_version"].(string)
		host.Memory = int(result["ansible_memtotal_mb"].(float64))
		host.CpuCore = int(result["ansible_processor_vcpus"].(float64))

		devices := result["ansible_devices"].(map[string]interface{})

		var volumes []hostModel.Volume
		for index, _ := range devices {
			device := devices[index].(map[string]interface{})
			if "Virtual disk" == device["model"] {
				v := hostModel.Volume{
					ID:     uuid.NewV4().String(),
					Name:   "/dev/" + index,
					Size:   device["size"].(string),
					HostID: host.ID,
				}
				volumes = append(volumes, v)
			}
		}
		host.Volumes = volumes
		if err = Save(host); err != nil {
			return err
		}
	}
	return nil
}

func GetHostPasswordAndPrivateKey(host *hostModel.Host) (string, []byte, error) {
	var err error = nil
	password := ""
	privateKey := []byte("")
	if "password" == host.Credential.Type {
		pwd, err := encrypt.StringDecrypt(host.Credential.Password)
		password = pwd
		if err != nil {
			log.Fatalf(getHostConfigError, err.Error())
			return password, privateKey, err
		}
	}
	if "privateKey" == host.Credential.Type {
		privateKey = []byte(host.Credential.PrivateKey)
	}
	return password, privateKey, err
}

func ListHostByCredentialID(credentialID string) ([]hostModel.Host, error) {
	var host []hostModel.Host
	err := db.DB.Model(hostModel.Host{
		CredentialID: credentialID,
	}).Find(&host).Error
	return host, err
}
