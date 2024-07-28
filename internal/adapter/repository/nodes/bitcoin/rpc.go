package bitcoin

import (
	"context"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entity"
)

type RpcClient struct {
}

func NewRpcClient() *RpcClient {
	return &RpcClient{}
}

func (c *RpcClient) GetTransactionByTxId(ctx context.Context, txId entity.TxId) (*entity.Transaction, error) {
	return nil, nil
}

func (c *RpcClient) GetBlockByHash(ctx context.Context, blockHash entity.BlockHash) (*entity.Block, error) {
	return nil, nil
}
