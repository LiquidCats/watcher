package usecase

import (
	"context"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/bus"
	"github.com/LiquidCats/watcher/v2/internal/app/port/runner"
	"github.com/LiquidCats/watcher/v2/internal/app/port/state"

	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
)

type BlockHandler struct {
	cfg           configs.ChainConfig
	blockPub      bus.Publisher[entities.Block]
	state         state.State[entities.BlockHash]
	transactionCh runner.ChanWrite[entities.Transaction]
}

func NewBlockHandler(
	cfg configs.ChainConfig,
	blockPub bus.Publisher[entities.Block],
	state state.State[entities.BlockHash],
	transactionCh runner.ChanWrite[entities.Transaction],
) *BlockHandler {
	return &BlockHandler{
		cfg:           cfg,
		blockPub:      blockPub,
		state:         state,
		transactionCh: transactionCh,
	}
}

func (uc *BlockHandler) Handle(ctx context.Context, block entities.Block) error {
	logger := zerolog.Ctx(ctx).With().Any("block_hash", block.GetHash()).Logger()

	for _, tx := range block.GetTransactions() {
		uc.transactionCh <- tx
	}

	err := uc.blockPub.PublishTo(ctx, uc.cfg.Topics.Blocks, block)
	if err != nil {
		return eris.Wrap(err, "publish block")
	}

	blocksState, err := uc.state.Get(ctx, uc.cfg.Key("blocks"))
	if err != nil {
		return eris.Wrap(err, "get state")
	}

	if len(blocksState) >= uc.cfg.Persist.Capacity {
		blocksState = append(blocksState[1:], block.GetHash())
	} else {
		blocksState = append(blocksState, block.GetHash())
	}

	if err = uc.state.Set(
		ctx,
		uc.cfg.Key("blocks"),
		blocksState,
		uc.cfg.Persist.Duration,
	); err != nil {
		logger.Error().Err(err).Msg("set state")
	}

	logger.Info().Msg("published block")

	return nil
}
