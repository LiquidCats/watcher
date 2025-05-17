package usecase

import (
	"context"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
	"github.com/LiquidCats/watcher/v2/internal/app/port/runner"
	"github.com/LiquidCats/watcher/v2/internal/app/port/state"
	"github.com/go-faster/errors"
	"github.com/rs/zerolog"
)

type MempoolJob struct {
	cfg       configs.App
	state     state.State[entities.TxID]
	rpcClient rpc.Client
	txIdCh    runner.ChanWrite[entities.TxID]
}

func NewMempoolJob(
	cfg configs.App,
	state state.State[entities.TxID],
	rpcClient rpc.Client,
	txIdCh runner.ChanWrite[entities.TxID],
) *MempoolJob {
	return &MempoolJob{
		cfg:       cfg,
		state:     state,
		rpcClient: rpcClient,
		txIdCh:    txIdCh,
	}
}

func (uc *MempoolJob) Handle(ctx context.Context) error {
	logger := zerolog.Ctx(ctx).
		With().
		Any("chain", uc.cfg.Chain).
		Any("driver", uc.cfg.Driver).
		Any("type", uc.cfg.Type).
		Str("module", "mempool_processor").
		Logger()

	oldMempool, err := uc.state.Get(ctx, uc.cfg.Key("mempool"))
	if err != nil {
		return errors.Wrap(err, "get old mempool")
	}

	logger.Debug().Any("old_mempool_len", len(oldMempool)).Msg("old mempool")

	newMempool, err := uc.rpcClient.GetMempool(ctx)
	if err != nil {
		return errors.Wrap(err, "get new mempool")
	}

	if len(newMempool) == 0 && len(oldMempool) == 0 {
		return nil
	}

	m := make(map[entities.TxID]struct{}, len(newMempool))

	for _, txID := range oldMempool {
		m[txID] = struct{}{}
	}

	var diff []entities.TxID

	for _, txID := range newMempool {
		_, ok := m[txID]
		if !ok {
			diff = append(diff, txID)
		}
	}

	logger.Info().Any("diff_len", len(diff)).Msg("found new transactions")

	for _, txID := range diff {
		uc.txIdCh <- txID
	}

	if err = uc.state.Set(
		ctx,
		uc.cfg.Key("mempool"),
		newMempool,
		uc.cfg.PersistDuration,
	); err != nil {
		return errors.Wrap(err, "set new mempool")
	}

	return nil
}
