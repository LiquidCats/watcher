package port

import (
	"context"
	"watcher/internal/app/domain/entity"
)

type RpcRepository interface {
	GetLatestBlock(ctx context.Context) (*entity.Block, error)
	GetBlockByHash(ctx context.Context, hash entity.BlockHash) (*entity.Block, error)
	GetBlockByHeight(ctx context.Context, height entity.BlockHeight) (*entity.Block, error)
	GetBlockTransactions(ctx context.Context, hash entity.BlockHash) ([]entity.Transaction, error)
}
