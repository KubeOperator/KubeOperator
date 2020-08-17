package model

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

var (
	AdminCanNotDelete = "ADMIN_CAN_NOT_DELETE"
	AdminCanNotUpdate = "ADMIN_CAN_NOT_UPDATE"
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
	Email    string `json:"email" gorm:"type:varchar(256);not null;unique"`
	IsActive bool   `json:"isActive" gorm:"type:boolean;default:true"`
	Language string `json:"language" gorm:"type:varchar(64)"`
	IsAdmin  bool   `json:"isAdmin" gorm:"type:boolean;default:false"`
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
	return err
}

func (u *User) BeforeUpdate() (err error) {
	return err
}

func (u User) TableName() string {
	return "ko_user"
}



func (u *User) ValidateOldPassword(password string) (bool, error) {
	oldPassword, err := encrypt.StringDecrypt(u.Password)
	if err != nil {
		return false, err
	}
	if oldPassword != password {
		return false, err
	}
	return true, err
}
