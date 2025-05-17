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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWatchBlocksUseCase_Execute(t *testing.T) {
	cfg := configs.App{
		Driver:          entities.DriverRPC,
		Type:            entities.TypeUtxo,
		Chain:           "bitcoin",
		ScanDepth:       2,
		PersistBocks:    6,
		PersistDuration: time.Hour,
	}

	state := mocks.NewMockState[entities.BlockHash](t)
	client := mocks.NewMockClient(t)

	block1 := &data.Block{
		Hash:              "block1",
		Height:            1,
		PreviousBlockHash: "block0",
		Tx: []*data.Transaction{
			{
				TxID:          "tx1",
				Vin:           nil,
				Vout:          nil,
				Fee:           decimal.RequireFromString("0.001"),
				Confirmations: 1,
				BlockHash:     "block1",
			},
		},
	}
	block2 := &data.Block{
		Hash:              "block2",
		Height:            2,
		PreviousBlockHash: "block1",
		Tx: []*data.Transaction{
			{
				TxID:          "tx2",
				Vin:           nil,
				Vout:          nil,
				Fee:           decimal.RequireFromString("0.001"),
				Confirmations: 2,
				BlockHash:     "block2",
			},
		},
	}

	state.On("Get", mock.Anything, "utxo.rpc.bitcoin.blocks").Once().Return(nil, nil)
	client.On("GetLatestBlockHash", mock.Anything).Once().Return(block2.Hash, nil)
	client.On("GetBlockByHash", mock.Anything, block2.Hash).Once().Return(block2, nil)
	client.On("GetBlockByHash", mock.Anything, block1.Hash).Once().Return(block1, nil)

	testCh := make(chan entities.Block, 2)
	defer close(testCh)

	uc := usecase.NewBlocksJob(cfg, state, client, testCh)

	err := uc.Handle(t.Context())
	require.NoError(t, err)

	b1, b2 := <-testCh, <-testCh

	assert.Equal(t, block1.GetHash(), b1.GetHash())
	assert.Equal(t, block2.GetHash(), b2.GetHash())
}
