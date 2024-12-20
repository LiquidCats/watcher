package rpc

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/adapter/utxo/rpc/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/shared/entities"
	sharedEntities "github.com/LiquidCats/watcher/v2/internal/app/domain/utxo/entities"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetMempool(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (c *Client) GetBlockByHash(ctx context.Context, hash entities.BlockHash) (*data.Block, error) {
	return nil, nil
}

func (c *Client) GetTransactionByTxId(ctx context.Context, hash sharedEntities.TxHash) (*data.Transaction, error) {
	return nil, nil
}
