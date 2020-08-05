package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type BackupAccount struct {
	common.BaseModel
	ID         string `json:"_"`
	Name       string `json:"name" gorm:"type:varchar(256)"`
	Bucket     string `json:"bucket" gorm:"type:varchar(256)"`
	Credential string `json:"credential" gorm:"type:text(65535)"`
	Type       string `json:"type" gorm:"type:varchar(64)"`
	Status     string `json:"status" gorm:"type:varchar(64)"`
}

func (b *BackupAccount) BeforeCreate() (err error) {
	b.ID = uuid.NewV4().String()
	return err
}

func (b BackupAccount) TableName() string {
	return "ko_backup_account"
}
