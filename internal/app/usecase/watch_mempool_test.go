package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWatchMempoolUseCase_Execute(t *testing.T) {
	cfg := configs.App{
		Driver: entities.DriverRPC,
		Type:   entities.TypeUtxo,
		Chain:  entities.ChainBitcoin,

		PersistDuration: time.Hour,
	}
	state := mocks.NewState[[]entities.TxID](t)
	transactionPublisher := mocks.NewTransactionPublisher(t)
	client := mocks.NewClient(t)

	uc := usecase.NewWatchMempoolUseCase(cfg, state, client, transactionPublisher)

	newMempool := []entities.TxID{"tx1", "tx3"}
	client.
		On("GetMempool", mock.Anything).
		Return(newMempool, nil)

	client.
		On("GetTransactionByTxId", mock.Anything, entities.TxID("tx1")).
		Return(&data.Transaction{Txid: "tx1"}, nil).
		On("GetTransactionByTxId", mock.Anything, entities.TxID("tx3")).
		Return(&data.Transaction{Txid: "tx3"}, nil)

	transactionPublisher.
		On("PublishTransaction", mock.Anything, mock.MatchedBy(func(tx *entities.UtxoTransaction) bool {
			return tx.TxID == "tx1"
		})).
		Return(nil).
		On("PublishTransaction", mock.Anything, mock.MatchedBy(func(tx *entities.UtxoTransaction) bool {
			return tx.TxID == "tx3"
		})).
		Return(nil)

	state.On("Get", mock.Anything, "rpc.utxo.bitcoin.mempool").Once().Return([]entities.TxID{}, nil)
	state.On("Set", mock.Anything, "rpc.utxo.bitcoin.mempool", newMempool, cfg.PersistDuration).
		Once().
		Return(nil)

	err := uc.Execute(context.Background())
	require.NoError(t, err)
}
