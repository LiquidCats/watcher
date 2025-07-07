package usecase_test

import (
	"testing"
	"time"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestTxIDHandler_Handle(t *testing.T) {
	ctx := t.Context()
	cfg := configs.ChainConfig{
		Driver: entities.DriverRPC,
		Type:   entities.TypeUtxo,
		Chain:  "bitcoin",
		Scan: configs.ScanConfig{
			Depth: 2,
		},
		Persist: configs.PersistConfig{
			Capacity: 6,
			Duration: time.Hour,
		},
		Topics: configs.TopicsConfig{
			Transactions: "test-transactions",
		},
	}

	tx := &data.Transaction{
		TxID:          "tx_hash_1",
		Vin:           nil,
		Vout:          nil,
		Fee:           decimal.RequireFromString("0.01"),
		Confirmations: 0,
		BlockHash:     "hash1",
	}

	rpc := mocks.NewMockClient(t)
	publisher := mocks.NewMockPublisher[entities.Transaction](t)
	requestToNodeCounter := mocks.NewMockRequestToNodeCounter(t)

	rpc.On("GetTransactionByTxID", ctx, tx.TxID).Once().Return(tx, nil)
	publisher.On("PublishTo", ctx, cfg.Topics.Transactions, tx).Once().Return(nil)
	requestToNodeCounter.On("Inc", cfg.Chain).Once().Return(nil)

	uc := usecase.NewTxIDHandler(cfg, rpc, publisher, usecase.TxIDHandlerMetrics{
		RequestToNodeCounter: requestToNodeCounter,
	})

	err := uc.Handle(ctx, tx.TxID)
	require.NoError(t, err)
}
