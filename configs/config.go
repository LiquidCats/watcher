package configs

import (
	"fmt"
	"watcher/internal/app/domain/entity"
)

type Config struct {
	Blockchain entity.Blockchain
	NodeUrl    string `envconfig:"NODE_URL"`
	Interval   int
	Gap        entity.BlockHeight

	DB          DB
	Broadcaster Broadcaster
}

type DB struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type Broadcaster struct {
	Host string
}

func (d *DB) ToPostgres() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
	)
}
