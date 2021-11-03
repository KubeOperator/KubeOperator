package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type NtpServerRepository interface {
	Get(name string) (model.NtpServer, error)
	List() ([]model.NtpServer, error)
	GetAddressStr() (string, error)
	Page(num, size int) (int, []model.NtpServer, error)
	Save(ntpServer *model.NtpServer) error
	Delete(name string) error
}

func NewNtpServerRepository() NtpServerRepository {
	return &ntpServerRepository{}
}

type ntpServerRepository struct {
}

func (n ntpServerRepository) Get(name string) (model.NtpServer, error) {
	var ntpServer model.NtpServer
	if err := db.DB.Where("name = ?", name).First(&ntpServer).Error; err != nil {
		return ntpServer, err
	}
	return ntpServer, nil
}

func (n ntpServerRepository) GetAddressStr() (string, error) {
	addressStr := ""
	var datas []model.NtpServer
	if err := db.DB.Where("status = ?", "enable").Find(&datas).Error; err != nil {
		return addressStr, err
	}
	for _, n := range datas {
		addressStr += (n.Address + ",")
	}
	if len(addressStr) != 0 {
		addressStr = addressStr[0 : len(addressStr)-1]
	}
	return addressStr, nil
}

func (n ntpServerRepository) List() ([]model.NtpServer, error) {
	var ntpServer []model.NtpServer
	err := db.DB.Find(&ntpServer).Error
	return ntpServer, err
}

func (n ntpServerRepository) Page(num, size int) (int, []model.NtpServer, error) {
	var total int
	var ntpServer []model.NtpServer
	err := db.DB.Model(&model.NtpServer{}).Order("created_at desc").Count(&total).Find(&ntpServer).Offset((num - 1) * size).Limit(size).Error
	return total, ntpServer, err
}

func (n ntpServerRepository) Save(ntpServer *model.NtpServer) error {
	if db.DB.NewRecord(ntpServer) {
		return db.DB.Create(&ntpServer).Error
	} else {
		return db.DB.Save(&ntpServer).Error
	}
}

func (n ntpServerRepository) Delete(name string) error {
	return db.DB.Where("name = ?", name).Delete(&model.NtpServer{}).Error
}
