package credential

import (
	uuid "github.com/satori/go.uuid"
	"ko3-gin/pkg/model/common"
	"time"
)

type Credential struct {
	common.BaseModel
	Username   string
	Password   string
	PrivateKey string
	Type       string
}

func (c *Credential) BeforeCreate() error {
	c.ID = uuid.NewV4().String()
	c.CreatedDate = time.Now()
	c.UpdatedDate = time.Now()
	return nil
}

func (c *Credential) BeforeUpdate() error {
	c.UpdatedDate = time.Now()
	return nil
}

func (c Credential) TableName() string {
	return "ko_credential"
}
