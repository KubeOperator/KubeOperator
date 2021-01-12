package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
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
	List() ([]dto.IpPool, error)
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
	ipPoolDTO.IpPool = ipPool
	return ipPoolDTO, nil
}

func (i ipPoolService) Page(num, size int) (page.Page, error) {
	var page page.Page
	var ipPoolDTOS []dto.IpPool
	var total int
	var ipPools []model.IpPool
	err := db.DB.Model(model.IpPool{}).Preload("Ips").Count(&total).Find(&ipPools).Offset((num - 1) * size).Limit(size).Error
	if err != nil {
		return page, err
	}
	for _, mo := range ipPools {
		ipUsed := 0
		for _, ip := range mo.Ips {
			if ip.Status != constant.IpAvailable {
				ipUsed++
			}
		}
		ipPoolDTOS = append(ipPoolDTOS, dto.IpPool{
			IpPool: mo,
			IpUsed: ipUsed,
		})
	}
	page.Total = total
	page.Items = ipPoolDTOS
	return page, err
}

func (i ipPoolService) List() ([]dto.IpPool, error) {
	var ipPoolDTOS []dto.IpPool
	var ipPools []model.IpPool
	err := db.DB.Model(model.IpPool{}).Preload("Ips").Find(&ipPools).Error
	if err != nil {
		return ipPoolDTOS, err
	}
	for _, mo := range ipPools {
		ipUsed := 0
		for _, ip := range mo.Ips {
			if ip.Status != constant.IpAvailable {
				ipUsed++
			}
		}
		ipPoolDTOS = append(ipPoolDTOS, dto.IpPool{
			IpPool: mo,
			IpUsed: ipUsed,
		})
	}
	return ipPoolDTOS, nil
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
		IpStart:    creation.IpStart,
		IpEnd:      creation.IpEnd,
		Gateway:    creation.Gateway,
		Subnet:     creation.Subnet,
		IpPoolName: ipPool.Name,
		DNS1:       creation.DNS1,
		DNS2:       creation.DNS2,
	}, tx)
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
	tx := db.DB.Begin()
	switch op.Operation {
	case constant.BatchOperationDelete:
		for i := range opItems {
			var ipPool model.IpPool
			if err := tx.Where(model.IpPool{Name: opItems[i].Name}).First(&ipPool).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Delete(&ipPool).Error; err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Where("ip_pool_id = ?", ipPool.ID).Delete(&model.Ip{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	default:
		return constant.NotSupportedBatchOperation
	}
	tx.Commit()
	return nil
}
