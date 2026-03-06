package database

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type StateDB interface {
	SetBlockState(
		ctx context.Context,
		chain entities.Chain,
		blocks []entities.BlockHash,
	) error
}
