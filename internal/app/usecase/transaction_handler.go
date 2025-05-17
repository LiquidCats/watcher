package usecase

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/bus"
	"github.com/go-faster/errors"
	"github.com/rs/zerolog"
)

type TransactionHandler struct {
	transactionPub bus.TransactionPublisher
}

func NewTransactionHandler(transactionPub bus.TransactionPublisher) *TransactionHandler {
	return &TransactionHandler{
		transactionPub: transactionPub,
	}
}

func (uc *TransactionHandler) Handle(ctx context.Context, transaction entities.Transaction) error {
	logger := zerolog.Ctx(ctx).With().Any("txid", transaction.GetBlockHash()).Logger()

	err := uc.transactionPub.PublishTransaction(ctx, transaction)
	if err != nil {
		return errors.Wrap(err, "publish transaction")
	}

	logger.Info().Msg("published transaction")

	return nil
}
