package evm

import (
	"context"

	"github.com/LiquidCats/jsonrpc"
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/evm/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
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

func (c *Client) GetBlockByHash(ctx context.Context, hash entities.BlockHash) (*entities.Block, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Any("block_hash", hash).Msg("DEBUG")

	req, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeURL, "eth_getBlockByHash", []any{string(hash), false})
	if err != nil {
		return nil, eris.Wrapf(err, "prepare get block by hash %s", hash)
	}

	req.Header.Set("Content-Type", "application/json")

	block, err := jsonrpc.Execute[data.Block[entities.TxID]](req)
	if err != nil {
		return nil, eris.Wrapf(err, "execute get block by hash %s", hash)
	}

	return &entities.Block{
		Height:       entities.BlockHeight(block.Number.Uint64()),
		Hash:         block.Hash,
		Transactions: block.Transactions,
	}, nil
}

func (c *Client) GetBlockByHashWithTransactions(ctx context.Context, hash entities.BlockHash) (*entities.BlockWithTransactions[entities.TransactionAccountInput], error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Any("block_hash", hash).Msg("DEBUG")

	blockReq, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeURL, "eth_getBlockByHash", []any{string(hash), true})
	if err != nil {
		return nil, eris.Wrapf(err, "prepare get block by hash %s", hash)
	}

	blockReq.Header.Set("Content-Type", "application/json")

	block, err := jsonrpc.Execute[data.Block[data.Transaction]](blockReq)
	if err != nil {
		return nil, eris.Wrapf(err, "execute get block by hash %s", hash)
	}

	receiptsReq, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeURL, "eth_getBlockReceipts", []any{block.Number.String(), true})
	if err != nil {
		return nil, eris.Wrapf(err, "prepare get block by hash %s", hash)
	}

	receiptsReq.Header.Set("Content-Type", "application/json")

	receipts, err := jsonrpc.Execute[[]data.TransactionReceipt](receiptsReq)
	if err != nil {
		return nil, eris.Wrapf(err, "execute get block by hash %s", hash)
	}

	receiptsMap := make(map[entities.TxID]data.TransactionReceipt, len(*receipts))

	for _, receipt := range *receipts {
		receiptsMap[receipt.TransactionHash] = receipt
	}

	blockEntity := entities.BlockWithTransactions[entities.TransactionAccountInput]{
		Height:       entities.BlockHeight(block.Number.Uint64()),
		Hash:         block.Hash,
		Transactions: make([]entities.Transaction[entities.TransactionAccountInput], len(block.Transactions)),
	}

	for idx, tx := range block.Transactions {
		blockEntity.Transactions[idx] = entities.Transaction[entities.TransactionAccountInput]{
			TxID: tx.Hash,
			Inputs: []entities.TransactionAccountInput{
				{
					Address: tx.From,
				},
			},
			Outputs: []entities.TransactionOutput{},
			Fee:     decimal.RequireFromString("0"),
		}
	}

	return &blockEntity, nil
}

func (c *Client) GetTransactionByTxID(ctx context.Context, hash entities.TxID) (*entities.Transaction[entities.TransactionAccountInput], error) {
	txReq, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeURL, "eth_getTransactionByHash", []any{hash})
	if err != nil {
		return nil, eris.Wrapf(err, "prepare get block by hash %s", hash)
	}

	txReq.Header.Set("Content-Type", "application/json")

	tx, err := jsonrpc.Execute[data.Transaction](txReq)
	if err != nil {
		return nil, eris.Wrapf(err, "execute get block by hash %s", hash)
	}

	receiptReq, err := jsonrpc.Prepare[[]any](ctx, c.cfg.NodeURL, "eth_getTransactionReceipt", []any{hash})
	if err != nil {
		return nil, eris.Wrapf(err, "prepare get block by hash %s", hash)
	}

	receiptReq.Header.Set("Content-Type", "application/json")

	txReceipt, err := jsonrpc.Execute[data.TransactionReceipt](receiptReq)
	if err != nil {
		return nil, eris.Wrapf(err, "execute get block by hash %s", hash)
	}

	txEntity := entities.Transaction[entities.TransactionAccountInput]{
		TxID: tx.Hash,
		Inputs: []entities.TransactionAccountInput{
			{
				Address: tx.From,
			},
		},
		Outputs: make([]entities.TransactionOutput, 0, len(txReceipt.Logs)),
		Fee:     decimal.Decimal{},
	}

	for _, log := range txReceipt.Logs {

		txEntity.Outputs = append(txEntity.Outputs, entities.TransactionOutput{
			Address: log.Address,
			Value:   decimal.RequireFromString(log.Data),
		})
	}

	return &txEntity, nil
}

func (c *Client) GetMempool(_ context.Context) ([]entities.TxID, error) {
	return nil, nil
}
