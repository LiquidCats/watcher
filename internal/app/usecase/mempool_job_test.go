package usecase_test

import (
	"testing"
	"time"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWatchMempoolUseCase_Execute(t *testing.T) {
	cfg := configs.ChainConfig{
		Driver: entities.DriverRPC,
		Type:   entities.TypeUtxo,
		Chain:  "bitcoin",

		Persist: configs.PersistConfig{
			Interval: time.Hour,
		},
	}
	client := mocks.NewMockClient(t)
	requestToNodeCounter := mocks.NewMockRequestToNodeCounter(t)

	requestToNodeCounter.On("Inc", cfg.Chain).Once().Return(nil)

	testCh := make(chan entities.TxID, 2)
	uc := usecase.NewMempoolJob(cfg, client, testCh, []entities.TxID{}, usecase.MempoolJobMetrics{
		RequestToNodeCounter: requestToNodeCounter,
	})

	newMempool := []entities.TxID{"tx1", "tx3"}
	client.
		On("GetMempool", mock.Anything).
		Return(newMempool, nil)

	err := uc.Handle(t.Context())
	require.NoError(t, err)

	assert.Equal(t, entities.TxID("tx1"), <-testCh)
	assert.Equal(t, entities.TxID("tx3"), <-testCh)
}
