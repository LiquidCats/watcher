package utxo

import (
	"context"

	"github.com/LiquidCats/jsonrpc"
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"

	"github.com/rotisserie/eris"
)

type Client struct {
	cfg configs.RPCConfig
}

func NewClient(cfg configs.RPCConfig) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) GetMempool(ctx context.Context) ([]entities.TxID, error) {
	req, err := jsonrpc.Prepare[any](ctx, c.cfg.NodeURL, "getrawmempool", nil)
	if err != nil {
		return nil, eris.Wrap(err, "GetMempool: prepare")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := jsonrpc.Execute[[]entities.TxID](req)
	if err != nil {
		return nil, eris.Wrap(err, "GetMempool: execute")
	}

	return *resp, nil
}

func (c *Client) GetLatestBlockHash(ctx context.Context) (entities.BlockHash, error) {
	req, err := jsonrpc.Prepare[any](ctx, c.cfg.NodeURL, "getbestblockhash", nil)
	if err != nil {
		return "", eris.Wrap(err, "GetLatestBlockHash: prepare")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := jsonrpc.Execute[entities.BlockHash](req)
	if err != nil {
		return "", eris.Wrap(err, "GetLatestBlockHash: execute")
	}

	return *resp, nil
}

func (c *Client) GetBlockByHash(ctx context.Context, hash entities.BlockHash, withTx bool) (entities.Block, error) {
	var verbosity int
	if withTx {
		verbosity = 2
	} else {
		verbosity = 1
	}

	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeURL, "getblock", []any{hash, verbosity})
	if err != nil {
		return nil, eris.Wrap(err, "GetBlockByHash: prepare")
	}

	req.Header.Set("Content-Type", "application/json")

	var block entities.Block
	if withTx {
		block, err = jsonrpc.Execute[data.Block[*data.Transaction]](req)
	} else {
		block, err = jsonrpc.Execute[data.Block[entities.TxID]](req)
	}
	if err != nil {
		return nil, eris.Wrap(err, "GetBlockByHash: execute")
	}

	return block, nil
}

func (c *Client) GetTransactionByTxID(ctx context.Context, hash entities.TxID) (entities.Transaction, error) {
	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeURL, "getrawtransaction", []any{hash, 2})
	if err != nil {
		return nil, eris.Wrap(err, "GetBlockByHash: prepare")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := jsonrpc.Execute[data.Transaction](req)
	if err != nil {
		return nil, eris.Wrap(err, "GetBlockByHash: execute")
	}

	return resp, nil
}
