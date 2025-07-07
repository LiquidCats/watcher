package metrics

import "github.com/LiquidCats/watcher/v2/internal/app/domain/entities"

type RequestToNodeCounter interface {
	Inc(chain entities.Chain)
}
