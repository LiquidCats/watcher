package configs

import (
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
)

type Config struct {
	Driver entities.Driver

	Utxo Utxo

	DB DB `yaml:"db" envconfig:"DB"`
}
