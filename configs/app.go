package configs

import "github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"

type App struct {
	Driver entities.Driver
	Type   entities.Type
	Chain  entities.Chain
}
