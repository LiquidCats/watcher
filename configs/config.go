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
	Cleanup    Cleanup

	DB          DB
	Broadcaster Broadcaster
}

type Cleanup struct {
	Enabled  bool               `default:"false"`
	Interval int                `default:"120"`
	Gap      entity.BlockHeight `default:"10"` // how many block from the last confirmed
}

type DB struct {
	Driver   string `default:"postgres"`
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type Broadcaster struct {
	Host string
}

func (d *DB) ToDSN() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		d.Driver,
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
	)
}
