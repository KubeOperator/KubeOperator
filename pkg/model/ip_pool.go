package model

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type IpPool struct {
	common.BaseModel
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Subnet      string `json:"subnet"`
	Ips         []Ip   `json:"ips"`
}

func (i *IpPool) BeforeCreate() (err error) {
	i.ID = uuid.NewV4().String()
	return err
}

func (i *IpPool) BeforeDelete() (err error) {
	var zones []Zone
	db.DB.Where(Zone{IpPoolID: i.ID}).Find(&zones)
	if len(zones) > 0 {
		return errors.New("IP_POOL_DELETE_FAILED")
	}
	return nil
}
