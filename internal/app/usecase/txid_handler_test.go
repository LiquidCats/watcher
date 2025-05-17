package usecase_test

import (
	"testing"

	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestTxIDHandler_Handle(t *testing.T) {
	ctx := t.Context()
	tx := &data.Transaction{
		TxID:          "tx_hash_1",
		Vin:           nil,
		Vout:          nil,
		Fee:           decimal.RequireFromString("0.01"),
		Confirmations: 0,
		BlockHash:     "hash1",
	}

	rpc := mocks.NewMockClient(t)
	ch := make(chan entities.Transaction, 1)

	rpc.On("GetTransactionByTxID", ctx, tx.TxID).Once().Return(tx, nil)

	uc := usecase.NewTxIDHandler(rpc, ch)

	err := uc.Handle(ctx, tx.TxID)
	require.NoError(t, err)

	tx1 := <-ch
	require.Equal(t, tx.GetTxID(), tx1.GetTxID())
}
