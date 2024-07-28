package repository

import (
	"context"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entity"
)

type RpcCaller interface {
	CallFor(ctx context.Context, dest any, method string, params ...any) error
}

type RpcRepository interface {
	GetBlockByHash(ctx context.Context, blockHash entity.BlockHash) (*entity.Block, error)
	GetTransactionByTxId(ctx context.Context, txId entity.TxId) (*entity.Transaction, error)
}
