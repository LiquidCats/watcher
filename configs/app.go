package configs

import (
	"fmt"
	"time"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type App struct {
	Driver entities.Driver `yaml:"driver" envconfig:"DRIVER"`
	Type   entities.Type   `yaml:"type" envconfig:"TYPE"`
	Chain  entities.Chain  `yaml:"chain" envconfig:"CHAIN"`

	TxIDWorkerCount        uint `default:"5" yaml:"txid_worker" envconfig:"TXID_WORKER"`
	BlockWorkerCount       uint `default:"3" yaml:"block_worker" envconfig:"BLOCK_WORKER"`
	TransactionWorkerCount uint `default:"1" yaml:"transaction_worker" envconfig:"TRANSACTION_WORKER"`

	ScanDepth    int           `yaml:"scan_depth" envconfig:"SCAN_DEPTH"`
	ScanInterval time.Duration `yaml:"scan_interval" envconfig:"SCAN_INTERVAL"`

	PersistBocks    int           `yaml:"persist_bocks" envconfig:"PERSIST_BOCKS"`
	PersistDuration time.Duration `yaml:"persist_duration" envconfig:"PERSIST_DURATION"`
}

func (app App) Key(k string) string {
	return fmt.Sprintf("%s.%s.%s.%s", app.Type, app.Driver, app.Chain, k)
}
