package model

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/jinzhu/gorm"
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

func (u *User) BeforeCreate() (err error) {
	u.ID = uuid.NewV4().String()
	return err
}

func (u *User) BeforeDelete() (err error) {
	if u.Name == "admin" {
		return errors.New(AdminCanNotDelete)
	}
	var member ProjectMember
	err = db.DB.Where(ProjectMember{UserID: u.ID}).Find(&member).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		} else {
			return err
		}
	}
	err = db.DB.Delete(&member).Error
	if err != nil {
		return err
	}
	var userMessage UserMessage
	err = db.DB.Where(UserMessage{UserID: u.ID}).Find(&userMessage).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		} else {
			return err
		}
	}
	err = db.DB.Delete(&userMessage).Error
	if err != nil {
		return err
	}
	var config UserNotificationConfig
	err = db.DB.Where(UserNotificationConfig{UserID: u.ID}).Find(&config).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		} else {
			return err
		}
	}
	err = db.DB.Delete(&config).Error
	if err != nil {
		return err
	}
	var receiver UserReceiver
	err = db.DB.Where(UserReceiver{UserID: u.ID}).Find(&receiver).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		} else {
			return err
		}
	}
	err = db.DB.Delete(&receiver).Error
	if err != nil {
		return err
	}
	return err
}

func (u *User) BeforeUpdate() (err error) {
	if u.Type == constant.Ldap {
		return errors.New(LdapCanNotUpdate)
	}
	return err
}
