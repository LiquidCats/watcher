package usecase

import (
	"context"
	"slices"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/bus"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
	"github.com/LiquidCats/watcher/v2/internal/app/port/runner"
	"github.com/LiquidCats/watcher/v2/internal/app/port/state"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
)

type BlocksJob struct {
	cfg       configs.ChainConfig
	state     state.State[entities.BlockHash]
	rpcClient rpc.Client
	workerCh  runner.ChanWrite[entities.Block]
	publisher bus.Publisher[entities.Block]
}

func NewBlocksJob(
	cfg configs.ChainConfig,
	state state.State[entities.BlockHash],
	rpcClient rpc.Client,
	workerCh runner.ChanWrite[entities.Block],
	publisher bus.Publisher[entities.Block],
) *BlocksJob {
	return &BlocksJob{
		cfg:       cfg,
		state:     state,
		rpcClient: rpcClient,
		workerCh:  workerCh,
		publisher: publisher,
	}
}

func (uc *BlocksJob) Handle(ctx context.Context) error {
	var blockHash entities.BlockHash
	var block entities.Block
	var err error

	logger := zerolog.Ctx(ctx).With().
		Str("name", "blocks_job").
		Any("driver", uc.cfg.Driver).
		Any("type", uc.cfg.Type).
		Any("chain", uc.cfg.Chain).
		Logger()

	blocksState, err := uc.state.Get(ctx, uc.cfg.Key("blocks"))
	if err != nil {
		return eris.Wrap(err, "get state")
	}

	blockHash, err = uc.rpcClient.GetLatestBlockHash(ctx)
	if err != nil {
		return eris.Wrap(err, "get latest block hash")
	}

	logger.Info().Any("block_hash", blockHash).Msg("starting form")
	if slices.Contains(blocksState, blockHash) {
		return nil
	}

	var blocks []entities.Block

	for {
		block, err = uc.rpcClient.GetBlockByHash(ctx, blockHash, false)
		if err != nil {
			return eris.Wrapf(err, "get block [%s]", blockHash)
		}

		blocks = append(blocks, block)

		exists := slices.Contains(blocksState, block.GetPrevHash())
		if exists {
			break
		}

		blockHash = block.GetPrevHash()

		if len(blocks) >= uc.cfg.Scan.Depth {
			logger.Debug().Msg("scan block depth")
			break
		}
	}

	logger.Info().Any("blocks_len", len(blocks)).Msg("blocks collected")

	slices.Reverse(blocks)

	for _, b := range blocks {
		err = uc.publisher.PublishTo(ctx, uc.cfg.Topics.Blocks, b)
		if err != nil {
			return eris.Wrap(err, "publish block")
		}

		if len(blocksState) >= uc.cfg.Persist.Capacity {
			blocksState = append(blocksState[1:], b.GetHash())
		} else {
			blocksState = append(blocksState, b.GetHash())
		}

		if err = uc.state.Set(
			ctx,
			uc.cfg.Key("blocks"),
			blocksState,
			uc.cfg.Persist.Duration,
		); err != nil {
			logger.Error().Err(err).Msg("set state")
		}

		uc.workerCh <- b
	}

	return nil
}
