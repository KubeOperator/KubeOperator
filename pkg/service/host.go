package service

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"strings"
	"time"
)

type HostService interface {
	Get(name string) (dto.Host, error)
	List(projectName string) ([]dto.Host, error)
	Page(num, size int) (page.Page, error)
	Create(creation dto.HostCreate) (dto.Host, error)
	Delete(name string) error
	Sync(name string) (dto.Host, error)
	Batch(op dto.HostOp) error
}

type hostService struct {
	hostRepo repository.HostRepository
}

func NewHostService() HostService {
	return &hostService{
		hostRepo: repository.NewHostRepository(),
	}
}

func (h hostService) Get(name string) (dto.Host, error) {
	var hostDTO dto.Host
	mo, err := h.hostRepo.Get(name)
	if err != nil {
		return hostDTO, err
	}
	hostDTO = dto.Host{
		Host:        mo,
		ClusterName: mo.Cluster.Name,
		ZoneName:    mo.Zone.Name,
	}
	return hostDTO, err
}

func (h hostService) List(projectName string) ([]dto.Host, error) {
	var hostDTOs []dto.Host
	mos, err := h.hostRepo.List(projectName)
	if err != nil {
		return hostDTOs, err
	}
	for _, mo := range mos {
		hostDTOs = append(hostDTOs, dto.Host{
			Host:        mo,
			ClusterName: mo.Cluster.Name,
			ZoneName:    mo.Zone.Name,
		})
	}
	return hostDTOs, err
}

func (h hostService) Page(num, size int) (page.Page, error) {
	var page page.Page
	var hostDTOs []dto.Host
	total, mos, err := h.hostRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		hostDTOs = append(hostDTOs, dto.Host{
			Host:        mo,
			ClusterName: mo.Cluster.Name,
			ZoneName:    mo.Zone.Name,
		})
	}
	page.Total = total
	page.Items = hostDTOs
	return page, err
}

func (h hostService) Delete(name string) error {
	err := h.hostRepo.Delete(name)
	if err != nil {
		return err
	}
	return nil
}

func (h hostService) Create(creation dto.HostCreate) (dto.Host, error) {

	credential, err := repository.NewCredentialRepository().GetById(creation.CredentialID)
	if err != nil {
		return dto.Host{}, err
	}
	host := model.Host{
		BaseModel:    common.BaseModel{},
		Name:         creation.Name,
		Ip:           creation.Ip,
		Port:         creation.Port,
		CredentialID: creation.CredentialID,
		Credential:   credential,
		Status:       constant.ClusterInitializing,
	}

	err = h.hostRepo.Save(&host)
	if err != nil {
		return dto.Host{}, err
	}
	go h.RunGetHostConfig(host)
	go h.GetHostGpu(&host)
	return dto.Host{Host: host}, err
}

func (h hostService) Sync(name string) (dto.Host, error) {
	host, err := h.hostRepo.Get(name)
	if err != nil {
		return dto.Host{Host: host}, err
	}
	err = h.GetHostConfig(&host)
	if err != nil {
		host.Status = constant.ClusterFailed
		host.Message = err.Error()
		_ = h.hostRepo.Save(&host)
		return dto.Host{Host: host}, err
	}
	err = h.GetHostGpu(&host)
	if err != nil {
		return dto.Host{Host: host}, err
	}
	host.Status = constant.ClusterRunning
	err = h.hostRepo.Save(&host)
	if err != nil {
		return dto.Host{Host: host}, err
	}
	return dto.Host{Host: host}, nil
}

func (h hostService) Batch(op dto.HostOp) error {
	var deleteItems []model.Host
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.Host{
			BaseModel: common.BaseModel{},
			ID:        item.ID,
			Name:      item.Name,
		})
	}
	return h.hostRepo.Batch(op.Operation, deleteItems)
}

func (h hostService) GetHostGpu(host *model.Host) error {
	password, privateKey, err := host.GetHostPasswordAndPrivateKey()
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
	result, _, _, err := client.Exec("lspci|grep -i NVIDIA")
	if err != nil {
		host.HasGpu = false
		host.GpuNum = 0
	}
	host.GpuNum = strings.Count(result, "NVIDIA")
	if host.GpuNum > 0 {
		host.HasGpu = true
		s := strings.Index(result, "[")
		t := strings.Index(result, "]")
		host.Gpus = result[s+1 : t]
	}
	_ = h.hostRepo.Save(host)
	return err
}

func (h hostService) RunGetHostConfig(host model.Host) {
	host.Status = constant.ClusterInitializing
	_ = h.hostRepo.Save(&host)
	err := h.GetHostConfig(&host)
	if err != nil {
		host.Status = constant.ClusterFailed
		host.Message = err.Error()
		_ = h.hostRepo.Save(&host)
		return
	}
	host.Status = constant.ClusterRunning
	_ = h.hostRepo.Save(&host)
}

func (h hostService) GetHostConfig(host *model.Host) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error("gather fact error!")
		}
	}()

	password, privateKey, err := host.GetHostPasswordAndPrivateKey()
	if err != nil {
		return err
	}
	ansible := kobe.NewAnsible(&kobe.Config{
		Inventory: &api.Inventory{
			Hosts: []*api.Host{
				{
					Ip:         host.Ip,
					Name:       host.Name,
					Port:       int32(host.Port),
					User:       host.Credential.Username,
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
					Hosts:    []string{host.Name},
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
		if result["ansible_memtotal_mb"] != nil {
			host.Memory = int(result["ansible_memtotal_mb"].(float64))
		}
		if result["ansible_processor_vcpus"] != nil {
			host.CpuCore = int(result["ansible_processor_vcpus"].(float64))
		}
		devices := result["ansible_devices"].(map[string]interface{})
		var volumes []model.Volume
		for i := range devices {
			device := devices[i].(map[string]interface{})
			if "Virtual disk" == device["model"] {
				v := model.Volume{
					ID:     uuid.NewV4().String(),
					Name:   "/dev/" + i,
					Size:   device["size"].(string),
					HostID: host.ID,
				}
				volumes = append(volumes, v)
			}
		}
		host.Volumes = volumes
	}
	return nil
}
