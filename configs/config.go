package configs

import (
	"fmt"
	"time"
	"watcher/internal/app/domain/entity"
)

type Config struct {
	Blockchain entity.Blockchain
	NodeUrl    string
	Interval   time.Duration
	Gap        entity.BlockHeight
	Workers    int

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
		"postgres://%s:%s@%s:%s/%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
	)
}
