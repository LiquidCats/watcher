package configs

import (
	"fmt"
	"time"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type ChainsConfig []ChainConfig

type ChainConfig struct {
	Driver entities.Driver `yaml:"driver"`
	Type   entities.Type   `yaml:"type"`
	Chain  entities.Chain  `yaml:"chain"`

	Persist PersistConfig `yaml:"persist"`
	Scan    ScanConfig    `yaml:"scan"`

	Workers WorkersConfig `yaml:"workers"`

	RPC RPCConfig `yaml:"rpc"`

	Topics TopicsConfig `yaml:"topics"`
}

type TopicsConfig struct {
	Transactions string `yaml:"transactions"`
	Blocks       string `yaml:"blocks"`
}

type RPCConfig struct {
	NodeURL string `yaml:"node_url"`
}

type ScanConfig struct {
	Depth    int           `yaml:"depth"`
	Interval time.Duration `yaml:"interval"`
}

type PersistConfig struct {
	Capacity int           `yaml:"capacity"`
	Interval time.Duration `yaml:"interval"`
}

type WorkersConfig struct {
	TxIDWorkerCount              uint `default:"3" yaml:"txid_worker_count"`
	BlockTransactionsWorkerCount uint `default:"5" yaml:"block_transactions_worker_count"`
}

func (app ChainConfig) Key(k string) string {
	return fmt.Sprintf("%s.%s.%s.%s", app.Type, app.Driver, app.Chain, k)
}
