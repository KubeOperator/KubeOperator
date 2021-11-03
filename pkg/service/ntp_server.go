package service

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

var (
	DeleteNtpServerError = "DELETE_REGION_FAILED_RESOURCE"
	NtpServerNameExist   = "NAME_EXISTS"
)

type NtpServerService interface {
	Get(name string) (dto.NtpServer, error)
	List() ([]dto.NtpServer, error)
	Page(num, size int) (*page.Page, error)
	Delete(name string) error
	Create(creation dto.NtpServerCreate) (*dto.NtpServer, error)
	Update(name string, update dto.NtpServerUpdate) (*dto.NtpServer, error)
}

type ntpServerService struct {
	ntpServerRepo repository.NtpServerRepository
}

func NewNtpServerService() NtpServerService {
	return &ntpServerService{
		ntpServerRepo: repository.NewNtpServerRepository(),
	}
}

func (n ntpServerService) Get(name string) (dto.NtpServer, error) {
	var ntpServerDTO dto.NtpServer
	mo, err := n.ntpServerRepo.Get(name)
	if err != nil {
		return ntpServerDTO, err
	}
	ntpServerDTO.NtpServer = mo
	return ntpServerDTO, err
}

func (n ntpServerService) List() ([]dto.NtpServer, error) {
	var ntpServers []dto.NtpServer
	datas, err := n.ntpServerRepo.List()
	if err != nil {
		return nil, err
	}
	for _, n := range datas {
		ntpServers = append(ntpServers, dto.NtpServer{NtpServer: n})
	}
	return ntpServers, nil
}

func (n ntpServerService) Page(num, size int) (*page.Page, error) {
	var p page.Page

	total, datas, err := n.ntpServerRepo.Page(num, size)
	if err != nil {
		return nil, err
	}
	p.Total = total
	p.Items = datas
	return &p, nil
}

func (n ntpServerService) Delete(name string) error {
	return n.ntpServerRepo.Delete(name)
}

func (n ntpServerService) Create(creation dto.NtpServerCreate) (*dto.NtpServer, error) {
	old, _ := n.Get(creation.Name)
	if old.ID != "" {
		return nil, errors.New(NtpServerNameExist)
	}

	ntpServer := model.NtpServer{
		BaseModel: common.BaseModel{},
		Name:      creation.Name,
		Address:   creation.Address,
		Status:    creation.Status,
	}

	err := n.ntpServerRepo.Save(&ntpServer)
	if err != nil {
		return nil, err
	}
	return &dto.NtpServer{NtpServer: ntpServer}, err
}

func (n ntpServerService) Update(name string, update dto.NtpServerUpdate) (*dto.NtpServer, error) {
	ntpServer, err := n.ntpServerRepo.Get(name)
	if err != nil {
		return nil, err
	}
	ntpServer.Status = update.Status
	ntpServer.Address = update.Address
	if err := n.ntpServerRepo.Save(&ntpServer); err != nil {
		return nil, err
	}
	return &dto.NtpServer{NtpServer: ntpServer}, nil
}
