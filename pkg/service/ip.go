package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/ipaddr"
	"strconv"
	"strings"
)

type IpService interface {
	Get(ip string) (dto.Ip, error)
	Create(create dto.IpCreate) error
	Page(num, size int, ipPoolName string) (page.Page, error)
	Batch(op dto.IpOp) error
}

type ipService struct {
}

func NewIpService() IpService {
	return &ipService{}
}

func (i ipService) Get(ip string) (dto.Ip, error) {
	var ipDTO dto.Ip
	var ipM model.Ip
	err := db.DB.Where(model.Ip{Address: ip}).First(&ipM).Error
	if err != nil {
		return ipDTO, err
	}
	ipDTO = dto.Ip{
		Ip: ipM,
	}
	return ipDTO, nil
}

func (i ipService) Create(create dto.IpCreate) error {
	var ipPool model.IpPool
	err := db.DB.Where(model.IpPool{Name: create.IpPoolName}).First(&ipPool).Error
	if err != nil {
		return err
	}
	cs := strings.Split(create.Subnet, "/")
	mask, _ := strconv.Atoi(cs[1])
	ips := ipaddr.GenerateIps(cs[0], mask, create.StartIp, create.EndIp)
	tx := db.DB.Begin()
	for _, ip := range ips {
		err := tx.Create(&model.Ip{
			Address:  ip,
			Gateway:  create.Gateway,
			DNS1:     create.DNS1,
			DNS2:     create.DNS2,
			IpPoolID: ipPool.ID,
			Status:   constant.IpAvailable,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (i ipService) Page(num, size int, ipPoolName string) (page.Page, error) {
	var page page.Page
	var ipDTOS []dto.Ip
	var total int
	var ips []model.Ip
	var ipPool model.IpPool
	err := db.DB.Where(model.IpPool{Name: ipPoolName}).First(&ipPool).Error
	if err != nil {
		return page, err
	}
	err = db.DB.Model(model.Ip{}).Where(model.Ip{IpPoolID: ipPool.ID}).Order("inet_aton(address)").Count(&total).Find(&ips).Offset((num - 1) * size).Limit(size).Error
	if err != nil {
		return page, err
	}
	for _, mo := range ips {
		ipDTOS = append(ipDTOS, dto.Ip{
			Ip: mo,
		})
	}
	page.Total = total
	page.Items = ipDTOS
	return page, nil
}

func (i ipService) Batch(op dto.IpOp) error {
	tx := db.DB.Begin()
	switch op.Operation {
	case constant.BatchOperationDelete:
		for i := range op.Items {
			var ip model.Ip
			if err := tx.Where(model.Ip{Address: op.Items[i].Address}).First(&ip).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Delete(&ip).Error; err != nil {
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
