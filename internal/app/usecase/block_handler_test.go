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

func TestBlockHandler_Handle(t *testing.T) {
	ctx := t.Context()
	cfg := configs.ChainConfig{
		Driver: entities.DriverRPC,
		Type:   entities.TypeUtxo,
		Chain:  "bitcoin",
		Persist: configs.PersistConfig{
			Capacity: 3,
			Duration: time.Hour,
		},
	}
	block := &data.Block{
		Hash:              "hash3",
		Height:            3,
		PreviousBlockHash: "hash2",
		Tx: []*data.Transaction{
			{
				TxID:          "tx_hash_1",
				Vin:           nil,
				Vout:          nil,
				Fee:           decimal.RequireFromString("0.1"),
				Confirmations: 1,
				BlockHash:     "hash3",
			},
			{
				TxID:          "tx_hash_2",
				Vin:           nil,
				Vout:          nil,
				Fee:           decimal.RequireFromString("0.33"),
				Confirmations: 1,
				BlockHash:     "hash3",
			},
		},
	}

	pub := mocks.NewMockPublisher[entities.Block](t)
	st := mocks.NewMockState[entities.BlockHash](t)
	ch := make(chan entities.Transaction, 2)

	pub.On("PublishTo", mock.Anything, cfg.Topics.Blocks, block).Once().Return(nil)
	st.On("Get", mock.Anything, "utxo.rpc.bitcoin.blocks").Return([]entities.BlockHash{"hash1", "hash2"}, nil)
	st.On("Set", mock.Anything, "utxo.rpc.bitcoin.blocks", []entities.BlockHash{"hash1", "hash2", "hash3"}, time.Hour).Return(nil)

	uc := usecase.NewBlockHandler(cfg, pub, st, ch)

	err := uc.Handle(ctx, block)
	require.NoError(t, err)

	tx1, tx2 := <-ch, <-ch
	require.Equal(t, entities.TxID("tx_hash_1"), tx1.GetTxID())
	require.Equal(t, entities.TxID("tx_hash_2"), tx2.GetTxID())
}
