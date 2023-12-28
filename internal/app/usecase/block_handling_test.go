package usecase

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"watcher/configs"
	"watcher/internal/adapter/logger"
	"watcher/internal/app/domain/entity"
	"watcher/test/mocks"
)

func TestBlockHandlingUsecase_Handle(t *testing.T) {
	ctx := context.Background()

	cfg := configs.Config{
		Blockchain: entity.Ethereum,
		NodeUrl:    "http://test",
		Interval:   10,
		Gap:        6,
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
	storage.On("Transaction", ctx, mock.Anything).Return(func(ctx context.Context, cb func(context.Context) error) error {
		return cb(ctx)
	})
	rpc.On("GetBlockByHeight", ctx, blockchainBlock.Height).Once().Return(blockchainBlock, nil)
	storage.On("StoreBlock", ctx, entity.Ethereum, blockchainBlock).Once().Return(nil)
	storage.On("UpdateHeight", ctx, entity.Ethereum, blockchainBlock.Height).Once().Return(nil)
	publisher.On("NewBlock", ctx, entity.Ethereum, blockchainBlock).Once().Return(nil)

	usecase := NewBlockHandlingUsecase(cfg, rpc, storage, publisher, logger.NewLogger("test"))

	err := usecase.Handle(ctx)

	require.NoError(t, err)

}
