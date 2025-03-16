package usecase

import (
	"context"
	"fmt"
	"slices"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/bus"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
	"github.com/LiquidCats/watcher/v2/internal/app/port/state"
	"github.com/go-faster/errors"
	"github.com/rs/zerolog"
)

type WatchBlocksUseCase struct {
	cfg                  configs.App
	state                state.State[[]entities.BlockHash]
	rpcClient            rpc.UtxoClient
	blockPublisher       bus.BlockPublisher
	transactionPublisher bus.TransactionPublisher
}

func NewWatchBlocksUseCase(
	cfg configs.App,
	state state.State[[]entities.BlockHash],
	rpcClient rpc.UtxoClient,
	blockPublisher bus.BlockPublisher,
	transactionPublisher bus.TransactionPublisher,
) *WatchBlocksUseCase {
	return &WatchBlocksUseCase{
		cfg:                  cfg,
		state:                state,
		rpcClient:            rpcClient,
		blockPublisher:       blockPublisher,
		transactionPublisher: transactionPublisher,
	}
}

func (uc *WatchBlocksUseCase) Execute(ctx context.Context) error {
	var blockHash entities.BlockHash
	var block *data.Block
	var err error

	logger := zerolog.Ctx(ctx)
	blocksState, err := uc.state.Get(ctx, uc.getStateKey())
	if err != nil {
		return errors.Wrap(err, "get state")
	}

	blockHash, err = uc.rpcClient.GetLatestBlockHash(ctx)
	if err != nil {
		return errors.Wrap(err, "get latest block hash")
	}

	if slices.Contains(blocksState, blockHash) {
		return nil
	}

	var blocks []*entities.UtxoBlock

	for {
		block, err = uc.rpcClient.GetBlockByHash(ctx, blockHash)
		if err != nil {
			return errors.Wrap(err, "get latest block")
		}

		blocks = append(blocks, block.ToEntity())

		if slices.Contains(blocksState, block.PreviousBlockHash) {
			break
		}

		if len(blocks) >= uc.cfg.ScanDepth {
			logger.Debug().Msg("scan block depth")
			break
		}

		blockHash = block.PreviousBlockHash
	}

	slices.Reverse(blocks)

	for _, block := range blocks {
		for _, tx := range block.Transactions {
			if err := uc.transactionPublisher.PublishTransaction(ctx, tx); err != nil {
				logger.Error().Err(err).Any("txid", tx.TxID).Msg("publish transaction")
			}
		}
		if err := uc.blockPublisher.PublishBlock(ctx, block); err != nil {
			logger.Error().Err(err).Any("hash", block.Hash).Msg("publish transaction")
		}

		if len(blocksState) >= uc.cfg.PersistBocks {
			blocksState = append(blocksState[1:], block.Hash)
		} else {
			blocksState = append(blocksState, block.Hash)
		}

		if err := uc.state.Set(
			ctx,
			uc.getStateKey(),
			blocksState,
			uc.cfg.PersistDuration,
		); err != nil {
			logger.Error().Err(err).Msg("set state")
		}
	}

	return nil
}

func (uc *WatchBlocksUseCase) getStateKey() string {
	return fmt.Sprint(uc.cfg.Driver, ".", uc.cfg.Type, ".", uc.cfg.Chain, ".", "blocks")
}
