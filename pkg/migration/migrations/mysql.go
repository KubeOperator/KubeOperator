package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	ErrDatabaseDirty = fmt.Errorf("database is dirty")
	DefaultTable     = "schema_migrations"
)

type Config struct {
	MigrationsTable string
	DatabaseName    string
}

type Mysql struct {
	conn   *sql.Conn
	db     *sql.DB
	config *Config
}

func NewMysql(url string) (*Mysql, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	conn, err := db.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	my := &Mysql{
		db:   db,
		conn: conn,
		config: &Config{
			MigrationsTable: DefaultTable,
		},
	}
	return my, nil
}

func (m *Mysql) Close() error {
	connErr := m.conn.Close()
	dbErr := m.db.Close()
	if connErr != nil || dbErr != nil {
		return fmt.Errorf("conn: %v, db: %v", connErr, dbErr)
	}
	return nil
}

func (m *Mysql) Run(migration io.Reader) error {
	migr, err := ioutil.ReadAll(migration)
	if err != nil {
		return err
	}
	query := string(migr[:])
	if _, err := m.conn.ExecContext(context.Background(), query); err != nil {
		return fmt.Errorf("runErr: %v, migration failed: %v", err, query)
	}
	return nil
}

func (m *Mysql) Init() error {
	create := `CREATE TABLE IF NOT EXISTS schema_migrations (
        version bigint NOT NULL,
        dirty  tinyint(1) NOT NULL,
        PRIMARY KEY (version)  
    );`

	if _, err := m.db.Exec(create); err != nil {
		fmt.Println("Unable to create `schema_migrations` table", err)
		return err
	}
	return nil
}

func (m *Mysql) Version() (version int, dirty bool, err error) {
	query := "SELECT version, dirty FROM `" + m.config.MigrationsTable + "` LIMIT 1"
	err = m.conn.QueryRowContext(context.Background(), query).Scan(&version, &dirty)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return 0, false, nil
		} else {
			return 0, false, err
		}
	}
	return
}

func (m *Mysql) SetVersion(version int, dirty bool) error {
	tx, err := m.conn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	query := "TRUNCATE `" + m.config.MigrationsTable + "`"
	if _, err := tx.ExecContext(context.Background(), query); err != nil {
		return err
	}
	if version >= 0 || (version == -1 && dirty) {
		query := "INSERT INTO `" + m.config.MigrationsTable + "` (version, dirty) VALUES (?, ?)"
		if _, err := tx.ExecContext(context.Background(), query, version, dirty); err != nil {
			return err
		}
	}
	return tx.Commit()
}
