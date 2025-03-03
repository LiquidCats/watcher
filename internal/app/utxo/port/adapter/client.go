package adapter

import (
	"context"

	data2 "github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	entities2 "github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
)

type NodeClient[T any] interface {
	GetMempool(ctx context.Context) ([]entities2.TxID, error)
	GetLatestBlockHash(ctx context.Context) (entities2.BlockHash, error)
	GetBlockByHash(ctx context.Context, hash entities2.BlockHash) (*data2.Block, error)
	GetTransactionByTxId(ctx context.Context, hash entities2.TxID) (*data2.Transaction, error)
}
