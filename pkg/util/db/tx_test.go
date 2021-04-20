package db

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/jinzhu/gorm"
	"testing"
)

func TestTx(t *testing.T) {

	var afs []AtomicFunc

	afs = append(afs, func(db *gorm.DB) error {
		return db.Create(model.Host{}).Error
	})
	afs = append(afs, func(db *gorm.DB) error {
		return db.Create(model.Cluster{}).Error
	})

	if err := Tx(afs...); err != nil {

	}

}
