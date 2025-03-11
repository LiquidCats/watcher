package configs

import (
	"time"

	"github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
)

type App struct {
	Driver entities.Driver `yaml:"driver" envconfig:"DRIVER"`
	Type   entities.Type   `yaml:"type" envconfig:"TYPE"`
	Chain  entities.Chain  `yaml:"chain" envconfig:"CHAIN"`

	ScanDepth       int           `yaml:"scan_depth" envconfig:"SCAN_DEPTH"`
	PersistBocks    int           `yaml:"persist_bocks" envconfig:"PERSIST_BOCKS"`
	PersistDuration time.Duration `yaml:"persist_duration" envconfig:"PERSIST_DURATION"`
}
