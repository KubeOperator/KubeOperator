package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/errorf"
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
	Get(name string) (*dto.Host, error)
	List(projectName string, conditions condition.Conditions) ([]dto.Host, error)
	Page(num, size int, projectName string, conditions condition.Conditions) (*page.Page, error)
	Create(creation dto.HostCreate) (*dto.Host, error)
	Update(host dto.HostUptate) (*dto.Host, error)
	Delete(name string) error
	SyncList(names []dto.HostSync) error
	Sync(name string) (dto.Host, error)
	Batch(op dto.HostOp) error
	DownloadTemplateFile() error
	RunGetHostConfig(host *model.Host)
	ImportHosts(file []byte) error
}

type hostService struct {
	hostRepo          repository.HostRepository
	credentialRepo    repository.CredentialRepository
	credentialService CredentialService
	projectRepository repository.ProjectRepository
}

func NewHostService() HostService {
	return &hostService{
		hostRepo:          repository.NewHostRepository(),
		credentialRepo:    repository.NewCredentialRepository(),
		credentialService: NewCredentialService(),
		projectRepository: repository.NewProjectRepository(),
	}
}

func (h *hostService) Get(name string) (*dto.Host, error) {
	var (
		mo      model.Host
		hostDTO dto.Host
	)

	if err := db.DB.Where("name = ?", name).
		Preload("Volumes").
		Preload("Credential").
		Preload("Zone").
		Preload("Cluster").
		First(&mo).Error; err != nil {
		return nil, err
	}

	hostDTO = dto.Host{
		Host:        mo,
		ClusterName: mo.Cluster.Name,
		ZoneName:    mo.Zone.Name,
	}
	return &hostDTO, nil
}

func (h *hostService) List(projectName string, conditions condition.Conditions) ([]dto.Host, error) {
	var (
		mos       []model.Host
		hostDTOs  []dto.Host
		projects  []model.Project
		resources []model.ProjectResource
	)

	d := db.DB.Model(model.Host{})
	if err := dbUtil.WithConditions(&d, model.Host{}, conditions); err != nil {
		return hostDTOs, nil
	}

	if len(projectName) != 0 {
		res, err := dbUtil.WithProjectResource(&d, projectName, constant.ResourceHost)
		if err != nil {
			return hostDTOs, nil
		}
		resources = res[:]
	} else {
		if err := db.DB.Where(model.ProjectResource{ResourceType: constant.ResourceHost}).Find(&resources).Error; err != nil {
			return hostDTOs, nil
		}
	}

	if err := d.
		Preload("Volumes").
		Preload("Cluster").
		Preload("Zone").
		Find(&mos).Error; err != nil {
		return hostDTOs, nil
	}
	if err := db.DB.Find(&projects).Error; err != nil {
		return hostDTOs, nil
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
	return hostDTOs, nil
}

func (h *hostService) Page(num, size int, projectName string, conditions condition.Conditions) (*page.Page, error) {
	var (
		p                page.Page
		hostDTOs         []dto.Host
		mos              []model.Host
		projects         []model.Project
		clusters         []model.Cluster
		projectResources []model.ProjectResource
		clusterResources []model.ClusterResource
	)
	d := db.DB.Model(model.Host{})
	if err := dbUtil.WithConditions(&d, model.Host{}, conditions); err != nil {
		return &p, err
	}
	if len(projectName) != 0 {
		res, err := dbUtil.WithProjectResource(&d, projectName, constant.ResourceHost)
		if err != nil {
			return nil, err
		}
		projectResources = res[:]
	} else {
		if err := db.DB.Where(model.ProjectResource{ResourceType: constant.ResourceHost}).Find(&projectResources).Error; err != nil {
			return &p, err
		}
	}
	if err := db.DB.Find(&clusterResources).Error; err != nil {
		return &p, err
	}

	if err := d.
		Count(&p.Total).
		Order("name").
		Offset((num - 1) * size).
		Limit(size).
		Preload("Volumes").
		Preload("Zone").
		Preload("Credential").
		Find(&mos).Error; err != nil {
		return &p, err
	}

	if err := db.DB.Find(&projects).Error; err != nil {
		return &p, err
	}
	if err := db.DB.Find(&clusters).Error; err != nil {
		return &p, err
	}

	for _, mo := range mos {
		hostItem := dto.Host{Host: mo, ZoneName: mo.Zone.Name, CredentialName: mo.Credential.Name}
		for _, res := range projectResources {
			if mo.ID == res.ResourceID {
				for _, pro := range projects {
					if pro.ID == res.ProjectID {
						hostItem.ProjectName = pro.Name
						break
					}
				}
				break
			}
		}
		for _, res := range clusterResources {
			if mo.ID == res.ResourceID {
				for _, clu := range clusters {
					if clu.ID == res.ClusterID {
						hostItem.ClusterName = clu.Name
						break
					}
				}
			}
		}
		hostDTOs = append(hostDTOs, hostItem)
	}
	p.Items = hostDTOs
	return &p, nil
}

func (h *hostService) Delete(name string) error {
	var host model.Host
	if err := db.DB.Where(model.Host{Name: name}).First(&host).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&host).Error; err != nil {
		return err
	}
	return nil
}

