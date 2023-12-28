package usecase

import (
	"context"
	"go.uber.org/zap"
	"watcher/configs"
	"watcher/internal/port"
)

type CleanupUsecase struct {
	cfg     configs.Config
	storage port.StorageRepository
	logger  *zap.Logger
}

func NewCleanupUsecase(cfg configs.Config, storage port.StorageRepository, logger *zap.Logger) *CleanupUsecase {
	return &CleanupUsecase{
		cfg:     cfg,
		storage: storage,
		logger:  logger,
	}
}

func (c *CleanupUsecase) Clean(ctx context.Context) error {
	if !c.cfg.Cleanup.Enabled {
		return nil
	}
	return c.storage.Transaction(ctx, func(ctx context.Context) error {
		savedHeight, err := c.storage.GetHeight(ctx, c.cfg.Blockchain)
		if err != nil {
			return err
		}

		fromHeight := savedHeight - c.cfg.Gap - c.cfg.Cleanup.Gap

		if err := c.storage.RemoveConfirmedBlocks(ctx, c.cfg.Blockchain, fromHeight); err != nil {
			return err
		}

		c.logger.Info("confirmed blocks cleaned", zap.Int("from_block", int(fromHeight)))

		return nil
	})
}
