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
	block := &data.Block[*data.Transaction]{
		Hash: "test_hash_1",
		Tx:   []*data.Transaction{tx},
	}
	cfg := configs.ChainConfig{
		Driver: entities.DriverRPC,
		Type:   entities.TypeUtxo,
		Chain:  "bitcoin",
		Topics: configs.TopicsConfig{
			Transactions: "test-transactions",
		},
	}

	pub := mocks.NewMockPublisher[entities.Transaction](t)
	rpc := mocks.NewMockClient(t)
	requestToNodeCounter := mocks.NewMockRequestToNodeCounter(t)

	rpc.On("GetBlockByHash", mock.Anything, block.GetHash(), true).Once().Return(block, nil)
	pub.On("PublishTo", mock.Anything, cfg.Topics.Transactions, tx).Once().Return(nil)
	requestToNodeCounter.On("Inc", cfg.Chain).Once().Return(nil)

	uc := usecase.NewBlockTransactionsHandler(cfg, rpc, pub, usecase.BlockTransactionsHandlerMetrics{
		RequestToNodeCounter: requestToNodeCounter,
	})

	err := uc.Handle(ctx, block)
	require.NoError(t, err)
}
