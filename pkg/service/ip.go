package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/ipaddr"
	"github.com/jinzhu/gorm"
)

type IpService interface {
	Get(ip string) (dto.Ip, error)
	Create(create dto.IpCreate, tx *gorm.DB) error
	Page(num, size int, ipPoolName string) (page.Page, error)
	Batch(op dto.IpOp) error
	Update(update dto.IpUpdate) (*dto.Ip, error)
	Sync(ipPoolName string) error
}

type ipService struct {
}

func NewIpService() IpService {
	return &ipService{}
}

func (i ipService) Get(ip string) (dto.Ip, error) {
	var ipDTO dto.Ip
	var ipM model.Ip
	if err := db.DB.Where("address = ?", ip).First(&ipM).Error; err != nil {
		return ipDTO, err
	}
	ipDTO = dto.Ip{
		Ip: ipM,
	}
	return ipDTO, nil
}

func (i ipService) Create(create dto.IpCreate, tx *gorm.DB) error {
	if tx == nil {
		tx = db.DB.Begin()
	}
	var ipPool model.IpPool
	if err := tx.Where("name = ?", create.IpPoolName).First(&ipPool).Error; err != nil {
		return err
	}
	cs := strings.Split(create.Subnet, "/")
	mask, _ := strconv.Atoi(cs[1])
	startIp := strings.Replace(create.IpStart, " ", "", -1)
	endIp := strings.Replace(create.IpEnd, " ", "", -1)
	ips := ipaddr.GenerateIps(cs[0], mask, startIp, endIp)
	for _, ip := range ips {
		var old model.Ip
		tx.Where("address = ?", ip).First(&old)
		if old.ID != "" {
			tx.Rollback()
			return errors.New("IP_EXISTS")
		}
		insert := model.Ip{
			Address:  ip,
			Gateway:  create.Gateway,
			DNS1:     create.DNS1,
			DNS2:     create.DNS2,
			IpPoolID: ipPool.ID,
			Status:   constant.IpAvailable,
		}
		err := tx.Create(&insert).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		go func() {
			err := ipaddr.Ping(insert.Address)
			if err == nil {
				insert.Status = constant.IpReachable
				db.DB.Save(&insert)
			}
		}()
	}
	tx.Commit()
	return nil
}

func (i ipService) Page(num, size int, ipPoolName string) (page.Page, error) {
	var (
		page   page.Page
		ipDTOS []dto.Ip
		total  int
		ips    []model.Ip
		ipPool model.IpPool
	)
	err := db.DB.Where("name = ?", ipPoolName).First(&ipPool).Error
	if err != nil {
		return page, err
	}
	err = db.DB.Model(&model.Ip{}).Where("ip_pool_id = ?", ipPool.ID).Order("inet_aton(address)").Count(&total).Offset((num - 1) * size).Limit(size).Find(&ips).Error
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
			if err := tx.Where("address = ?", op.Items[i].Address).First(&ip).Error; err != nil {
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

func (i ipService) Update(update dto.IpUpdate) (*dto.Ip, error) {
	tx := db.DB.Begin()
	var ip model.Ip
	err := tx.Where("address = ?", update.Address).First(&ip).Error
	if err != nil {
		return nil, err
	}
	switch update.Operation {
	case "LOCK":
		ip.Status = constant.IpLock
		break
	case "UNLOCK":
		ip.Status = constant.IpAvailable
		break
	default:
		break
	}
	err = tx.Save(&ip).Error
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return &dto.Ip{Ip: ip}, err
}

func (i ipService) Sync(ipPoolName string) error {
	var ipPool model.IpPool
	err := db.DB.Where("name = ?", ipPoolName).First(&ipPool).Error
	if err != nil {
		return err
	}
	var ips []model.Ip
	err = db.DB.Where("ip_pool_id = ?", ipPool.ID).Find(&ips).Error
	if err != nil {
		return err
	}
	for i := range ips {
		if ips[i].Status != constant.IpUsed && ips[i].Status != constant.IpLock {
			go func(i int) {
				err := ipaddr.Ping(ips[i].Address)
				if err == nil && ips[i].Status != constant.IpReachable {
					ips[i].Status = constant.IpReachable
					db.DB.Save(&ips[i])
				}
				if err != nil && ips[i].Status == constant.IpReachable {
					ips[i].Status = constant.IpAvailable
					db.DB.Save(&ips[i])
				}
			}(i)
		}
	}
	return nil
}
