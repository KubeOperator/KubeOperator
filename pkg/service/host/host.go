package host

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/KubeOperator/kobe/api"
	"log"
	"os"
	"time"
)

var (
	getHostConfigError = "get host config error"
)

func Page(num, size int) (host []hostModel.Host, total int, err error) {
	err = db.DB.Model(hostModel.Host{}).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Find(&host).
		Error
	return
}

func List() (host []hostModel.Host, err error) {
	err = db.DB.Model(hostModel.Host{}).Find(&host).Error
	return
}

func Get(name string) (*hostModel.Host, error) {
	var result hostModel.Host
	err := db.DB.Model(hostModel.Host{}).Where(&result).First(&result).Error
	return &result, err
}

func Save(item *hostModel.Host) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}

func Delete(name string) error {
	var h hostModel.Host
	h.Name = name
	return db.DB.Delete(&h).Error
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

	password, privateKey, err := getHostPasswordAndPrivateKey(host)
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
		return err
	}
	if err := client.Ping(); err != nil {
		return err
	}
	return err
}

func GetHostConfig(host *hostModel.Host) error {

	//TODO
	password, _, err := getHostPasswordAndPrivateKey(host)
	if err != nil {
		log.Fatal(getHostConfigError, err)
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
		log.Fatal(getHostConfigError, err)
		return err
	}
	err = ansible.Watch(os.Stdout, resultId)
	if err != nil {
		log.Fatal(getHostConfigError, err)
		return err
	}
	res, err := ansible.GetResult(resultId)
	if err != nil {
		log.Fatal(getHostConfigError, err)
		return err
	}
	result, err := kobe.ParseResult(res.Content)
	if err != nil {
		log.Fatal(getHostConfigError, err)
		return err
	}
	facts := result.Plays[0].Tasks[0].Hosts[host.Name]["ansible_facts"]
	if facts == nil {
		log.Fatal(getHostConfigError, err)
		return err
	} else {
		result, ok := facts.(map[string]interface{})
		if !ok {
			log.Fatal(getHostConfigError, err)
			return err
		}
		host.Os = result["ansible_distribution"].(string)
		host.OsVersion = result["ansible_distribution_version"].(string)
		host.Memory = int(result["ansible_memtotal_mb"].(float64))
		host.CpuCore = int(result["ansible_processor_vcpus"].(float64))
		err = Save(host)
		if err != nil {
			log.Fatal(getHostConfigError, err)
			return err
		}
	}
	return nil
}

func getHostPasswordAndPrivateKey(host *hostModel.Host) (string, []byte, error) {
	var err error = nil
	password := ""
	privateKey := []byte("")
	if "password" == host.Credential.Type {
		pwd, err := encrypt.StringDecrypt(host.Credential.Password)
		password = pwd
		if err != nil {
			log.Fatal(getHostConfigError, err)
			return password, privateKey, err
		}
	}
	if "privateKey" == host.Credential.Type {
		privateKey = []byte(host.Credential.PrivateKey)
	}
	return password, privateKey, err
}
