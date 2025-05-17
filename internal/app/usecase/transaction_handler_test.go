package usecase_test

import (
	"context"
	"testing"

	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTransactionHandler_Handle(t *testing.T) {
	ctx := context.Background()
	tx := &data.Transaction{
		TxID:          "tx_hash_1",
		Vin:           nil,
		Vout:          nil,
		Fee:           decimal.Decimal{},
		Confirmations: 0,
		BlockHash:     "",
	}

	pub := mocks.NewMockTransactionPublisher(t)

	pub.On("PublishTransaction", mock.Anything, tx).Return(nil)

	uc := usecase.NewTransactionHandler(pub)

	err := uc.Handle(ctx, tx)
	require.NoError(t, err)

}
