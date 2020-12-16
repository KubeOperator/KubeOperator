package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type IpPool struct {
	common.BaseModel
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Subnet      string `json:"subnet"`
}

func (i *IpPool) BeforeCreate() (err error) {
	i.ID = uuid.NewV4().String()
	return err
}
