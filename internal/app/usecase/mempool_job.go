package usecase

import (
	"context"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
	"github.com/LiquidCats/watcher/v2/internal/app/port/runner"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
)

type MempoolJob struct {
	cfg       configs.ChainConfig
	rpcClient rpc.Client
	txIDCh    runner.ChanWrite[entities.TxID]

	oldMempool []entities.TxID
}

func NewMempoolJob(
	cfg configs.ChainConfig,
	rpcClient rpc.Client,
	txIDCh runner.ChanWrite[entities.TxID],
	oldMempool []entities.TxID,
) *MempoolJob {
	return &MempoolJob{
		cfg:        cfg,
		rpcClient:  rpcClient,
		txIDCh:     txIDCh,
		oldMempool: oldMempool,
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

	newMempool, err := uc.rpcClient.GetMempool(ctx)
	if err != nil {
		return eris.Wrap(err, "get new mempool")
	}

	if len(newMempool) == 0 {
		return nil
	}

	m := make(map[entities.TxID]struct{}, len(newMempool))

	for _, txID := range uc.oldMempool {
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
		uc.txIDCh <- txID
	}

	uc.oldMempool = newMempool

	return nil
}
