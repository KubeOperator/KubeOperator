package migration

import (
	"github.com/KubeOperator/KubeOperator/pkg/migration/migrations"
	"os"
)

type Migrate struct {
	Mysql  *migrations.Mysql
	Source *migrations.SourceDriver
}

func New(sourceUrl, databaseUrl string) (*Migrate, error) {

	mysql, err := migrations.NewMysql(databaseUrl)
	if err != nil {
		return nil, err
	}
	source, err := migrations.NewSourceDriver(sourceUrl)
	if err != nil {
		return nil, err
	}
	migrate := &Migrate{
		Mysql:  mysql,
		Source: source,
	}
	return migrate, nil
}

func (m *Migrate) Run() error {

	if err := m.Mysql.Init(); err != nil {
		return err
	}
	version, dirty, err := m.Mysql.Version()
	if err != nil {
		return err
	}
	if dirty {
		return migrations.ErrDatabaseDirty
	}
	mss, err := m.Source.ReadUp(version)
	if err != nil {
		return err
	}
	for _, v := range mss.Index {
		mi := mss.Migrations[v]
		if err := m.Exec(mi); err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrate) Version() (int, error) {
	version, _, err := m.Mysql.Version()
	if err != nil {
		return 0, err
	}
	return version, err
}

func (m *Migrate) Exec(migration migrations.Migration) error {

	f, err := os.Open(m.Source.Dir + "/" + migration.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := m.Mysql.SetVersion(migration.Version, true); err != nil {
		return err
	}
	if err := m.Mysql.Run(f); err != nil {
		return err
	}
	if err := m.Mysql.SetVersion(migration.Version, false); err != nil {
		return err
	}
	return nil
}

func (m *Migrate) Close() {
	_ = m.Mysql.Close()
}
