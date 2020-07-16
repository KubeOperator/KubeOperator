package repository

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

var (
	DeleteHostFailed = "DELETE_HOST_FAILED"
)

type HostRepository interface {
	Get(name string) (model.Host, error)
	List() ([]model.Host, error)
	ListByClusterId(clusterId string) ([]model.Host, error)
	Page(num, size int) (int, []model.Host, error)
	Save(host *model.Host) error
	BatchSave(hosts []*model.Host) error
	Delete(name string) error
	ListByCredentialID(credentialID string) ([]model.Host, error)
	Batch(operation string, items []model.Host) error
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
	err := db.DB.Model(model.Host{}).Preload("Volumes").Preload("Cluster").Find(&hosts).Error
	return hosts, err
}

func (h hostRepository) Page(num, size int) (int, []model.Host, error) {
	var total int
	var hosts []model.Host
	err := db.DB.Model(model.Host{}).
		Count(&total).
		Preload("Volumes").
		Preload("Cluster").
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

func (h hostRepository) ListByClusterId(clusterId string) ([]model.Host, error) {
	var cluster model.Cluster
	var hosts []model.Host
	cluster.ID = clusterId
	if err := db.DB.First(&cluster).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Where(model.Host{ClusterID: clusterId}).Find(&hosts).Error; err != nil {
		return nil, err
	}
	return hosts, nil

}

func (h hostRepository) BatchSave(hosts []*model.Host) error {
	tx := db.DB.Begin()
	for i, _ := range hosts {
		if db.DB.NewRecord(hosts[i]) {
			if err := tx.Create(hosts[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.Save(hosts[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (h hostRepository) Delete(name string) error {
	host, err := h.Get(name)
	if err != nil {
		return err
	}
	if host.ClusterID != "" {
		return errors.New(DeleteHostFailed)
	}
	return db.DB.Delete(&host).Error
}

func (h hostRepository) ListByCredentialID(credentialID string) ([]model.Host, error) {
	var host []model.Host
	err := db.DB.Model(model.Host{
		CredentialID: credentialID,
	}).Find(&host).Error
	return host, err
}

func (h hostRepository) Batch(operation string, items []model.Host) error {
	switch operation {
	case constant.BatchOperationDelete:
		var clusterIds []string
		for _, item := range items {
			clusterIds = append(clusterIds, item.ClusterID)
		}
		var clusters []model.Cluster
		err := db.DB.Where("id in (?)", clusterIds).Find(&clusters).Error
		if err != nil {
			return err
		}
		if len(clusters) > 0 {
			return errors.New(DeleteFailedError)
		}
		var ids []string
		for _, item := range items {
			ids = append(ids, item.ID)
		}
		err = db.DB.Where("id in (?)", ids).Delete(&items).Error
		if err != nil {
			return err
		}
	default:
		return constant.NotSupportedBatchOperation
	}
	return nil
}
