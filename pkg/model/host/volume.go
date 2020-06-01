package host

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Volume struct {
	common.BaseModel
	ID     string
	HostID string
	Size   string
	Name   string
}

func (v *Volume) BeforeCreate() (err error) {
	v.ID = uuid.NewV4().String()
	return err
}

func (v Volume) TableName() string {
	return "ko_volume"
}
