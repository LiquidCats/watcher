package utxo

import (
	"context"

	"github.com/LiquidCats/watcher/v2/configs"
	data2 "github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	entities2 "github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
	"github.com/LiquidCats/watcher/v2/pkg/jsonrpc"
	"github.com/go-faster/errors"
)

type Client struct {
	cfg configs.UtxoRpc
}

func NewClient(cfg configs.UtxoRpc) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) GetMempool(ctx context.Context) ([]entities2.TxID, error) {
	req, err := jsonrpc.Prepare[any](ctx, c.cfg.URL, "getrawmempool", nil)
	if err != nil {
		return nil, errors.Wrap(err, "GetMempool: prepare")
	}

	resp, err := jsonrpc.Execute[[]entities2.TxID](req)
	if err != nil {
		return nil, errors.Wrap(err, "GetMempool: execute")
	}

	return *resp, nil
}

func (c *Client) GetLatestBlockHash(ctx context.Context) (entities2.BlockHash, error) {
	req, err := jsonrpc.Prepare[any](ctx, c.cfg.URL, "getbestblockhash", nil)
	if err != nil {
		return "", errors.Wrap(err, "GetLatestBlockHash: prepare")
	}

	resp, err := jsonrpc.Execute[entities2.BlockHash](req)
	if err != nil {
		return "", errors.Wrap(err, "GetLatestBlockHash: execute")
	}

	return *resp, nil
}

func (c *Client) GetBlockByHash(ctx context.Context, hash entities2.BlockHash) (*data2.Block, error) {
	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.URL, "getblockbyhash", []any{hash, 2})
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: prepare")
	}

	resp, err := jsonrpc.Execute[data2.Block](req)
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: execute")
	}

	return resp, nil
}

func (c *Client) GetTransactionByTxId(ctx context.Context, hash entities2.TxID) (*data2.Transaction, error) {
	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.URL, "getrawtransaction", []any{hash, 1})
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: prepare")
	}

	resp, err := jsonrpc.Execute[data2.Transaction](req)
	if err != nil {
		return nil, errors.Wrap(err, "GetBlockByHash: execute")
	}

	return resp, nil
}
