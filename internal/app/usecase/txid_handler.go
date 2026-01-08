package usecase

import (
	"context"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/bus"
	"github.com/LiquidCats/watcher/v2/internal/app/port/metrics"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
)

type TxIDHandler struct {
	cfg       configs.ChainConfig
	rpcClient rpc.Client[any]
	publisher bus.Publisher[entities.Transaction[any]]
	metrics   TxIDHandlerMetrics
}

type TxIDHandlerMetrics struct {
	RequestToNodeCounter metrics.RequestToNodeCounter
}

func NewTxIDHandler(
	cfg configs.ChainConfig,
	rpcClient rpc.Client[any],
	publisher bus.Publisher[entities.Transaction[any]],
	metrics TxIDHandlerMetrics,
) *TxIDHandler {
	return &TxIDHandler{
		cfg:       cfg,
		rpcClient: rpcClient,
		publisher: publisher,
		metrics:   metrics,
	}
}

func (uc *TxIDHandler) Handle(ctx context.Context, txid entities.TxID) error {
	logger := zerolog.Ctx(ctx).With().
		Str("name", "txid_handler").
		Any("driver", uc.cfg.Driver).
		Any("type", uc.cfg.Type).
		Any("chain", uc.cfg.Chain).
		Any("txid", txid).
		Logger()

	tx, err := uc.rpcClient.GetTransactionByTxID(ctx, txid)
	uc.metrics.RequestToNodeCounter.Inc(uc.cfg.Chain)
	if err != nil {
		return eris.Wrap(err, "get transaction by txid")
	}

	logger.Info().Msg("got transaction by hash")

	err = uc.publisher.PublishTo(ctx, uc.cfg.Topics.Transactions, *tx)
	if err != nil {
		return eris.Wrap(err, "publish mempool transaction")
	}

	return nil
}