func (h *hostService) Create(creation dto.HostCreate) (*dto.Host, error) {
	tx := db.DB.Begin()
	var num int
	if err := db.DB.Model(model.SystemRegistry{}).Where("hostname = ?", creation.Ip).Count(&num).Error; err != nil {
		return nil, err
	}
	if num != 0 {
		return nil, errors.New("IS_LOCAL_HOST")
	}
	var credential model.Credential
	if creation.CredentialID != "" {
		if err := db.DB.Where(model.Credential{ID: creation.CredentialID}).First(&credential).Error; err != nil {
			return nil, err
		}
	} else {
		var password string
		if creation.Credential.Password != "" {
			p, err := encrypt.StringEncrypt(creation.Credential.Password)
			if err != nil {
				return nil, err
			}
			password = p
		}
		c := model.Credential{
			Name:       creation.Credential.Name,
			Password:   password,
			Username:   creation.Credential.Username,
			PrivateKey: creation.Credential.PrivateKey,
			Type:       creation.Credential.Type,
		}
		if err := tx.Create(&c).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		credential = c
	}

	host := model.Host{
		BaseModel:    common.BaseModel{},
		Name:         creation.Name,
		Ip:           creation.Ip,
		Port:         creation.Port,
		CredentialID: credential.ID,
		Credential:   credential,
		Status:       constant.ClusterInitializing,
	}
	if err := tx.Create(&host).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if creation.Cluster != "" {
		var cluster model.Cluster
		if err := tx.Where("name = ?", creation.Cluster).Find(&cluster).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Create(&model.ClusterResource{
			ResourceType: constant.ResourceHost,
			ResourceID:   host.ID,
			ClusterID:    cluster.ID,
		}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	var project model.Project
	if err := tx.Where("name = ?", creation.Project).Find(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Create(&model.ProjectResource{
		ResourceType: constant.ResourceHost,
		ResourceID:   host.ID,
		ProjectID:    project.ID,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	go h.RunGetHostConfig(&host)
	return &dto.Host{Host: host}, nil
}

func (h *hostService) Update(host dto.HostUptate) (*dto.Host, error) {
	tx := db.DB.Begin()
	var num int
	if err := db.DB.Model(model.SystemRegistry{}).Where("hostname = ?", host.Ip).Count(&num).Error; err != nil {
		return nil, err
	}
	if num != 0 {
		return nil, errors.New("IS_LOCAL_HOST")
	}
	var credential model.Credential
	if host.CredentialID != "" {
		if err := db.DB.Where(model.Credential{ID: host.CredentialID}).First(&credential).Error; err != nil {
			return nil, err
		}
	} else {
		var password string
		if host.Credential.Password != "" {
			p, err := encrypt.StringEncrypt(host.Credential.Password)
			if err != nil {
				return nil, err
			}
			password = p
		}
		c := model.Credential{
			Name:       host.Credential.Name,
			Password:   password,
			Username:   host.Credential.Username,
			PrivateKey: host.Credential.PrivateKey,
			Type:       host.Credential.Type,
		}
		if err := tx.Create(&c).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		credential = c
	}
	var oldHost model.Host
	if err := tx.Where("name = ?", host.Name).First(&oldHost).Error; err != nil {
		return nil, err
	}
	if err := tx.Model(&model.Host{}).Where("id = ?", oldHost.ID).Updates(map[string]interface{}{
		"Ip":           host.Ip,
		"Port":         host.Port,
		"CredentialID": credential.ID,
		"Status":       constant.ClusterInitializing,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	newHost := model.Host{
		ID:           oldHost.ID,
		Name:         host.Name,
		Ip:           host.Ip,
		Port:         host.Port,
		CredentialID: credential.ID,
		Credential:   credential,
		Status:       constant.ClusterInitializing,
	}

	tx.Commit()
	go h.RunGetHostConfig(&newHost)
	return &dto.Host{Host: newHost}, nil
}

func (h *hostService) SyncList(hosts []dto.HostSync) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 2)
	for _, host := range hosts {
		if host.HostStatus == constant.ClusterCreating || host.HostStatus == constant.ClusterInitializing || host.HostStatus == constant.ClusterSynchronizing {
			continue
		}
		// 先更新所有待同步主机状态
		if err := db.DB.Model(&model.Host{}).Where("name = ?", host.HostName).Update("status", constant.ClusterSynchronizing).Error; err != nil {
			logger.Log.Errorf("update host status to synchronizing error: %s", err.Error())
		}

		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			logger.Log.Infof("gather host [%s] info", name)
			_, err := h.Sync(name)
			if err != nil {
				logger.Log.Errorf("gather host info error: %s", err.Error())
			}
		}(host.HostName)
	}
	return nil
}

func (h *hostService) Sync(name string) (dto.Host, error) {
	host, err := h.hostRepo.Get(name)
	if err != nil {
		logger.Log.Errorf("host of %s not found error: %s", name, err.Error())
		return dto.Host{Host: host}, err
	}
	if err := h.GetHostConfig(&host); err != nil {
		host.Status = constant.ClusterFailed
		host.Message = err.Error()
		if err := syncHostInfoWithDB(&host); err != nil {
			logger.Log.Errorf("update host info error: %s", err.Error())
			return dto.Host{Host: host}, err
		}
		return dto.Host{Host: host}, err
	}
	host.Status = constant.ClusterRunning
	if err := syncHostInfoWithDB(&host); err != nil {
		logger.Log.Errorf("update host info error: %s", err.Error())
		return dto.Host{Host: host}, err
	}
	return dto.Host{Host: host}, nil
}

func (h *hostService) Batch(op dto.HostOp) error {
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

func (h *hostService) GetHostGpu(host *model.Host) error {
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

func (h *hostService) GetHostMem(host *model.Host) error {
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
	if len(result) == 0 {
		result, _, _, err = client.Exec("sudo dmidecode -t 17 | grep \"Size.*GB\" | awk '{s+=$2} END {print s}'")
		if err != nil {
			return err
		}
		me, _ := strconv.Atoi(strings.Trim(result, "\n"))
		host.Memory = me * 1024
	} else {
		me, err := strconv.Atoi(strings.Trim(result, "\n"))
		if err != nil {
			return err
		}
		host.Memory = me
	}
	return err
}

func (h *hostService) RunGetHostConfig(host *model.Host) {
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

func (h *hostService) GetHostConfig(host *model.Host) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Log.Error("gather fact error!")
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
			if device["model"] == "Virtual disk" {
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

func (h *hostService) DownloadTemplateFile() error {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "name (中文、大小写英文、数字和-)")
	f.SetCellValue("Sheet1", "B1", "ip")
	f.SetCellValue("Sheet1", "C1", "port")
	f.SetCellValue("Sheet1", "D1", "credential (系统设置-凭据中的名称)")
	f.SetCellValue("Sheet1", "E1", "project (项目-项目名称)")
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

func (h *hostService) ImportHosts(file []byte) error {
	f, err := os.Create("./import.xlsx")
	if err != nil {
		return err
	}
	defer f.Close()
	err = ioutil.WriteFile("./import.xlsx", file, 0775)
	if err != nil {
		return err
	}
	xlsx, err := excelize.OpenFile("./import.xlsx")
	if err != nil {
		return err
	}
	rows := xlsx.GetRows("Sheet1")
	if len(rows) == 0 {
		return errors.New("HOST_IMPORT_ERROR_NULL")
	}
	var hosts []model.Host
	//var errMsg string
	var failedNum int
	var errs errorf.CErrFs
	for index, row := range rows {
		if index == 0 {
			continue
		}
		if len(row) < 5 {
			errs = errs.Add(errorf.New("HOST_IMPORT_NOT_COMPLETE_VALUE", strconv.Itoa(index)))
			failedNum++
			continue
		}
		if row[0] == "" || row[1] == "" || row[2] == "" || row[3] == "" || row[4] == "" {
			errs = errs.Add(errorf.New("HOST_IMPORT_NOT_COMPLETE_VALUE", strconv.Itoa(index)))
			failedNum++
			continue
		}
		port, err := strconv.Atoi(row[2])
		if err != nil {
			errs = errs.Add(errorf.New("HOST_IMPORT_WRONG_FORMAT", strconv.Itoa(index)))
			failedNum++
			continue
		}
		credential, err := h.credentialRepo.Get(row[3])
		if err != nil {
			errs = errs.Add(errorf.New("HOST_IMPORT_CREDENTIAL_NOT_FOUND", strconv.Itoa(index)))
			failedNum++
			continue
		}
		project, err := h.projectRepository.Get(row[4])
		if err != nil {
			errs = errs.Add(errorf.New("HOST_IMPORT_PROJECT_NOT_FOUND", strconv.Itoa(index)))
			failedNum++
			continue
		}
		host := model.Host{
			Name:         strings.Trim(row[0], " "),
			Ip:           strings.Trim(row[1], " "),
			Port:         port,
			CredentialID: credential.ID,
			Status:       constant.ClusterInitializing,
			Credential:   credential,
			ClusterID:    project.ID,
		}
		hosts = append(hosts, host)
	}

	if len(errs) > 0 {
		errs = errs.Add(errorf.New("HOST_IMPORT_FAILED_NUM", strconv.Itoa(failedNum)))
	}

	itemProjectName := ""
	for _, host := range hosts {
		itemProjectName = host.ClusterID
		host.ClusterID = ""
		if err := h.hostRepo.Save(&host); err != nil {
			errs = errs.Add(errorf.New("HOST_IMPORT_FAILED_SAVE", host.Name, err.Error()))
			continue
		}
		if err := db.DB.Create(&model.ProjectResource{
			ResourceType: constant.ResourceHost,
			ResourceID:   host.ID,
			ProjectID:    itemProjectName,
		}).Error; err != nil {
			errs = errs.Add(errorf.New("HOST_RESOURCE_FAILED_BIND", host.Name, err.Error()))
			continue
		}
		saveHost := host
		go h.RunGetHostConfig(&saveHost)
		var ip model.Ip
		if err := db.DB.Where("address = ?", host.Ip).First(&ip).Error; err != nil {
			continue
		}
		if ip.ID != "" {
			ip.Status = constant.IpUsed
			if err := db.DB.Save(&ip).Error; err != nil {
				continue
			}
		}
	}
	if len(errs) > 0 {
		return errs
	} else {
		return nil
	}
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
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
