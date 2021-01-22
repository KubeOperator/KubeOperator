package repository

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

var (
	DeleteHostFailed          = "DELETE_HOST_FAILED"
	DeleteHostFailedByProject = "DELETE_HOST_FAILED_BY_PROJECT"
)

type HostRepository interface {
	Get(name string) (model.Host, error)
	List(projectName string) ([]model.Host, error)
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

func (h hostRepository) List(projectName string) ([]model.Host, error) {
	var hosts []model.Host
	if projectName == "" {
		err := db.DB.Model(&model.Host{}).
			Preload("Volumes").
			Preload("Cluster").
			Preload("Zone").
			Find(&hosts).
			Error
		return hosts, err
	} else {
		var project model.Project
		err := db.DB.Model(&model.Project{}).Where(&model.Project{Name: projectName}).First(&project).Error
		if err != nil {
			return nil, err
		}
		var projectResources []model.ProjectResource
		err = db.DB.Model(&model.ProjectResource{}).Where(&model.ProjectResource{ProjectID: project.ID, ResourceType: constant.ResourceHost}).Find(&projectResources).Error
		if err != nil {
			return nil, err
		}
		var resourceIds []string
		for _, pr := range projectResources {
			resourceIds = append(resourceIds, pr.ResourceID)
		}
		err = db.DB.Model(&model.Host{}).Where("id in (?)", resourceIds).Find(&hosts).Error
		return hosts, err
	}
}

func (h hostRepository) Page(num, size int) (int, []model.Host, error) {
	var total int
	var hosts []model.Host
	err := db.DB.Model(&model.Host{}).
		Order("name asc").
		Count(&total).
		Preload("Volumes").
		Preload("Cluster").
		Preload("Zone").
		Find(&hosts).
		Offset((num - 1) * size).
		Limit(size).
		Error
	return total, hosts, err
}

func (h hostRepository) Save(host *model.Host) error {
	if host.Name == "" {
		return nil
	}
	tx := db.DB.Begin()
	if db.DB.NewRecord(host) {
		if err := tx.Create(&host).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// lock
		var lock model.Host
		lock.ID = host.ID
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&lock).Error; err != nil {
			return err
		}
		if len(host.Volumes) > 0 {
			for i := range host.Volumes {
				var volume model.Volume
				if notFound := tx.Where(&model.Volume{HostID: host.ID, Name: host.Volumes[i].Name}).
					First(&volume).RecordNotFound(); notFound {
					if err := tx.Create(&host.Volumes[i]).Error; err != nil {
						tx.Rollback()
						return err
					}
				} else {
					host.Volumes[i].ID = volume.ID
					if err := tx.Save(&host.Volumes[i]).Error; err != nil {
						tx.Rollback()
						return err
					}

				}
			}
		}
		if err := tx.Save(&host).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	var ip model.Ip
	tx.Where(&model.Ip{Address: host.Ip}).First(&ip)
	if ip.ID != "" && ip.Status != constant.IpUsed {
		ip.Status = constant.IpUsed
		if err := tx.Save(&ip).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	// unlock
	tx.Commit()
	return nil
}

func (h hostRepository) ListByClusterId(clusterId string) ([]model.Host, error) {
	var cluster model.Cluster
	var hosts []model.Host
	cluster.ID = clusterId
	if err := db.DB.First(&cluster).Error; err != nil {
		return nil, err
	}
	if err := db.DB.Where(&model.Host{ClusterID: clusterId}).Find(&hosts).Error; err != nil {
		return nil, err
	}
	return hosts, nil

}

func (h hostRepository) BatchSave(hosts []*model.Host) error {
	tx := db.DB.Begin()
	for i := range hosts {
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
		var ip model.Ip
		tx.Where(&model.Ip{Address: hosts[i].Ip}).First(&ip)
		if ip.ID != "" && ip.Status != constant.IpUsed {
			ip.Status = constant.IpUsed
			if err := tx.Save(&ip).Error; err != nil {
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
	err := db.DB.Model(&model.Host{
		CredentialID: credentialID,
	}).Find(&host).Error
	return host, err
}

func (h hostRepository) Batch(operation string, items []model.Host) error {

	tx := db.DB.Begin()
	switch operation {
	case constant.BatchOperationDelete:
		for i := range items {

			var host model.Host
			if err := db.DB.Where(&model.Host{Name: items[i].Name}).First(&host).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := db.DB.Delete(&host).Error; err != nil {
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
