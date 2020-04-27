package host

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"ko3-gin/internal/db"
	"ko3-gin/pkg/model"
)

var Service = &service{}

type service struct{}

func (d *service) Create(host model.Host) error {
	host.Id = uuid.NewV4().String()
	if err := db.DB.Create(&host).Error; err != nil {
		return err
	}
	return nil
}

func (d *service) Update(host model.Host) error {
	if err := db.DB.Update(&host).Error; err != nil {
		return err
	}
	return nil
}

func (d *service) List() []model.Host {
	var hosts []model.Host
	db.DB.Find(&hosts)
	return hosts
}

func (d *service) Page(num, size int) (hosts []model.Host, total int) {
	db.DB.Model(model.Host{}).
		Find(&hosts).
		Offset((num - 1) * size).
		Limit(size).
		Count(total)
	return
}

func (d *service) Get(id string) (h *model.Host, error error) {
	var host model.Host
	if db.DB.First(&host, id).RecordNotFound() {
		return nil, gorm.ErrRecordNotFound
	}
	return &host, nil
}

func (d *service) Delete(id string) error {
	host, err := d.Get(id)
	if err != nil {
		return err
	}
	if err := db.DB.Delete(host).Error; err != nil {
		return err
	}
	return nil
}
