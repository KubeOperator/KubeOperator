package migrate

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/file"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	phaseName = "migrate"
)

var migrationDirs = []string{
	"/usr/local/lib/migration",
	"./migration",
}

var log = logger.Default

type InitMigrateDBPhase struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

func (i *InitMigrateDBPhase) Init() error {
	url := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?charset=utf8",
		i.User,
		i.Password,
		i.Host,
		i.Port,
		i.Name)
	var path string
	for _, d := range migrationDirs {
		if file.Exists(d) {
			path = d
		}
	}
	if path == "" {
		return errors.New("can not find migration in ['/usr/local/lib/migration','./migration']")
	}
	filePath := fmt.Sprintf("file://%s", path)
	m, err := migrate.New(
		filePath, url)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no databases change,skip migrate")
			return nil
		}
		return err
	}
	return nil
}

func (i *InitMigrateDBPhase) PhaseName() string {
	return phaseName
}
