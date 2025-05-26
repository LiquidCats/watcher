package usecase_test

import (
	"testing"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTransactionHandler_Handle(t *testing.T) {
	ctx := t.Context()
	tx := &data.Transaction{
		TxID:          "tx_hash_1",
		Vin:           nil,
		Vout:          nil,
		Fee:           decimal.Decimal{},
		Confirmations: 0,
		BlockHash:     "",
	}
	cfg := configs.ChainConfig{
		Topics: configs.TopicsConfig{
			Transactions: "test-transactions",
		},
	}

	pub := mocks.NewMockPublisher[entities.Transaction](t)

	pub.On("PublishTo", mock.Anything, cfg.Topics.Transactions, tx).Return(nil)

	uc := usecase.NewTransactionHandler(cfg, pub)

	err := uc.Handle(ctx, tx)
	require.NoError(t, err)
}
