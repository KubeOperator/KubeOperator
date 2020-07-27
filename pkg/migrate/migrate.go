package migrate

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	phaseName = "migrate"
)

type InitMigrateDBPhase struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

func (i *InitMigrateDBPhase) Init() error {
	var log = logger.Default
	for _, m := range model.Models {
		log.Infof("migrate table: %s", m.TableName())
		db.DB.AutoMigrate(m)
	}
	for _, d := range model.InitData {
		switch v := d.(type) {
		case model.User:
			op, ok := d.(model.User)
			if ok {
				user := model.User{}
				db.DB.Model(model.User{}).Where("name = ?", op.Name).First(&user)
				if db.DB.NewRecord(user) {
					db.DB.Create(d)
				} else {
					db.DB.Save(d)
				}
			}
		case model.Credential:
			op, ok := d.(model.Credential)
			if ok {
				credential := model.Credential{}
				db.DB.Model(model.Credential{}).Where("name = ?", op.Name).First(&credential)
				if db.DB.NewRecord(credential) {
					db.DB.Create(d)
				} else {
					db.DB.Save(d)
				}
			}

		default:
			log.Infof("insert data failed: %s", v)
		}
	}
	return nil
}

func (i *InitMigrateDBPhase) PhaseName() string {
	return phaseName
}
