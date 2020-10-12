package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type VmConfigService interface {
	Page(num, size int) (page.Page, error)
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
