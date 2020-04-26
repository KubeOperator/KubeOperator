package migrate

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrateDir = "/Users/shenchenyang/go/src/ko3-gin/resource/migration"
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
	url := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?charset=utf8",
		i.User,
		i.Password,
		i.Host,
		i.Port,
		i.Name)
	filePath := fmt.Sprintf("file://%s", migrateDir)
	m, err := migrate.New(
		filePath, url)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}
	return nil
}

func (i *InitMigrateDBPhase) PhaseName() string {
	return phaseName
}
