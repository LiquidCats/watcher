package usecase

import (
	"context"
	"fmt"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/bus"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
	"github.com/LiquidCats/watcher/v2/internal/app/port/state"
	"github.com/go-faster/errors"
	"github.com/rs/zerolog"
)

type WatchMempoolUseCase struct {
	cfg                  configs.App
	state                state.State[[]entities.TxID]
	rpcClient            rpc.UtxoClient
	transactionPublisher bus.TransactionPublisher
}

func NewWatchMempoolUseCase(
	cfg configs.App,
	state state.State[[]entities.TxID],
	rpcClient rpc.UtxoClient,
	transactionPublisher bus.TransactionPublisher,
) *WatchMempoolUseCase {
	return &WatchMempoolUseCase{
		cfg:                  cfg,
		state:                state,
		rpcClient:            rpcClient,
		transactionPublisher: transactionPublisher,
	}
}

func (uc *WatchMempoolUseCase) Execute(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	oldMempool, err := uc.state.Get(ctx, uc.getStateKey())
	if err != nil {
		return errors.Wrap(err, "get old mempool")
	}

	newMempool, err := uc.rpcClient.GetMempool(ctx)
	if err != nil {
		return errors.Wrap(err, "get new mempool")
	}

	m := make(map[entities.TxID]struct{}, len(newMempool))

	for _, txID := range oldMempool {
		m[txID] = struct{}{}
	}

	diff := make([]entities.TxID, 0, len(m))

	for _, txID := range newMempool {
		_, ok := m[txID]
		if !ok {
			diff = append(diff, txID)
		}
	}

	for _, txID := range diff {
		tx, err := uc.rpcClient.GetTransactionByTxId(ctx, txID)
		if err != nil {
			logger.Error().
				Ctx(ctx).
				Err(err).
				Any("txID", txID).
				Msg("get transaction")
			continue
		}

		if err := uc.transactionPublisher.PublishTransaction(ctx, tx.ToEntity()); err != nil {
			logger.Error().
				Ctx(ctx).
				Err(err).
				Any("txID", txID).
				Msg("publish transaction")
			continue
		}
	}

	if err := uc.state.Set(
		ctx,
		uc.getStateKey(),
		newMempool,
		uc.cfg.PersistDuration,
	); err != nil {
		return errors.Wrap(err, "set new mempool")
	}

	return nil
}

func (uc *WatchMempoolUseCase) getStateKey() string {
	return fmt.Sprint(uc.cfg.Driver, ".", uc.cfg.Type, ".", uc.cfg.Chain, ".", "mempool")
}
