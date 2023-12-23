package port

import (
	"context"
	"watcher/internal/app/domain/entity"
)

type EventPublisher interface {
	NewBlock(ctx context.Context, block *entity.Block) error
	ConfirmBlock(ctx context.Context, block *entity.Block) error
	RejectBlock(ctx context.Context, block *entity.Block) error
}
