package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type IpPoolService interface {
	Get(name string) (dto.IpPool, error)
	Page(num, size int) (page.Page, error)
	Create(creation dto.IpPoolCreate) (dto.IpPool, error)
	Batch(op dto.IpPoolOp) error
}

type ipPoolService struct {
	ipPoolRepo repository.IpPoolRepository
	ipService  IpService
}

func NewIpPoolService() IpPoolService {
	return &ipPoolService{
		ipPoolRepo: repository.NewIpPoolRepository(),
		ipService:  NewIpService(),
	}
}

func (i ipPoolService) Get(name string) (dto.IpPool, error) {
	var ipPoolDTO dto.IpPool
	ipPool, err := i.ipPoolRepo.Get(name)
	if err != nil {
		return ipPoolDTO, err
	}
	ipPoolDTO = dto.IpPool{
		ipPool,
	}
	return ipPoolDTO, nil
}

func (i ipPoolService) Page(num, size int) (page.Page, error) {
	var page page.Page
	var ipPoolDTOS []dto.IpPool
	total, mos, err := i.ipPoolRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		ipPoolDTOS = append(ipPoolDTOS, dto.IpPool{
			IpPool: mo,
		})
	}
	page.Total = total
	page.Items = ipPoolDTOS
	return page, err
}

func (i ipPoolService) Create(creation dto.IpPoolCreate) (dto.IpPool, error) {
	var ipPoolDTO dto.IpPool
	ipPool := model.IpPool{
		BaseModel:   common.BaseModel{},
		Name:        creation.Name,
		Description: creation.Description,
		Subnet:      creation.Subnet,
	}
	tx := db.DB.Begin()
	err := tx.Create(&ipPool).Error
	if err != nil {
		tx.Rollback()
		return ipPoolDTO, err
	}
	err = i.ipService.Create(dto.IpCreate{
		StartIp:  creation.IpStart,
		EndIp:    creation.IpEnd,
		Gateway:  creation.Gateway,
		Subnet:   creation.Subnet,
		IpPoolID: ipPool.ID,
		DNS1:     creation.DNS1,
		DNS2:     creation.DNS2,
	})
	if err != nil {
		tx.Rollback()
		return ipPoolDTO, err
	}
	tx.Commit()
	ipPoolDTO.IpPool = ipPool
	return ipPoolDTO, err
}

func (i ipPoolService) Batch(op dto.IpPoolOp) error {
	var opItems []model.IpPool
	for _, item := range op.Items {
		opItems = append(opItems, model.IpPool{
			BaseModel: common.BaseModel{},
			Name:      item.Name,
		})
	}
	return i.ipPoolRepo.Batch(op.Operation, opItems)
}
