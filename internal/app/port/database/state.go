package database

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/adapter/database"
)

type StateDB interface {
	GetStateByKey(ctx context.Context, key string) (database.State, error)
	SetState(ctx context.Context, arg database.SetStateParams) error
}
