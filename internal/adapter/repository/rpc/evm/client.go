package evm

import (
	"context"

	"github.com/LiquidCats/jsonrpc"
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/evm/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/pkg/errors"
)

type Client struct {
	cfg configs.EvmRPC
}

func NewClient(cfg configs.EvmRPC) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c Client) GetLatestBlockHash(ctx context.Context) (entities.BlockHash, error) {
	req, err := jsonrpc.Prepare[any](ctx, c.cfg.URL, "eth_getBlockByNumber", []any{"latest", false})
	if err != nil {
		return "", errors.Wrap(err, "prepare get latest block hash")
	}

	block, err := jsonrpc.Execute[data.LatestBlock](req)
	if err != nil {
		return "", errors.Wrap(err, "execute get latest block hash")
	}

	return block.Hash, nil
}

func (c Client) GetBlockByHash(ctx context.Context, hash entities.BlockHash) (entities.Block, error) {
	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.URL, "eth_getBlockByHash", []any{hash, true})
	if err != nil {
		return nil, errors.Wrapf(err, "prepare get block by hash %s", hash)
	}

	block, err := jsonrpc.Execute[data.Block](req)
	if err != nil {
		return nil, errors.Wrapf(err, "execute get block by hash %s", hash)
	}

	return block, nil
}

func (c Client) GetTransactionByTxID(ctx context.Context, hash entities.TxID) (entities.Transaction, error) {
	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.URL, "eth_getTransactionByHash", []any{hash})
	if err != nil {
		return nil, errors.Wrapf(err, "prepare get block by hash %s", hash)
	}

	tx, err := jsonrpc.Execute[data.Transaction](req)
	if err != nil {
		return nil, errors.Wrapf(err, "execute get block by hash %s", hash)
	}

	return tx, nil
}
