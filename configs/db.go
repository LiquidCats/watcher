package configs

import (
	"fmt"

	"github.com/LiquidCats/watcher/v2/pkg/docker"
	"github.com/rs/zerolog"
)

type DB struct {
	Driver   string `yaml:"driver" envconfig:"DRIVER" default:"postgres"`
	Host     string `yaml:"host" envconfig:"HOST"`
	Port     string `yaml:"port" envconfig:"PORT"`
	Database string `yaml:"database" envconfig:"DATABASE"`
	User     string `yaml:"user" envconfig:"USER"`
	Password string `yaml:"password" envconfig:"PASSWORD"`
}

func (d *DB) ToDSN(logger *zerolog.Logger) string {
	pwd, err := docker.GetSecret(d.Password)
	logger.Fatal().Err(err).Msg("cant get db password from file")

	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		d.Driver,
		d.User,
		pwd,
		d.Host,
		d.Port,
		d.Database,
	)
}
