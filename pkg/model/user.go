package model

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

var (
	AdminCanNotDelete = "ADMIN_CAN_NOT_DELETE"
	LdapCanNotUpdate  = "LDAP_CAN_NOT_UPDATE"
)

const (
	EN string = "en-US"
	ZH string = "zh-CN"
)

type User struct {
	common.BaseModel
	ID       string `json:"id" gorm:"type:varchar(64)"`
	Name     string `json:"name" gorm:"type:varchar(256);not null;unique"`
	Password string `json:"password" gorm:"type:varchar(256)"`
	Email    string `json:"email" gorm:"type:varchar(256)"`
	IsActive bool   `json:"isActive" gorm:"type:boolean;default:true"`
	Language string `json:"language" gorm:"type:varchar(64)"`
	IsAdmin  bool   `json:"isAdmin" gorm:"type:boolean;default:false"`
	IsSystem bool   `json:"isSystem" gorm:"type:boolean;default:false"`
	Type     string `json:"type" gorm:"type:varchar(64)"`
	IsFirst  bool   `json:"isFirst" gorm:"type:boolean;default:true"`
	ErrCount int    `json:"errCount" gorm:"type:int(64)"`
}

type Token struct {
	Token string `json:"access_token"`
}

func (u *User) BeforeCreate() error {
	u.ID = uuid.NewV4().String()
	return nil
}

func (u *User) BeforeDelete() (err error) {
	if u.Name == "admin" {
		return errors.New(AdminCanNotDelete)
	}
	if err := db.DB.Where("user_id = ?", u.ID).Delete(&ProjectMember{}).Error; err != nil {
		return err
	}
	if err := db.DB.Where("user_id = ?", u.ID).Delete(&UserMessage{}).Error; err != nil {
		return err
	}
	if err := db.DB.Where("user_id = ?", u.ID).Delete(&UserNotificationConfig{}).Error; err != nil {
		return err
	}
	if err := db.DB.Where("user_id = ?", u.ID).Delete(&UserReceiver{}).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeUpdate() error {
	if u.Type == constant.Ldap {
		return errors.New(LdapCanNotUpdate)
	}
	return nil
}
