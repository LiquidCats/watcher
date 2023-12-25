package usecase

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"watcher/configs"
	"watcher/internal/app/domain/entity"
	"watcher/test/mocks"
)

func TestBlockHandlingUsecase_Handle(t *testing.T) {
	ctx := context.Background()

	blockChan := make(chan entity.BlockHeight, 3)

	cfg := configs.Config{
		Blockchain: entity.Ethereum,
		NodeUrl:    "http://test",
		Interval:   10,
		Gap:        6,
		Workers:    3,
	}

	rpc := mocks.NewRpcRepository(t)
	storage := mocks.NewStorageRepository(t)
	publisher := mocks.NewEventPublisher(t)

	storedHeight := entity.BlockHeight(9999)

	blockchainBlock := &entity.Block{
		Height:   10000,
		Hash:     "current_hash",
		Previous: "previous_hash",
	}

	storage.On("GetHeight", ctx, entity.Ethereum).Once().Return(storedHeight, nil)
	rpc.On("GetLatestBlock", ctx).Once().Return(blockchainBlock, nil)
	rpc.On("GetBlockByHeight", ctx, blockchainBlock.Height).Once().Return(blockchainBlock, nil)
	storage.On("StoreBlock", ctx, entity.Ethereum, blockchainBlock).Once().Return(nil)
	publisher.On("NewBlock", ctx, entity.Ethereum, blockchainBlock).Once().Return(nil)

	usecase := NewBlockHandlingUsecase(cfg, rpc, storage, publisher)

	err := usecase.Handle(ctx, blockChan)

	require.Len(t, blockChan, 1)

	blockHeight := <-blockChan
	close(blockChan)

	require.NoError(t, err)
	require.Equal(t, entity.BlockHeight(10000), blockHeight)

}
