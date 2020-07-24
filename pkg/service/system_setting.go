package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/jinzhu/gorm"
)

type SystemSettingService interface {
	Get(name string) (dto.SystemSetting, error)
	GetLocalHostName() string
	List() ([]dto.SystemSettingResult, error)
	Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error)
}

type systemSettingService struct {
	systemSettingRepo repository.SystemSettingRepository
}

func NewSystemSettingService() SystemSettingService {
	return &systemSettingService{
		systemSettingRepo: repository.NewSystemSettingRepository(),
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

func (s systemSettingService) List() ([]dto.SystemSettingResult, error) {
	var systemSettingResult []dto.SystemSettingResult
	vars := make(map[string]string)
	mos, err := s.systemSettingRepo.List()
	if err != nil {
		return systemSettingResult, err
	}
	for _, mo := range mos {
		vars[mo.Key] = mo.Value
	}
	systemSettingResult = append(systemSettingResult, dto.SystemSettingResult{Vars: vars})
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
			err := s.systemSettingRepo.Save(&systemSetting)
			if err != nil {
				return result, err
			}
			result = append(result, dto.SystemSetting{SystemSetting: systemSetting})
		}
	}
	return result, nil
}

func (s systemSettingService) GetLocalHostName() string {
	mo, err := s.systemSettingRepo.Get("ip")
	if err != nil || mo.Value == "" {
		return ""
	}
	return mo.Value
}
