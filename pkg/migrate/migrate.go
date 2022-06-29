package migrate

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/migration"

	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/file"
)

const (
	phaseName = "migrate"
)

const (
	releaseMigrationDir = "/usr/local/lib/ko/migration"
	localMigrationDir   = "./migration"
)

var migrationDirs = []string{
	localMigrationDir,
	releaseMigrationDir,
}

type InitMigrateDBPhase struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

func (i *InitMigrateDBPhase) Init() error {
	p, err := encrypt.StringDecrypt(i.Password)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Asia%%2FShanghai&multiStatements=true",
		i.User,
		p,
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
		return fmt.Errorf("can not find migration in [%s,%s]", localMigrationDir, releaseMigrationDir)
	}

	m, err := migration.New(path, url)
	if err != nil {
		return err
	}
	err = m.Run()
	if err != nil {
		return err
	}
	v, _ := m.Version()
	dp, err := encrypt.StringEncrypt(constant.DefaultPassword)
	if err != nil {
		return fmt.Errorf("can not init default user")
	}
	if !(v > 0) {
		if err := db.DB.Model(&model.User{}).Where("name = ?", "admin").Updates(map[string]interface{}{"Password": dp}).Error; err != nil {
			return fmt.Errorf("can not update default user")
		}
	}
	m.Close()
	return nil
}

func (i *InitMigrateDBPhase) PhaseName() string {
	return phaseName
}
