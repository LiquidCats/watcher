package database

import (
	"embed"

	"github.com/go-faster/errors"
	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(conn *pgx.Conn) error {
	// Create a new source driver using the embedded migrations.
	// Note: The second argument must match the directory name within the embedded FS.
	sourceDriver, err := iofs.New(migrations, "sql")
	if err != nil {
		return err
	}

	dbConn := stdlib.OpenDB(*conn.Config())

	// Create a new pgx migration driver instance.
	dbDriver, err := pgxmigrate.WithInstance(dbConn, &pgxmigrate.Config{})
	if err != nil {
		return err
	}

	// Create the migrate instance using the source and database drivers.
	m, err := migrate.NewWithInstance(
		"iofs", sourceDriver,
		"pgx", dbDriver,
	)
	if err != nil {
		return err
	}

	// Run the up migrations.
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
