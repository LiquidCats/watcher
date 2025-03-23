package usecase

import (
	"context"
	"fmt"
	"slices"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/bus"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
	"github.com/LiquidCats/watcher/v2/internal/app/port/state"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type BlocksProcessor struct {
	cfg                  configs.App
	state                state.State[entities.BlockHash]
	rpcClient            rpc.Client
	blockPublisher       bus.BlockPublisher
	transactionPublisher bus.TransactionPublisher
}

func NewBlocksProcessor(
	cfg configs.App,
	state state.State[entities.BlockHash],
	rpcClient rpc.Client,
	blockPublisher bus.BlockPublisher,
	transactionPublisher bus.TransactionPublisher,
) *BlocksProcessor {
	return &BlocksProcessor{
		cfg:                  cfg,
		state:                state,
		rpcClient:            rpcClient,
		blockPublisher:       blockPublisher,
		transactionPublisher: transactionPublisher,
	}
}

func (uc *BlocksProcessor) Execute(ctx context.Context) error {
	var blockHash entities.BlockHash
	var block entities.Block
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

	logger.Info().Any("block_hash", blockHash).Msg("starting form")
	if slices.Contains(blocksState, blockHash) {
		return nil
	}

	var blocks []entities.Block

	for {
		block, err = uc.rpcClient.GetBlockByHash(ctx, blockHash)
		if err != nil {
			return errors.Wrapf(err, "get block [%s]", blockHash)
		}

		blocks = append(blocks, block)

		if slices.Contains(blocksState, block.GetPrevHash()) {
			break
		}

		blockHash = block.GetPrevHash()

		if len(blocks) >= uc.cfg.ScanDepth {
			logger.Debug().Msg("scan block depth")
			break
		}
	}

	logger.Info().Any("blocks_len", len(blocks)).Msg("blocks collected")

	slices.Reverse(blocks)

	for _, block := range blocks {
		for _, tx := range block.GetTransactions() {
			if err = uc.transactionPublisher.PublishTransaction(ctx, tx); err != nil {
				logger.Error().Err(err).Any("txid", tx.GetTxID()).Msg("publish transaction")
			}
		}
		if err = uc.blockPublisher.PublishBlock(ctx, block); err != nil {
			logger.Error().Err(err).Any("hash", block.GetHash()).Msg("publish transaction")
		}

		logger.Info().Any("block_hash", block.GetHash()).Msgf("block published")

		if len(blocksState) >= uc.cfg.PersistBocks {
			blocksState = append(blocksState[1:], block.GetHash())
		} else {
			blocksState = append(blocksState, block.GetHash())
		}

		if err = uc.state.Set(
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

func (uc *BlocksProcessor) getStateKey() string {
	return fmt.Sprint(uc.cfg.Driver, ".", uc.cfg.Type, ".", uc.cfg.Chain, ".", "blocks")
}
