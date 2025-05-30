package evm

import (
	"context"

	"github.com/LiquidCats/jsonrpc"
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/evm/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
)

type Client struct {
	cfg configs.RPCConfig
}

func NewClient(cfg configs.RPCConfig) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) GetLatestBlockHash(ctx context.Context) (entities.BlockHash, error) {
	req, err := jsonrpc.Prepare[any](ctx, c.cfg.NodeURL, "eth_getBlockByNumber", []any{"latest", false})
	if err != nil {
		return "", eris.Wrap(err, "prepare get latest block hash")
	}

	req.Header.Set("Content-Type", "application/json")

	block, err := jsonrpc.Execute[data.LatestBlock](req)
	if err != nil {
		return "", eris.Wrap(err, "execute get latest block hash")
	}

	return block.Hash, nil
}

func (c *Client) GetBlockByHash(ctx context.Context, hash entities.BlockHash, withTx bool) (entities.Block, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Any("block_hash", hash).Msg("DEBUG")

	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeURL, "eth_getBlockByHash", []any{string(hash), withTx})
	if err != nil {
		return nil, eris.Wrapf(err, "prepare get block by hash %s", hash)
	}

	req.Header.Set("Content-Type", "application/json")

	var block entities.Block
	if withTx {
		block, err = jsonrpc.Execute[data.Block[*data.Transaction]](req)
	} else {
		block, err = jsonrpc.Execute[data.Block[entities.TxID]](req)
	}

	if err != nil {
		return nil, eris.Wrapf(err, "execute get block by hash %s", hash)
	}

	return block, nil
}

func (c *Client) GetTransactionByTxID(ctx context.Context, hash entities.TxID) (entities.Transaction, error) {
	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeURL, "eth_getTransactionByHash", []any{hash})
	if err != nil {
		return nil, eris.Wrapf(err, "prepare get block by hash %s", hash)
	}

	req.Header.Set("Content-Type", "application/json")

	tx, err := jsonrpc.Execute[data.Transaction](req)
	if err != nil {
		return nil, eris.Wrapf(err, "execute get block by hash %s", hash)
	}

	return tx, nil
}

func (c *Client) GetMempool(_ context.Context) ([]entities.TxID, error) {
	return nil, nil
}
