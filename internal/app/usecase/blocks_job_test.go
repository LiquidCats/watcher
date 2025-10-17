package usecase_test

import (
	"testing"
	"time"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatchBlocksUseCase_Execute(t *testing.T) {
	cfg := configs.ChainConfig{
		Driver: entities.DriverRPC,
		Type:   entities.TypeUtxo,
		Chain:  "bitcoin",
		Scan: configs.ScanConfig{
			Depth: 2,
		},
		Persist: configs.PersistConfig{
			Capacity: 6,
			Interval: time.Hour,
		},
		Topics: configs.TopicsConfig{
			Blocks: "blocks",
		},
	}

	state := mocks.NewMockState[entities.BlockHash](t)
	client := mocks.NewMockClient(t)
	publisher := mocks.NewMockPublisher[entities.Block](t)
	requestToNodeCounter := mocks.NewMockRequestToNodeCounter(t)

	block1 := &data.Block[*data.Transaction]{
		Hash: "block1",
	}
	block2 := &data.Block[*data.Transaction]{
		Hash:              "block2",
		Height:            2,
		PreviousBlockHash: "block1",
	}

	state.On("Get").Once().Return([]entities.BlockHash{block1.Hash})
	state.On("Set", block2.Hash).Once().Return()
	client.On("GetLatestBlockHash", t.Context()).Once().Return(block2.Hash, nil)
	client.On("GetBlockByHash", t.Context(), block2.Hash, false).Once().Return(block2, nil)
	publisher.On("PublishTo", t.Context(), "blocks", block2).Once().Return(nil)
	requestToNodeCounter.On("Inc", cfg.Chain).Twice().Return(nil)

	testCh := make(chan entities.Block, 2)
	defer close(testCh)

	uc := usecase.NewBlocksJob(cfg, state, client, testCh, publisher, usecase.BlocksJobMetrics{
		RequestToNodeCounter: requestToNodeCounter,
	})

	err := uc.Handle(t.Context())
	require.NoError(t, err)

	b2 := <-testCh

	assert.Equal(t, block2.GetHash(), b2.GetHash())
}
