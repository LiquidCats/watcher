package bus

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type TransactionPublisher interface {
	PublishTransaction(ctx context.Context, transaction entities.Transaction) error
}

type BlockPublisher interface {
	PublishBlock(ctx context.Context, block entities.Block) error
}
