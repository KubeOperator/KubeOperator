package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
)

type HostService interface {
	Get(name string) (dto.Host, error)
	List() ([]dto.Host, error)
	Page(num, size int) (int, []dto.Host, error)
	Create(creation dto.Host) error
	Delete(name string) error
}

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

func (h hostService) Create(creation dto.Host) error {
	host := model.Host{
		BaseModel: common.BaseModel{},
		Name:      creation.Name,
		Ip:        creation.Ip,
		Port:      creation.Port,
	}
	err := h.hostRepo.Save(&host)
	if err != nil {
		return err
	}
	return err
}
