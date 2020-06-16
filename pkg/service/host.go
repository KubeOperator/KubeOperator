package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/KubeOperator/kobe/api"
	"log"
	"os"
	"time"
)

type HostService interface {
	Get(name string) (dto.Host, error)
	List() ([]dto.Host, error)
	Page(num, size int) (int, []dto.Host, error)
	Create(creation dto.HostCreate) error
	Delete(name string) error
	Sync(name string) error
}

var (
	getHostConfigError = "get host config error,%s"
)

type hostService struct {
	hostRepo repository.HostRepository
}

func (h hostService) Get(name string) (dto.Host, error) {
	var hostDTO dto.Host
	mo, err := h.hostRepo.Get(name)
	if err != nil {
		return hostDTO, err
	}
	hostDTO.Host = mo
	return hostDTO, err
}

func (h hostService) List() ([]dto.Host, error) {
	var hostDTOs []dto.Host
	mos, err := h.hostRepo.List()
	if err != nil {
		return hostDTOs, err
	}
	for _, mo := range mos {
		hostDTOs = append(hostDTOs, dto.Host{Host: mo})
	}
	return hostDTOs, err
}

func (h hostService) Page(num, size int) (int, []dto.Host, error) {
	var total int
	var hostDTOs []dto.Host
	total, mos, err := h.hostRepo.Page(num, size)
	if err != nil {
		return total, hostDTOs, err
	}
	for _, mo := range mos {
		hostDTOs = append(hostDTOs, dto.Host{Host: mo})
	}
	return total, hostDTOs, err
}

func (h hostService) Delete(name string) error {
	err := h.hostRepo.Delete(name)
	if err != nil {
		return err
	}
	return nil
}

func (h hostService) Create(creation dto.HostCreate) error {

	credential, err := credentialService{}.credentialRepo.GetById(creation.CredentialID)
	if err != nil {
		return err
	}

	host := model.Host{
		BaseModel:    common.BaseModel{},
		Name:         creation.Name,
		Ip:           creation.Ip,
		Port:         creation.Port,
		CredentialID: creation.CredentialID,
		Credential:   credential,
	}
	err = h.GetHostGpu(&host)
	if err != nil {
		return err
	}
	go h.RunGetHostConfig(host)

	err = h.hostRepo.Save(&host)
	if err != nil {
		return err
	}
	return err
}

func (h hostService) Sync(name string) error {
	host, err := h.hostRepo.Get(name)
	if err != nil {
		return err
	}
	err = h.GetHostConfig(&host)
	if err != nil {
		return err
	}
	err = h.hostRepo.Save(&host)
	if err != nil {
		return err
	}
	return nil
}

func (h hostService) GetHostGpu(host *model.Host) error {

	password, privateKey, err := h.GetHostPasswordAndPrivateKey(*host)
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
		host.Status = model.SshError
		return err
	}
	if err := client.Ping(); err != nil {
		host.Status = model.Disconnect
		return err
	}
	return err
}
func (h hostService) RunGetHostConfig(host model.Host) {
	err := h.GetHostConfig(&host)
	if err != nil {
		if sErr := h.hostRepo.Save(&host); sErr != nil {
		}
		log.Fatalf("get host [%s] config failed reason: %s", host.Name, err.Error())
	}
	if sErr := h.hostRepo.Save(&host); sErr != nil {
	}
}

func (h hostService) GetHostConfig(host *model.Host) error {

	host.Status = model.AnsibleError
	//TODO
	password, _, err := h.GetHostPasswordAndPrivateKey(*host)
	if err != nil {
		return err
	}

	ansible := kobe.NewAnsible(&kobe.Config{
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
		return err
	}
	if err = ansible.Watch(os.Stdout, resultId); err != nil {
		return err
	}
	res, err := ansible.GetResult(resultId)
	if err != nil {
		return err
	}
	result, err := kobe.ParseResult(res.Content)
	if err != nil {
		return err
	}
	var facts interface{}
	if len(result.Plays) > 0 && len(result.Plays[0].Tasks) > 0 {
		facts = result.Plays[0].Tasks[0].Hosts[host.Name]["ansible_facts"]
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
		host.Os = result["ansible_distribution"].(string)
		host.OsVersion = result["ansible_distribution_version"].(string)
		host.Memory = int(result["ansible_memtotal_mb"].(float64))
		host.CpuCore = int(result["ansible_processor_vcpus"].(float64))

		devices := result["ansible_devices"].(map[string]interface{})

		var volumes []model.Volume
		for index, _ := range devices {
			device := devices[index].(map[string]interface{})
			if "Virtual disk" == device["model"] {
				v := model.Volume{
					Name:   "/dev/" + index,
					Size:   device["size"].(string),
					HostID: host.ID,
				}
				volumes = append(volumes, v)
			}
		}
		host.Volumes = volumes
		host.Status = model.Running
	}
	return nil
}

func (h hostService) GetHostPasswordAndPrivateKey(host model.Host) (string, []byte, error) {
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
