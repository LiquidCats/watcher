package usecase

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"watcher/configs"
	"watcher/internal/adapter/logger"
	"watcher/internal/app/domain/entity"
	mocks "watcher/test/mocks"
)

func TestBlockConfirmationUsecase_Confirm(t *testing.T) {
	cfg := configs.Config{
		Blockchain: entity.Bitcoin,
		NodeUrl:    "http://test",
		Interval:   10,
		Gap:        6,
	}

	savedHeight := entity.BlockHeight(10000)

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

		storage.On("Transaction", ctx, mock.Anything).Return(func(ctx context.Context, cb func(context.Context) error) error {
			return cb(ctx)
		})

		storage.On("GetHeight", ctx, cfg.Blockchain).
			Once().
			Return(savedHeight, nil)

		storage.On("GetAllUnconfirmedBlocks", ctx, cfg.Blockchain, savedHeight-cfg.Gap).
			Once().
			Return([]*entity.Block{
				storedBlock,
			}, nil)

		rpc.On("GetBlockByHeight", ctx, storedBlock.Height).
			Once().
			Return(blockchainBlock, nil)

		storage.On("ConfirmBlock", ctx, cfg.Blockchain, savedHeight-cfg.Gap).
			Once().
			Return(nil)

		publisher.On("ConfirmBlock", ctx, cfg.Blockchain, blockchainBlock).
			Once().
			Return(nil)

		usecase := NewBlockConfirmationUsecase(cfg, rpc, storage, publisher, logger.NewLogger("test"))

		err := usecase.Confirm(ctx)

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

		storage.On("Transaction", ctx, mock.Anything).Return(func(ctx context.Context, cb func(context.Context) error) error {
			return cb(ctx)
		})

		storage.On("GetHeight", ctx, cfg.Blockchain).
			Once().
			Return(savedHeight, nil)

		storage.On("GetAllUnconfirmedBlocks", ctx, cfg.Blockchain, savedHeight-cfg.Gap).
			Once().
			Return([]*entity.Block{
				storedBlock,
			}, nil)

		rpc.On("GetBlockByHeight", ctx, storedBlock.Height).
			Once().
			Return(blockchainBlock, nil)

		storage.On("UpdateBlock", ctx, entity.Bitcoin, blockchainBlock).Once().Return(nil)
		publisher.On("NewBlock", ctx, entity.Bitcoin, blockchainBlock).Once().Return(nil)
		publisher.On("RejectBlock", ctx, entity.Bitcoin, storedBlock).Once().Return(nil)

		storage.On("ConfirmBlock", ctx, entity.Bitcoin, entity.BlockHeight(9994)).Once().Return(nil)

		publisher.On("ConfirmBlock", ctx, entity.Bitcoin, blockchainBlock).Once().Return(nil)

		usecase := NewBlockConfirmationUsecase(cfg, rpc, storage, publisher, logger.NewLogger("test"))

		err := usecase.Confirm(ctx)

		require.NoError(t, err)
	})
}
