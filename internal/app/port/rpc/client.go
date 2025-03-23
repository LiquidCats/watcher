package rpc

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type Client interface {
	GetLatestBlockHash(ctx context.Context) (entities.BlockHash, error)
	GetBlockByHash(ctx context.Context, hash entities.BlockHash) (entities.Block, error)
	GetTransactionByTxId(ctx context.Context, hash entities.TxID) (entities.Transaction, error)
}

type UtxoClient interface {
	GetMempool(ctx context.Context) ([]entities.TxID, error)
	Client
	//GetLatestBlockHash(ctx context.Context) (entities.BlockHash, error)
	//GetBlockByHash(ctx context.Context, hash entities.BlockHash) (*utxoData.Block, error)
	//GetTransactionByTxId(ctx context.Context, hash entities.TxID) (*utxoData.Transaction, error)
}
