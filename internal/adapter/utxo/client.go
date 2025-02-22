package utxo

import (
	"context"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/entities"
	"github.com/LiquidCats/watcher/v2/pkg/jsonrpc"
	"github.com/go-faster/errors"
)

type Client[T any] struct {
	cfg configs.Utxo
}

func NewClient[T any](cfg configs.Utxo) *Client[T] {
	return &Client[T]{
		cfg: cfg,
	}
}

func (c *Client[T]) GetMempool(ctx context.Context) ([]entities.TxID, error) {
	req, err := jsonrpc.Prepare[any](ctx, c.cfg.NodeUrl, "getrawmempool", nil)
	if err != nil {
		return nil, errors.Wrap(err, "GetMempool: prepare")
	}

	resp, err := jsonrpc.Execute[[]entities.TxID](req)
	if err != nil {
		return nil, errors.Wrap(err, "GetMempool: execute")
	}

	return *resp, nil
}

func (c *Client[T]) GetBlockByHash(ctx context.Context, hash entities.BlockHash) (*data.Block[T], error) {
	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeUrl, "getblockbyhash", []any{hash, 1})
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: prepare")
	}

	resp, err := jsonrpc.Execute[data.Block[T]](req)
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: execute")
	}

	return resp, nil
}

func (c *Client[T]) GetTransactionByTxId(ctx context.Context, hash entities.TxID) (*data.Transaction, error) {
	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeUrl, "getrawtransaction", []any{hash, 1})
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: prepare")
	}

	resp, err := jsonrpc.Execute[data.Transaction](req)
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: execute")
	}

	return resp, nil
}
