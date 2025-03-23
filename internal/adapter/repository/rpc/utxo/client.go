package utxo

import (
	"context"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/pkg/jsonrpc"
	"github.com/pkg/errors"
)

type Client struct {
	cfg configs.UtxoRpc
}

func NewClient(cfg configs.UtxoRpc) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) GetMempool(ctx context.Context) ([]entities.TxID, error) {
	req, err := jsonrpc.Prepare[any](ctx, c.cfg.URL, "getrawmempool", nil)
	if err != nil {
		return nil, errors.Wrap(err, "GetMempool: prepare")
	}

	resp, err := jsonrpc.Execute[[]entities.TxID](req)
	if err != nil {
		return nil, errors.Wrap(err, "GetMempool: execute")
	}

	return *resp, nil
}

func (c *Client) GetLatestBlockHash(ctx context.Context) (entities.BlockHash, error) {
	req, err := jsonrpc.Prepare[any](ctx, c.cfg.URL, "getbestblockhash", nil)
	if err != nil {
		return "", errors.Wrap(err, "GetLatestBlockHash: prepare")
	}

	resp, err := jsonrpc.Execute[entities.BlockHash](req)
	if err != nil {
		return "", errors.Wrap(err, "GetLatestBlockHash: execute")
	}

	return *resp, nil
}

func (c *Client) GetBlockByHash(ctx context.Context, hash entities.BlockHash) (entities.Block, error) {
	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.URL, "getblock", []any{hash, 2})
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: prepare")
	}

	resp, err := jsonrpc.Execute[data.Block](req)
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: execute")
	}

	return resp, nil
}

func (c *Client) GetTransactionByTxId(ctx context.Context, hash entities.TxID) (entities.Transaction, error) {
	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.URL, "getrawtransaction", []any{hash, 2})
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: prepare")
	}

	resp, err := jsonrpc.Execute[data.Transaction](req)
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: execute")
	}

	return resp, nil
}
