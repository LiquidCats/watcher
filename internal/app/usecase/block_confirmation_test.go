package usecase

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"watcher/configs"
	"watcher/internal/app/domain/entity"
	mocks "watcher/test/mocks"
)

func TestBlockConfirmationUsecase_Confirm(t *testing.T) {
	cfg := configs.Config{
		Blockchain: entity.Bitcoin,
		NodeUrl:    "http://test",
		Interval:   10,
		Gap:        6,
		Workers:    3,
	}

	t.Run("confirm", func(t *testing.T) {
		rpc := mocks.NewRpcRepository(t)
		storage := mocks.NewStorageRepository(t)
		publisher := mocks.NewEventPublisher(t)

		ctx := context.Background()

		blockchainBlock := &entity.Block{
			Height:   9994,
			Hash:     "current_hash",
			Previous: "previous_hash",
		}
		storedBlock := &entity.Block{
			Height:   9994,
			Hash:     "current_hash",
			Previous: "previous_hash",
		}

		rpc.On("GetBlockByHeight", ctx, entity.BlockHeight(10000-6)).Once().Return(blockchainBlock, nil)

		storage.On("Transaction", ctx, mock.Anything).Return(func(ctx context.Context, cb func(context.Context) error) error {
			return cb(ctx)
		})
		storage.On("GetBlock", ctx, entity.Bitcoin, blockchainBlock.Height).Once().Return(storedBlock, nil)
		storage.On("ConfirmBlock", ctx, entity.Bitcoin, entity.BlockHeight(9994)).Once().Return(nil)

		publisher.On("ConfirmBlock", ctx, entity.Bitcoin, blockchainBlock).Once().Return(nil)

		usecase := NewBlockConfirmationUsecase(cfg, rpc, storage, publisher)

		err := usecase.Confirm(ctx, entity.BlockHeight(10000))

		require.NoError(t, err)
	})

	t.Run("reject", func(t *testing.T) {
		rpc := mocks.NewRpcRepository(t)
		storage := mocks.NewStorageRepository(t)
		publisher := mocks.NewEventPublisher(t)

		ctx := context.Background()

		blockchainBlock := &entity.Block{
			Height:   9994,
			Hash:     "current_hash",
			Previous: "previous_hash",
		}
		storedBlock := &entity.Block{
			Height:   9994,
			Hash:     "old_hash",
			Previous: "previous_hash",
		}

		rpc.On("GetBlockByHeight", ctx, entity.BlockHeight(10000-6)).Once().Return(blockchainBlock, nil)

		storage.On("Transaction", ctx, mock.Anything).Return(func(ctx context.Context, cb func(context.Context) error) error {
			return cb(ctx)
		})
		storage.On("GetBlock", ctx, entity.Bitcoin, blockchainBlock.Height).Once().Return(storedBlock, nil)
		storage.On("UpdateBlock", ctx, entity.Bitcoin, blockchainBlock).Once().Return(nil)
		publisher.On("NewBlock", ctx, entity.Bitcoin, blockchainBlock).Once().Return(nil)
		publisher.On("RejectBlock", ctx, entity.Bitcoin, storedBlock).Once().Return(nil)

		storage.On("ConfirmBlock", ctx, entity.Bitcoin, entity.BlockHeight(9994)).Once().Return(nil)

		publisher.On("ConfirmBlock", ctx, entity.Bitcoin, blockchainBlock).Once().Return(nil)

		usecase := NewBlockConfirmationUsecase(cfg, rpc, storage, publisher)

		err := usecase.Confirm(ctx, entity.BlockHeight(10000))

		require.NoError(t, err)
	})
}
