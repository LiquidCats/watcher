package usecase

import (
	"context"
	"fmt"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/bus"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
	"github.com/LiquidCats/watcher/v2/internal/app/port/state"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type MempoolProcessor struct {
	cfg                  configs.App
	state                state.State[entities.TxID]
	rpcClient            rpc.UtxoClient
	transactionPublisher bus.TransactionPublisher
}

func NewMempoolProcessor(
	cfg configs.App,
	state state.State[entities.TxID],
	rpcClient rpc.UtxoClient,
	transactionPublisher bus.TransactionPublisher,
) *MempoolProcessor {
	return &MempoolProcessor{
		cfg:                  cfg,
		state:                state,
		rpcClient:            rpcClient,
		transactionPublisher: transactionPublisher,
	}
}

func (uc *MempoolProcessor) Execute(ctx context.Context) error {
	logger := zerolog.Ctx(ctx).
		With().
		Any("chain", uc.cfg.Chain).
		Any("driver", uc.cfg.Driver).
		Any("type", uc.cfg.Type).
		Str("module", "mempool_processor").
		Logger()

	oldMempool, err := uc.state.Get(ctx, uc.getStateKey())
	if err != nil {
		return errors.Wrap(err, "get old mempool")
	}

	logger.Debug().Any("old_mempool_len", len(oldMempool)).Msg("old mempool")

	newMempool, err := uc.rpcClient.GetMempool(ctx)
	if err != nil {
		return errors.Wrap(err, "get new mempool")
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
		var tx entities.Transaction

		tx, err = uc.rpcClient.GetTransactionByTxID(ctx, txID)
		if err != nil {
			logger.Error().
				Ctx(ctx).
				Err(err).
				Any("txID", txID).
				Msg("get transaction")
			continue
		}

		if err = uc.transactionPublisher.PublishTransaction(ctx, tx); err != nil {
			logger.Error().
				Ctx(ctx).
				Err(err).
				Any("txID", txID).
				Msg("publish transaction")
			continue
		}
	}

	logger.Info().Any("diff_len", len(diff)).Msg("transactions published")

	if err = uc.state.Set(
		ctx,
		uc.getStateKey(),
		newMempool,
		uc.cfg.PersistDuration,
	); err != nil {
		return errors.Wrap(err, "set new mempool")
	}

	return nil
}

func (uc *MempoolProcessor) getStateKey() string {
	return fmt.Sprint(uc.cfg.Driver, ".", uc.cfg.Type, ".", uc.cfg.Chain, ".", "mempool")
}
