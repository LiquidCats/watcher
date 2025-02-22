package ports

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/adapter/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/entities"
)

type NodeClient[T any] interface {
	GetMempool(ctx context.Context) ([]entities.TxID, error)
	GetBlockByHash(ctx context.Context, hash entities.BlockHash) (*data.Block[T], error)
	GetTransactionByTxId(ctx context.Context, hash entities.TxID) (*data.Transaction, error)
}
