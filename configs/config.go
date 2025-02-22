package configs

import (
	"time"

	"github.com/LiquidCats/watcher/v2/internal/app/kernel/entities"
)

type Config struct {
	Driver   entities.Driver
	Interval time.Duration

	Evm  Evm
	Utxo Utxo

	DB DB
}
