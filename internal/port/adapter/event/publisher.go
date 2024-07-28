package event

import (
	"context"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entity"
)

type EventPublisher interface {
	Block(ctx context.Context, blockchain entity.Blockchain, block *entity.Block) error
	Transaction(ctx context.Context, blockchain entity.Blockchain, block *entity.Transaction) error
}
