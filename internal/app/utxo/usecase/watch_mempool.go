package usecase

import (
	"bytes"
	"context"
	"fmt"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/port/bus"
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/port/database"
	"github.com/LiquidCats/watcher/v2/internal/app/utxo/port/rpc"
	"github.com/bytedance/sonic"
	"github.com/go-faster/errors"
	"github.com/rs/zerolog"
)

type WatchMempoolUseCase struct {
	cfg                  configs.App
	stateDB              database.StateDB
	rpcClient            rpc.Client
	transactionPublisher bus.TransactionPublisher
	oldMempool           []entities.TxID
}

func NewWatchMempoolUseCase(
	cfg configs.App,
	stateDB database.StateDB,
	transactionPublisher bus.TransactionPublisher,
	rpcClient rpc.Client,
) *WatchMempoolUseCase {
	return &WatchMempoolUseCase{
		cfg:                  cfg,
		stateDB:              stateDB,
		rpcClient:            rpcClient,
		transactionPublisher: transactionPublisher,
		oldMempool:           []entities.TxID{},
	}
}

func (uc *WatchMempoolUseCase) Init(ctx context.Context) error {
	state, err := uc.stateDB.GetByKey(ctx, uc.getStateKey())
	if err != nil {
		return errors.Wrap(err, "get state")
	}
	reader := bytes.NewReader(state.Value)
	decoder := sonic.ConfigDefault.NewDecoder(reader)

	var oldMempool []entities.TxID

	if err := decoder.Decode(&oldMempool); err != nil {
		return errors.Wrap(err, "decode old mempool")
	}

	uc.oldMempool = oldMempool

	return nil
}

func (uc *WatchMempoolUseCase) Execute(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	newMempool, err := uc.rpcClient.GetMempool(ctx)
	if err != nil {
		return errors.Wrap(err, "get new mempool")
	}

	m := make(map[entities.TxID]struct{}, len(newMempool))

	for _, txID := range uc.oldMempool {
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

	uc.oldMempool = newMempool

	return nil
}

func (uc *WatchMempoolUseCase) getStateKey() string {
	return fmt.Sprint(uc.cfg.Driver, ".", uc.cfg.Type, ".", uc.cfg.Chain, ".", "mempool")
}
