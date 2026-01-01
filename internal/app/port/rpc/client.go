package rpc

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type Client[TxIn any] interface {
	GetLatestBlockHash(ctx context.Context) (entities.BlockHash, error)
	GetBlockByHash(ctx context.Context, hash entities.BlockHash) (*entities.Block, error)
	GetBlockByHashWithTransactions(ctx context.Context, hash entities.BlockHash) (*entities.BlockWithTransactions[TxIn], error)
	GetTransactionByTxID(ctx context.Context, hash entities.TxID) (*entities.Transaction[TxIn], error)
	GetMempool(ctx context.Context) ([]entities.TxID, error)
}
