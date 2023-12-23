package port

import (
	"context"
	"watcher/internal/app/domain/entity"
)

type StorageRepository interface {
	GetHeight(ctx context.Context, driver entity.Blockchain) (entity.BlockHeight, error)
	StoreBlock(ctx context.Context, block *entity.Block) error
	UpdateBlock(ctx context.Context, newBlock *entity.Block) error
	GetBlock(ctx context.Context, height entity.BlockHeight) (*entity.Block, error)
	ConfirmBlock(ctx context.Context, height entity.BlockHeight) error
	GetAllUnconfirmedBlocks(ctx context.Context, fromHeight entity.BlockHeight) ([]*entity.Block, error)
}
