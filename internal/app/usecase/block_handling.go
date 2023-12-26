package usecase

import (
	"context"
	"go.uber.org/zap"
	"watcher/configs"
	"watcher/internal/app/domain/entity"
	"watcher/internal/port"
)

type BlockHandlingUsecase struct {
	cfg       configs.Config
	rpc       port.RpcRepository
	storage   port.StorageRepository
	publisher port.EventPublisher
	logger    *zap.Logger
}

func NewBlockHandlingUsecase(
	cfg configs.Config,
	rpc port.RpcRepository,
	storage port.StorageRepository,
	publisher port.EventPublisher,
	logger *zap.Logger,
) *BlockHandlingUsecase {
	return &BlockHandlingUsecase{
		cfg:       cfg,
		rpc:       rpc,
		storage:   storage,
		publisher: publisher,
		logger:    logger,
	}
}

func (w *BlockHandlingUsecase) Handle(ctx context.Context) error {
	savedHeight, err := w.storage.GetHeight(ctx, w.cfg.Blockchain)
	if err != nil {
		return err
	}

	latestBlock, err := w.rpc.GetLatestBlock(ctx)
	if err != nil {
		return err
	}

	if latestBlock.Height <= savedHeight {
		return nil
	}

	for blockHeight := savedHeight + 1; blockHeight <= latestBlock.Height; blockHeight++ {
		if err := w.handleBlock(ctx, blockHeight); err != nil {
			return err
		}
	}

	return nil
}

func (w *BlockHandlingUsecase) handleBlock(ctx context.Context, blockHeight entity.BlockHeight) error {
	return w.storage.Transaction(ctx, func(ctx context.Context) error {
		block, err := w.rpc.GetBlockByHeight(ctx, blockHeight)
		if err != nil {
			return err
		}
		if err := w.storage.StoreBlock(ctx, w.cfg.Blockchain, block); err != nil {
			return err
		}
		if err := w.storage.UpdateHeight(ctx, w.cfg.Blockchain, blockHeight); err != nil {
			return err
		}
		if err := w.publisher.NewBlock(ctx, w.cfg.Blockchain, block); err != nil {
			return err
		}

		w.logger.Info(
			"block handled",
			zap.String("block_hash", string(block.Hash)),
			zap.Int("block_height", int(block.Height)),
		)

		return nil
	})

}
