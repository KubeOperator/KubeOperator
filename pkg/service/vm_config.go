package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
	"github.com/jinzhu/gorm"
)

var (
	ConfigNameExist = "NAME_EXISTS"
	ConfigExist     = "VM_CONFIG_EXISTS"
)

type VmConfigService interface {
	Page(num, size int, conditions condition.Conditions) (*page.Page, error)
	List(conditions condition.Conditions) ([]dto.VmConfig, error)
	Batch(op dto.VmConfigOp) error
	Create(creation dto.VmConfigCreate) (*dto.VmConfig, error)
	Update(name string, creation dto.VmConfigUpdate) (*dto.VmConfig, error)
	Get(name string) (*dto.VmConfig, error)
	Delete(name string) error
}

type vmConfigService struct {
	vmConfigRepo repository.VmConfigRepository
}

func NewVmConfigService() VmConfigService {
	return &vmConfigService{
		vmConfigRepo: repository.NewVmConfigRepository(),
	}
}

func (v vmConfigService) Page(num, size int, conditions condition.Conditions) (*page.Page, error) {

	var (
		p            page.Page
		vmConfigDTOs []dto.VmConfig
		vmConfigs    []model.VmConfig
	)
	d := db.DB.Model(model.VmConfig{})
	if err := dbUtil.WithConditions(&d, model.VmConfig{}, conditions); err != nil {
		return nil, err
	}
	if err := d.Count(&p.Total).Order("cpu").Offset((num - 1) * size).Limit(size).Find(&vmConfigs).Error; err != nil {
		return nil, err
	}
	for _, mo := range vmConfigs {
		vmConfigDTO := new(dto.VmConfig)
		vmConfigDTO.VmConfig = mo
		vmConfigDTOs = append(vmConfigDTOs, *vmConfigDTO)
	}
	p.Items = vmConfigDTOs
	return &p, nil
}

func (v vmConfigService) Get(name string) (*dto.VmConfig, error) {
	var vmConfigDTO dto.VmConfig
	vmConfig, err := v.vmConfigRepo.Get(name)
	if err != nil {
		return nil, err
	}
	vmConfigDTO.VmConfig = vmConfig
	return &vmConfigDTO, nil
}

func (v vmConfigService) List(conditions condition.Conditions) ([]dto.VmConfig, error) {
	var configDTOS []dto.VmConfig
	var configs []model.VmConfig
	d := db.DB.Model(model.VmConfig{})
	if err := dbUtil.WithConditions(&d, model.VmConfig{}, conditions); err != nil {
		return nil, err
	}
	if err := d.Order("cpu").Find(&configs).Error; err != nil {
		return nil, err
	}
	for _, config := range configs {
		configDTO := new(dto.VmConfig)
		configDTO.VmConfig = config
		configDTOS = append(configDTOS, *configDTO)
	}
	return configDTOS, nil
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
	if err != nil && !gorm.IsRecordNotFoundError(err) {
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

func (v vmConfigService) Update(name string, update dto.VmConfigUpdate) (*dto.VmConfig, error) {

	vmConfig, err := v.Get(name)
	if err != nil {
		return nil, err
	}
	if update.Name != "" {
		vmConfig.Name = update.Name
	}
	if update.Memory != 0 {
		vmConfig.Memory = update.Memory
	}
	if update.Cpu != 0 {
		vmConfig.Cpu = update.Cpu
	}
	if err := db.DB.Save(&vmConfig).Error; err != nil {
		return nil, err
	}
	return vmConfig, err
}

func (v vmConfigService) Delete(name string) error {
	vmConfig, err := v.Get(name)
	if err != nil {
		return err
	}
	if err := db.DB.Delete(&vmConfig).Error; err != nil {
		return err
	}
	return nil
}
