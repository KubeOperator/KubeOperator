package host

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"ko3-gin/internal/db"
)

var Service = &serviceManger{}

type serviceManger struct{}

func (d *serviceManger) Create(host Host) error {
	host.Id = uuid.NewV4().String()
	if err := db.DB.Create(&host).Error; err != nil {
		return err
	}
	return nil
}

func (d *serviceManger) Update(host Host) error {
	if err := db.DB.Update(&host).Error; err != nil {
		return err
	}
	return nil
}

func (d *serviceManger) List() []Host {
	var hosts []Host
	db.DB.Find(&hosts)
	return hosts
}

func (d *serviceManger) Page(num, size int) (hosts []Host, total int) {
	db.DB.Model(Host{}).
		Find(&hosts).
		Offset((num - 1) * size).
		Limit(size).
		Count(total)
	return
}

func (d *serviceManger) Get(id string) (h *Host, error error) {
	var host Host
	if db.DB.First(&host, id).RecordNotFound() {
		return nil, gorm.ErrRecordNotFound
	}
	return &host, nil
}

func (d *serviceManger) Delete(id string) error {
	host, err := d.Get(id)
	if err != nil {
		return err
	}
	if err := db.DB.Delete(host).Error; err != nil {
		return err
	}
	return nil
}
