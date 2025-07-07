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

type BlockTransactionsHandler struct {
	cfg       configs.ChainConfig
	rpc       rpc.Client
	publisher bus.Publisher[entities.Transaction]
	metrics   BlockTransactionsHandlerMetrics
}

type BlockTransactionsHandlerMetrics struct {
	RequestToNodeCounter metrics.RequestToNodeCounter
}

func NewBlockTransactionsHandler(
	cfg configs.ChainConfig,
	rpc rpc.Client,
	transactionPub bus.Publisher[entities.Transaction],
	metrics BlockTransactionsHandlerMetrics,
) *BlockTransactionsHandler {
	return &BlockTransactionsHandler{
		cfg:       cfg,
		rpc:       rpc,
		publisher: transactionPub,
		metrics:   metrics,
	}
}

func (uc *BlockTransactionsHandler) Handle(ctx context.Context, block entities.Block) error {
	logger := zerolog.Ctx(ctx).With().
		Str("name", "transaction_handler").
		Any("driver", uc.cfg.Driver).
		Any("type", uc.cfg.Type).
		Any("chain", uc.cfg.Chain).
		Any("block_hash", block.GetHash()).
		Logger()

	block, err := uc.rpc.GetBlockByHash(ctx, block.GetHash(), true)
	uc.metrics.RequestToNodeCounter.Inc(uc.cfg.Chain)
	if err != nil {
		return eris.Wrap(err, "get block by hash")
	}

	for _, transaction := range block.GetTransactions() {
		logger.Debug().Any("txid", transaction.GetTxID()).Msg("publish transaction")
		err := uc.publisher.PublishTo(ctx, uc.cfg.Topics.Transactions, transaction)
		if err != nil {
			logger.Error().
				Any("err", eris.ToJSON(err, true)).
				Any("txid", transaction.GetTxID()).
				Msg("publish transaction")
		}
	}

	logger.Debug().Int("txs", len(block.GetTransactions())).Msg("block transactions handled")

	return nil
}
