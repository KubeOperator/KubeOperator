package db

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/jinzhu/gorm"
)

type AtomicFunc func(db *gorm.DB) error

func Tx(afs ...AtomicFunc) (err error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("%v", err)
		}
	}()
	for _, f := range afs {
		err = f(tx)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Commit().Error
	return
}
