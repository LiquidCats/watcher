package configs

import (
	"time"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/shared/entities"
)

type Config struct {
	Drvier   entities.Driver
	Interval time.Duration

	Evm  Evm
	Utxo Utxo

	DB DB
}
