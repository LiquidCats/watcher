package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
	kernel "github.com/LiquidCats/watcher/v2/internal/app/utxo/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/utxo/usecase"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWatchMempoolUseCase_Execute(t *testing.T) {
	cfg := configs.App{
		Driver: entities.DriverRPC,
		Type:   entities.TypeUtxo,
		Chain:  entities.ChainBitcoin,
	}
	stateDB := mocks.NewStateDB(t)
	transactionPublisher := mocks.NewTransactionPublisher(t)
	client := mocks.NewClient(t)

	uc := usecase.NewWatchMempoolUseCase(cfg, stateDB, transactionPublisher, client)

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
		On("PublishTransaction", mock.Anything, mock.MatchedBy(func(tx *kernel.Transaction) bool {
			return tx.TxID == "tx1"
		})).
		Return(nil).
		On("PublishTransaction", mock.Anything, mock.MatchedBy(func(tx *kernel.Transaction) bool {
			return tx.TxID == "tx3"
		})).
		Return(nil)

	err := uc.Execute(context.Background())
	require.NoError(t, err)
}

// TestWatchMempoolUseCase_Init tests a successful initialization by
// simulating a stateDB that returns a pre-encoded mempool.
func TestWatchMempoolUseCase_Init(t *testing.T) {
	cfg := configs.App{
		Driver: entities.DriverRPC,
		Type:   entities.TypeUtxo,
		Chain:  entities.ChainBitcoin,
	}

	stateDB := mocks.NewStateDB(t)
	transactionPublisher := mocks.NewTransactionPublisher(t)
	client := mocks.NewClient(t)

	// Prepare an old mempool slice and encode it using sonic.
	oldMempool := []entities.TxID{"tx1", "tx2"}
	encoded, err := sonic.Marshal(oldMempool)
	require.NoError(t, err)

	// Expect stateDB.GetByKey to be called with the correct key.
	stateKey := fmt.Sprint(cfg.Driver, ".", cfg.Type, ".", cfg.Chain, ".", "mempool")
	stateDB.
		On("GetByKey", mock.Anything, stateKey).
		Return(database.State{Value: encoded}, nil)

	uc := usecase.NewWatchMempoolUseCase(cfg, stateDB, transactionPublisher, client)

	err = uc.Init(context.Background())
	require.NoError(t, err)
}
