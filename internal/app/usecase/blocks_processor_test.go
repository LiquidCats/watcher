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

	state := mocks.NewState[[]entities.BlockHash](t)
	transactionPublisher := mocks.NewTransactionPublisher(t)
	blockPublisher := mocks.NewBlockPublisher(t)
	client := mocks.NewClient(t)

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

	client.On("GetLatestBlockHash", mock.Anything).Once().Return(block2.Hash, nil)
	client.On("GetBlockByHash", mock.Anything, block2.Hash).Once().Return(block2, nil)
	client.On("GetBlockByHash", mock.Anything, block1.Hash).Once().Return(block1, nil)

	transactionPublisher.On("PublishTransaction", mock.Anything, block1.Tx[0]).Once().Return(nil)
	blockPublisher.On("PublishBlock", mock.Anything, block1).Once().Return(nil)

	transactionPublisher.On("PublishTransaction", mock.Anything, block2.Tx[0]).Once().Return(nil)
	blockPublisher.On("PublishBlock", mock.Anything, block2).Once().Return(nil)

	state.On("Get", mock.Anything, "rpc.utxo.bitcoin.blocks").Once().Return([]entities.BlockHash{
		block1.PreviousBlockHash,
	}, nil)
	state.On("Set", mock.Anything, "rpc.utxo.bitcoin.blocks", []entities.BlockHash{
		block1.PreviousBlockHash,
		block1.Hash,
	}, cfg.PersistDuration).Once().Return(nil)
	state.On("Set", mock.Anything, "rpc.utxo.bitcoin.blocks", []entities.BlockHash{
		block1.PreviousBlockHash,
		block1.Hash,
		block2.Hash,
	}, cfg.PersistDuration).Once().Return(nil)

	uc := usecase.NewBlocksProcessor(cfg, state, client, blockPublisher, transactionPublisher)

	err := uc.Execute(t.Context())
	require.NoError(t, err)
}
