package db

import (
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/escape"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

const phaseName = "db"

type InitDBPhase struct {
	Host         string
	Port         int
	Name         string
	User         string
	Password     string
	MaxOpenConns int
	MaxIdleConns int
}

func (i *InitDBPhase) Init() error {
	var err error
	i.Password, err = encrypt.StringDecryptWithSalt(i.Password)
	if err != nil {
		return err
	}
	url := []byte(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Asia%%2FShanghai", i.User, i.Password, i.Host, i.Port, i.Name))
	db, err := gorm.Open("mysql", string(url))
	if err != nil {
		return err
	}
	escape.Clean(url)
	gorm.DefaultTableNameHandler = func(DB *gorm.DB, defaultTableName string) string {
		return "ko_" + defaultTableName
	}
	db.SingularTable(true)
	db.DB().SetMaxOpenConns(i.MaxOpenConns)
	db.DB().SetMaxIdleConns(i.MaxIdleConns)
	DB = db
	DB.LogMode(false)
	return nil
}

func (i *InitDBPhase) PhaseName() string {
	return phaseName
}
