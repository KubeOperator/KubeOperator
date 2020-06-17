package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type HostRepository interface {
	Get(name string) (model.Host, error)
	List() ([]model.Host, error)
	Page(num, size int) (int, []model.Host, error)
	Save(host *model.Host) error
	Delete(name string) error
	ListByCredentialID(credentialID string) ([]model.Host, error)
}

func NewHostRepository() HostRepository {
	return &hostRepository{}
}

type hostRepository struct {
}

func (h hostRepository) Get(name string) (model.Host, error) {
	var host model.Host
	host.Name = name
	if err := db.DB.Where(host).First(&host).Error; err != nil {
		return host, err
	}
	if err := db.DB.First(&host).Related(&host.Volumes).Error; err != nil {
		return host, err
	}
	if err := db.DB.First(&host).Related(&host.Credential).Error; err != nil {
		return host, err
	}
	return host, nil
}

func (h hostRepository) List() ([]model.Host, error) {
	var hosts []model.Host
	err := db.DB.Model(model.Host{}).Preload("Volumes").Find(&hosts).Error
	return hosts, err
}

func (h hostRepository) Page(num, size int) (int, []model.Host, error) {
	var total int
	var hosts []model.Host
	err := db.DB.Model(model.Host{}).
		Count(&total).
		Preload("Volumes").
		Find(&hosts).
		Offset((num - 1) * size).
		Limit(size).
		Error
	return total, hosts, err
}

func (h hostRepository) Save(host *model.Host) error {
	if db.DB.NewRecord(host) {
		return db.DB.Create(&host).Error
	} else {
		err := db.DB.Where(model.Volume{HostID: host.ID}).Delete(model.Volume{}).Error
		if err != nil {
			return err
		}
		return db.DB.Save(&host).Error
	}
}

func (h hostRepository) Delete(name string) error {
	var host model.Host
	host.Name = name
	return db.DB.Delete(&host).Error
}

func (h hostRepository) ListByCredentialID(credentialID string) ([]model.Host, error) {
	var host []model.Host
	err := db.DB.Model(model.Host{
		CredentialID: credentialID,
	}).Find(&host).Error
	return host, err
}
