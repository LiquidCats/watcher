package usecase

import (
	"context"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/bus"

	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
)

type TransactionHandler struct {
	cfg            configs.ChainConfig
	transactionPub bus.Publisher[entities.Transaction]
}

func NewTransactionHandler(
	cfg configs.ChainConfig,
	transactionPub bus.Publisher[entities.Transaction],
) *TransactionHandler {
	return &TransactionHandler{
		cfg:            cfg,
		transactionPub: transactionPub,
	}
}

func (uc *TransactionHandler) Handle(ctx context.Context, transaction entities.Transaction) error {
	logger := zerolog.Ctx(ctx).With().Any("txid", transaction.GetBlockHash()).Logger()

	err := uc.transactionPub.PublishTo(ctx, uc.cfg.Topics.Transactions, transaction)
	if err != nil {
		return eris.Wrap(err, "publish transaction")
	}

	logger.Info().Msg("published transaction")

	return nil
}
