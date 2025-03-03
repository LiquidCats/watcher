package database

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
)

type StateDB interface {
	CreateState(ctx context.Context, arg database.CreateStateParams) error
	GetByKey(ctx context.Context, key string) (database.State, error)
	UpdateState(ctx context.Context, arg database.UpdateStateParams) error
}
