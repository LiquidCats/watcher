package configs

import "fmt"

type DB struct {
	Driver   string `default:"postgres"`
	Host     string
	Port     string
	Database string
	User     string
	Password string
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
