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

func TestCleanupUsecase_Clean(t *testing.T) {
	ctx := context.Background()

	cfg := configs.Config{
		Blockchain: entity.Ethereum,
		Gap:        6,
		Cleanup: configs.Cleanup{
			Enabled: true,
			Gap:     10,
		},
	}

	storedHeight := entity.BlockHeight(10000)
	fromHeight := entity.BlockHeight(10000) - cfg.Gap - cfg.Cleanup.Gap

	storage := mocks.NewStorageRepository(t)

	storage.On("Transaction", ctx, mock.Anything).Return(func(ctx context.Context, cb func(context.Context) error) error {
		return cb(ctx)
	})

	storage.On("GetHeight", ctx, cfg.Blockchain).Once().Return(storedHeight, nil)
	storage.On("RemoveConfirmedBlocks", ctx, cfg.Blockchain, fromHeight).Once().Return(nil)

	usecase := NewCleanupUsecase(cfg, storage, logger.NewLogger("test"))

	err := usecase.Clean(ctx)

	require.NoError(t, err)
}
