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
	phaseName  = "migrate"
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
	return nil
}

func (i *InitMigrateDBPhase) PhaseName() string {
	return phaseName
}
