package service

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/KubeOperator/KubeOperator/pkg/util/ssh"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
	"k8s.io/apimachinery/pkg/util/wait"
)

type HostService interface {
	Get(name string) (dto.Host, error)
	List(projectName string) ([]dto.Host, error)
	Page(num, size int) (page.Page, error)
	Create(creation dto.HostCreate) (dto.Host, error)
	Delete(name string) error
	SyncList(names []dto.HostSync) error
	Sync(name string) (dto.Host, error)
	Batch(op dto.HostOp) error
	DownloadTemplateFile() error
}

type hostService struct {
	hostRepo       repository.HostRepository
	credentialRepo repository.CredentialRepository
}

func NewHostService() HostService {
	return &hostService{
		hostRepo:       repository.NewHostRepository(),
		credentialRepo: repository.NewCredentialRepository(),
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
	var (
		page      page.Page
		hostDTOs  []dto.Host
		hostIDs   []string
		resources []model.ProjectResource
		projects  []model.Project
	)
	total, mos, err := h.hostRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		hostIDs = append(hostIDs, mo.ID)
	}
	if err := db.DB.Where("resource_id in (?) AND resource_type = ?", hostIDs, constant.ResourceHost).Find(&resources).Error; err != nil {
		return page, err
	}
	if err := db.DB.Find(&projects).Error; err != nil {
		return page, err
	}

	for _, mo := range mos {
		isExist := false
		for _, res := range resources {
			if mo.ID == res.ResourceID {
				isExist = true
				for _, pro := range projects {
					if pro.ID == res.ProjectID {
						hostDTOs = append(hostDTOs, dto.Host{
							Host:        mo,
							ProjectName: pro.Name,
							ClusterName: mo.Cluster.Name,
							ZoneName:    mo.Zone.Name,
						})
						break
					}
				}
				break
			}
		}
		if !isExist {
			hostDTOs = append(hostDTOs, dto.Host{
				Host:        mo,
				ClusterName: mo.Cluster.Name,
				ZoneName:    mo.Zone.Name,
			})
		}
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
	credential, err := h.credentialRepo.GetById(creation.CredentialID)
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
	go h.RunGetHostConfig(&host)
	return dto.Host{Host: host}, err
}

func (h hostService) SyncList(hosts []dto.HostSync) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 2)
	for _, host := range hosts {
		if host.HostStatus == constant.ClusterCreating || host.HostStatus == constant.ClusterInitializing || host.HostStatus == constant.ClusterSynchronizing {
			continue
		}
		// 先更新所有待同步主机状态
		if err := db.DB.Model(&model.Host{}).Where("name = ?", host.HostName).Update("status", constant.ClusterSynchronizing).Error; err != nil {
			log.Errorf("update host status to synchronizing error: %s", err.Error())
		}

		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			log.Infof("gather host [%s] info", name)
			_, err := h.Sync(name)
			if err != nil {
				log.Errorf("gather host info error: %s", err.Error())
			}
		}(host.HostName)
	}
	return nil
}

func (h hostService) Sync(name string) (dto.Host, error) {
	host, err := h.hostRepo.Get(name)
	if err != nil {
		log.Errorf("host of %s not found error: %s", name, err.Error())
		return dto.Host{Host: host}, err
	}
	if err := h.GetHostConfig(&host); err != nil {
		host.Status = constant.ClusterFailed
		host.Message = err.Error()
		if err := syncHostInfoWithDB(&host); err != nil {
			log.Errorf("update host info error: %s", err.Error())
			return dto.Host{Host: host}, err
		}
		return dto.Host{Host: host}, err
	}
	host.Status = constant.ClusterRunning
	if err := syncHostInfoWithDB(&host); err != nil {
		log.Errorf("update host info error: %s", err.Error())
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
	result, _, _, err := client.Exec("sudo lspci|grep -i NVIDIA")
	if err != nil {
		host.HasGpu = false
		host.GpuNum = 0
	}
	host.GpuNum = strings.Count(result, "NVIDIA")
	if host.GpuNum == 0 {
		host.HasGpu = false
		host.GpuInfo = ""
	}
	if host.GpuNum > 0 {
		host.HasGpu = true
		s := strings.Index(result, "[")
		t := strings.Index(result, "]")
		host.GpuInfo = result[s+1 : t]
	}
	return err
}

func (h hostService) GetHostMem(host *model.Host) error {
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
	result, _, _, err := client.Exec("sudo dmidecode -t 17 | grep \"Size.*MB\" | awk '{s+=$2} END {print s}'")
	if err != nil {
		return err
	}
	host.Memory, _ = strconv.Atoi(strings.Trim(result, "\n"))
	return err
}

func (h hostService) RunGetHostConfig(host *model.Host) {
	err := h.GetHostConfig(host)
	if err != nil {
		host.Status = constant.ClusterFailed
		host.Message = err.Error()
		_ = h.hostRepo.Save(host)
		return
	}
	host.Status = constant.ClusterRunning
	_ = h.hostRepo.Save(host)
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
		host.Architecture = result["ansible_architecture"].(string)
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
	err = h.GetHostMem(host)
	if err != nil {
		return err
	}
	err = h.GetHostGpu(host)
	if err != nil {
		host.GpuNum = 0
		host.GpuInfo = ""
		host.HasGpu = false
		return nil
	}
	return nil
}

func (h hostService) DownloadTemplateFile() error {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "name")
	f.SetCellValue("Sheet1", "B1", "ip")
	f.SetCellValue("Sheet1", "C1", "port")
	f.SetCellValue("Sheet1", "D1", "credential (系统设置-凭据中的名称)")
	file, err := os.Create("./demo.xlsx")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = f.WriteTo(file)
	if err != nil {
		return err
	}
	return nil
}

func syncHostInfoWithDB(host *model.Host) error {
	tx := db.DB.Begin()
	if host.Name == "" {
		return nil
	}

	if len(host.Volumes) > 0 {
		for i := range host.Volumes {
			var volume model.Volume
			if notFound := tx.Where("host_id = ? AND name = ?", host.ID, host.Volumes[i].Name).
				First(&volume).RecordNotFound(); notFound {
				if err := tx.Create(&host.Volumes[i]).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				host.Volumes[i].ID = volume.ID
				if err := tx.Save(&host.Volumes[i]).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	if err := tx.Model(&model.Ip{}).Where("address = ?", host.Ip).Update("status", constant.IpUsed).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&model.Host{}).Where("id = ?", host.ID).
		Updates(map[string]interface{}{
			"memory":       host.Memory,
			"cpu_core":     host.CpuCore,
			"os":           host.Os,
			"os_version":   host.OsVersion,
			"gpu_num":      host.GpuNum,
			"gpu_info":     host.GpuInfo,
			"has_gpu":      host.HasGpu,
			"status":       host.Status,
			"message":      host.Message,
			"datastore":    host.Datastore,
			"architecture": host.Architecture,
		}).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}
