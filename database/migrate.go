package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"watcher/configs"
)

var ErrUnsupportedDriver = errors.New("unsupported driver")

func Migrate(conn *sql.DB, cfg configs.Config) error {
	dbInstance, err := getDriver(conn, cfg.DB.Driver)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("file://%s", "database/migrations")

	m, err := migrate.NewWithDatabaseInstance(path, cfg.DB.Database, dbInstance)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func getDriver(conn *sql.DB, driver string) (database.Driver, error) {
	switch driver {
	case "postgres":
		return migratePostgresDriver(conn)
	default:
		return nil, ErrUnsupportedDriver
	}
}

func migratePostgresDriver(conn *sql.DB) (database.Driver, error) {
	return postgres.WithInstance(conn, &postgres.Config{})
}
