package port

import (
	"context"
	"watcher/internal/app/domain/entity"
)

type StorageRepository interface {
	Transaction(ctx context.Context, cb func(ctx context.Context) error) error
	GetHeight(ctx context.Context, blockchain entity.Blockchain) (entity.BlockHeight, error)
	StoreBlock(ctx context.Context, blockchain entity.Blockchain, block *entity.Block) error
	UpdateBlock(ctx context.Context, blockchain entity.Blockchain, newBlock *entity.Block) error
	GetBlock(ctx context.Context, blockchain entity.Blockchain, height entity.BlockHeight) (*entity.Block, error)
	ConfirmBlock(ctx context.Context, blockchain entity.Blockchain, height entity.BlockHeight) error
	GetAllUnconfirmedBlocks(ctx context.Context, blockchain entity.Blockchain, fromHeight entity.BlockHeight) ([]*entity.Block, error)
}
