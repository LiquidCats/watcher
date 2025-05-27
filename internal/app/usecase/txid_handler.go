package usecase

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
	"github.com/LiquidCats/watcher/v2/internal/app/port/runner"
	"github.com/rotisserie/eris"

	"github.com/rs/zerolog"
)

type TxIDHandler struct {
	rpcClient     rpc.Client
	transactionCh runner.ChanWrite[entities.Transaction]
}

func NewTxIDHandler(rpcClient rpc.Client, transactionCh runner.ChanWrite[entities.Transaction]) *TxIDHandler {
	return &TxIDHandler{
		rpcClient:     rpcClient,
		transactionCh: transactionCh,
	}
}

func (uc *TxIDHandler) Handle(ctx context.Context, txid entities.TxID) error {
	logger := zerolog.Ctx(ctx).With().Any("txid", txid).Logger()

	tx, err := uc.rpcClient.GetTransactionByTxID(ctx, txid)
	if err != nil {
		return eris.Wrap(err, "get transaction by txid")
	}

	logger.Info().Msg("got transaction by hash")

	uc.transactionCh <- tx

	return nil
}
