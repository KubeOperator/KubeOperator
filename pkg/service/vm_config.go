package service

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/jinzhu/gorm"
)

var (
	ConfigNameExist = "NAME_EXISTS"
	ConfigExist     = "VM_CONFIG_EXISTS"
)

type VmConfigService interface {
	Page(num, size int) (page.Page, error)
	List() ([]dto.VmConfig, error)
	Batch(op dto.VmConfigOp) error
	Create(creation dto.VmConfigCreate) (*dto.VmConfig, error)
	Update(creation dto.VmConfigUpdate) (*dto.VmConfig, error)
}

type vmConfigService struct {
	vmConfigRepo repository.VmConfigRepository
}

func NewVmConfigService() VmConfigService {
	return &vmConfigService{
		vmConfigRepo: repository.NewVmConfigRepository(),
	}
}

func (v vmConfigService) Page(num, size int) (page.Page, error) {
	var page page.Page
	var vmConfigDTOs []dto.VmConfig
	total, mos, err := v.vmConfigRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		vmConfigDTO := new(dto.VmConfig)
		vmConfigDTO.VmConfig = mo
		vmConfigDTOs = append(vmConfigDTOs, *vmConfigDTO)
	}
	page.Total = total
	page.Items = vmConfigDTOs
	return page, err
}

func (v vmConfigService) List() ([]dto.VmConfig, error) {
	var configDTOS []dto.VmConfig
	configs, err := v.vmConfigRepo.List()
	if err != nil {
		return nil, err
	}
	for _, config := range configs {
		configDTO := new(dto.VmConfig)
		configDTO.VmConfig = config
		configDTOS = append(configDTOS, *configDTO)
	}
	return configDTOS, err
}

func (v vmConfigService) Batch(op dto.VmConfigOp) error {
	var deleteItems []model.VmConfig
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.VmConfig{
			Name: item.Name,
		})
	}
	return v.vmConfigRepo.Batch(op.Operation, deleteItems)
}

func (v vmConfigService) Create(creation dto.VmConfigCreate) (*dto.VmConfig, error) {

	old, err := v.vmConfigRepo.Get(creation.Name)
	if !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if old.ID != "" {
		return nil, errors.New(ConfigNameExist)
	}
	var config model.VmConfig
	db.DB.Where("cpu = ? AND memory = ?", creation.Cpu, creation.Memory).Find(&config)
	if config.ID != "" {
		return nil, errors.New(ConfigExist)
	}
	vmConfig := model.VmConfig{
		Name:     creation.Name,
		Cpu:      creation.Cpu,
		Memory:   creation.Memory,
		Disk:     50,
		Provider: creation.Provider,
	}
	err = v.vmConfigRepo.Save(&vmConfig)
	if err != nil {
		return nil, err
	}
	return &dto.VmConfig{VmConfig: vmConfig}, err
}

func (v vmConfigService) Update(creation dto.VmConfigUpdate) (*dto.VmConfig, error) {
	old, err := v.vmConfigRepo.Get(creation.Name)
	if err != nil {
		return nil, err
	}

	vmConfig := model.VmConfig{
		ID:       old.ID,
		Name:     creation.Name,
		Cpu:      creation.Cpu,
		Memory:   creation.Memory,
		Disk:     50,
		Provider: creation.Provider,
	}
	err = v.vmConfigRepo.Save(&vmConfig)
	if err != nil {
		return nil, err
	}
	return &dto.VmConfig{VmConfig: vmConfig}, err
}
