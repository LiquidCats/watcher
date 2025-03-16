package rpc

import (
	"context"

	utxoData "github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type EvmClient interface {
	GetLatestBlockHash(ctx context.Context) (entities.BlockHash, error)
	GetBlockByHash(ctx context.Context, hash entities.BlockHash) (*any, error)
	GetTransactionByTxId(ctx context.Context, hash entities.TxID) (*any, error)
}

type UtxoClient interface {
	GetMempool(ctx context.Context) ([]entities.TxID, error)
	GetLatestBlockHash(ctx context.Context) (entities.BlockHash, error)
	GetBlockByHash(ctx context.Context, hash entities.BlockHash) (*utxoData.Block, error)
	GetTransactionByTxId(ctx context.Context, hash entities.TxID) (*utxoData.Transaction, error)
}
