package service

import (
	"errors"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"

	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/jinzhu/gorm"
)

type SystemSettingService interface {
	Get(name string) (dto.SystemSetting, error)
	GetLocalIPs() ([]model.SystemRegistry, error)
	List() (dto.SystemSettingResult, error)
	Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error)
	ListByTab(tabName string) (dto.SystemSettingResult, error)
	ListRegistry() ([]dto.SystemRegistry, error)
	GetRegistryByArch(arch string) (dto.SystemRegistry, error)
	CreateRegistry(creation dto.SystemRegistryCreate) (*dto.SystemRegistry, error)
	UpdateRegistry(creation dto.SystemRegistryUpdate) (*dto.SystemRegistry, error)
	PageRegistry(num, size int) (page.Page, error)
	BatchRegistry(op dto.SystemRegistryBatchOp) error
}

type systemSettingService struct {
	systemSettingRepo  repository.SystemSettingRepository
	systemRegistryRepo repository.SystemRegistryRepository
	userRepo           repository.UserRepository
}

func NewSystemSettingService() SystemSettingService {
	return &systemSettingService{
		systemSettingRepo:  repository.NewSystemSettingRepository(),
		systemRegistryRepo: repository.NewSystemRegistryRepository(),
		userRepo:           repository.NewUserRepository(),
	}
}

func (s systemSettingService) Get(key string) (dto.SystemSetting, error) {
	var systemSettingDTO dto.SystemSetting
	mo, err := s.systemSettingRepo.Get(key)
	if err != nil {
		return systemSettingDTO, err
	}
	systemSettingDTO.SystemSetting = mo
	return systemSettingDTO, err
}

func (s systemSettingService) List() (dto.SystemSettingResult, error) {
	var systemSettingResult dto.SystemSettingResult
	vars := make(map[string]string)
	mos, err := s.systemSettingRepo.List()
	if err != nil {
		return systemSettingResult, err
	}
	for _, mo := range mos {
		vars[mo.Key] = mo.Value
	}
	systemSettingResult.Vars = vars
	return systemSettingResult, err
}

func (s systemSettingService) ListByTab(tabName string) (dto.SystemSettingResult, error) {
	var systemSettingResult dto.SystemSettingResult
	vars := make(map[string]string)
	mos, err := s.systemSettingRepo.ListByTab(tabName)
	if err != nil {
		return systemSettingResult, err
	}
	for _, mo := range mos {
		vars[mo.Key] = mo.Value
	}
	if len(mos) > 0 {
		systemSettingResult.Tab = tabName
	}
	systemSettingResult.Vars = vars
	return systemSettingResult, err
}

func (s systemSettingService) Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error) {

	var result []dto.SystemSetting
	for k, v := range creation.Vars {
		systemSetting, err := s.systemSettingRepo.Get(k)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				systemSetting.Key = k
				systemSetting.Value = v
				systemSetting.Tab = creation.Tab
				err := s.systemSettingRepo.Save(&systemSetting)
				if err != nil {
					return result, err
				}
				result = append(result, dto.SystemSetting{SystemSetting: systemSetting})
			} else {
				return result, err
			}
		} else if systemSetting.ID != "" {
			systemSetting.Value = v
			if systemSetting.Tab == "" {
				systemSetting.Tab = creation.Tab
			}
			err := s.systemSettingRepo.Save(&systemSetting)
			if err != nil {
				return result, err
			}
			result = append(result, dto.SystemSetting{SystemSetting: systemSetting})
		}
	}
	return result, nil
}

func (s systemSettingService) GetLocalIPs() ([]model.SystemRegistry, error) {
	var sysRepo []model.SystemRegistry
	if err := db.DB.Find(&sysRepo).Error; err != nil {
		return sysRepo, fmt.Errorf("can't found repo from system registry, err %s", err.Error())
	}
	return sysRepo, nil
}

func (s systemSettingService) ListRegistry() ([]dto.SystemRegistry, error) {
	var systemRegistryDto []dto.SystemRegistry
	mos, err := s.systemRegistryRepo.List()
	if err != nil {
		return systemRegistryDto, err
	}
	for _, mo := range mos {
		systemRegistryDto = append(systemRegistryDto, dto.SystemRegistry{
			SystemRegistry: mo,
		})
	}
	return systemRegistryDto, nil
}

func (s systemSettingService) GetRegistryByArch(arch string) (dto.SystemRegistry, error) {
	r, err := s.systemRegistryRepo.Get(arch)
	if err != nil {
		return dto.SystemRegistry{}, err
	}
	systemRegistryDto := dto.SystemRegistry{
		SystemRegistry: model.SystemRegistry{
			ID:           r.ID,
			Hostname:     r.Hostname,
			Protocol:     r.Protocol,
			Architecture: r.Architecture,
		},
	}
	return systemRegistryDto, nil
}

func (s systemSettingService) CreateRegistry(creation dto.SystemRegistryCreate) (*dto.SystemRegistry, error) {
	old, err := s.systemRegistryRepo.Get(creation.Architecture)
	if !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if old.ID != "" {
		return nil, errors.New("NAME_EXISTS")
	}

	systemRegistry := model.SystemRegistry{
		ID:           old.ID,
		Architecture: creation.Architecture,
		Protocol:     creation.Protocol,
		Hostname:     creation.Hostname,
	}
	if err := s.systemRegistryRepo.Save(&systemRegistry); err != nil {
		return nil, err
	}
	return &dto.SystemRegistry{SystemRegistry: systemRegistry}, nil
}

func (s systemSettingService) UpdateRegistry(creation dto.SystemRegistryUpdate) (*dto.SystemRegistry, error) {
	old, err := s.systemRegistryRepo.Get(creation.Architecture)
	if err != nil {
		return nil, err
	}
	if old.ID == "" || old.ID != creation.ID {
		return nil, errors.New("NOT_FOUND")
	}

	systemRegistry := model.SystemRegistry{
		ID:           old.ID,
		Architecture: creation.Architecture,
		Protocol:     creation.Protocol,
		Hostname:     creation.Hostname,
	}
	if err := s.systemRegistryRepo.Save(&systemRegistry); err != nil {
		return nil, err
	}
	return &dto.SystemRegistry{SystemRegistry: systemRegistry}, nil
}

func (s systemSettingService) PageRegistry(num, size int) (page.Page, error) {
	var page page.Page
	var systemRegistryDto []dto.SystemRegistry
	total, mos, err := s.systemRegistryRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		systemRegistryDto = append(systemRegistryDto, dto.SystemRegistry{SystemRegistry: mo})
	}
	page.Total = total
	page.Items = systemRegistryDto
	return page, err
}

func (s systemSettingService) BatchRegistry(op dto.SystemRegistryBatchOp) error {
	var deleteItems []model.SystemRegistry
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.SystemRegistry{
			BaseModel:    common.BaseModel{},
			ID:           item.ID,
			Architecture: item.Architecture,
		})
	}
	err := s.systemRegistryRepo.Batch(op.Operation, deleteItems)
	if err != nil {
		return err
	}
	return nil
}
