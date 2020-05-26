package credential

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Credential struct {
	common.BaseModel
	ID         string
	Name       string `gorm:"not null;unique"`
	Username   string
	Password   string
	PrivateKey string
	Type       string
}

func (c *Credential) BeforeCreate() (err error) {
	c.ID = uuid.NewV4().String()
	return err
}

func (c Credential) TableName() string {
	return "ko_credential"
}
