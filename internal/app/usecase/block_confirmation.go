package usecase

import (
	"context"
	"go.uber.org/zap"
	"watcher/configs"
	"watcher/internal/app/domain/entity"
	"watcher/internal/port"
)

type BlockConfirmationUsecase struct {
	cfg       configs.Config
	rpc       port.RpcRepository
	storage   port.StorageRepository
	publisher port.EventPublisher
	logger    *zap.Logger
}

func NewBlockConfirmationUsecase(
	cfg configs.Config,
	rpc port.RpcRepository,
	storage port.StorageRepository,
	publisher port.EventPublisher,
	logger *zap.Logger,
) *BlockConfirmationUsecase {
	return &BlockConfirmationUsecase{
		cfg:       cfg,
		rpc:       rpc,
		storage:   storage,
		publisher: publisher,
		logger:    logger,
	}
}

func (b *BlockConfirmationUsecase) Confirm(ctx context.Context) error {
	return b.storage.Transaction(ctx, func(ctx context.Context) error {
		savedHeight, err := b.storage.GetHeight(ctx, b.cfg.Blockchain)
		if err != nil {
			return err
		}

		blocksToConfirm, err := b.storage.GetAllUnconfirmedBlocks(ctx, b.cfg.Blockchain, savedHeight-b.cfg.Gap)

		for _, storedBlock := range blocksToConfirm {
			blockchainBlock, err := b.rpc.GetBlockByHeight(ctx, storedBlock.Height)
			if err != nil {
				return err
			}

			if blockchainBlock.Hash != storedBlock.Hash {
				if err := b.reject(ctx, storedBlock, blockchainBlock); err != nil {
					return err
				}

				b.logger.Info(
					"block rejected",
					zap.String("block_hash", string(storedBlock.Hash)),
					zap.Int("block_height", int(storedBlock.Height)),
				)
			}

			if err := b.confirm(ctx, blockchainBlock); err != nil {
				return err
			}

			b.logger.Info(
				"block confirmed",
				zap.String("block_hash", string(blockchainBlock.Hash)),
				zap.Int("block_height", int(blockchainBlock.Height)),
			)
		}

		return nil
	})
}

func (b *BlockConfirmationUsecase) reject(ctx context.Context, storedBlock, newBlock *entity.Block) error {
	if err := b.storage.UpdateBlock(ctx, b.cfg.Blockchain, newBlock); err != nil {
		return err
	}
	if err := b.publisher.NewBlock(ctx, b.cfg.Blockchain, newBlock); err != nil {
		return err
	}
	if err := b.publisher.RejectBlock(ctx, b.cfg.Blockchain, storedBlock); err != nil {
		return err
	}

	return nil
}

func (b *BlockConfirmationUsecase) confirm(ctx context.Context, block *entity.Block) error {
	if err := b.storage.ConfirmBlock(ctx, b.cfg.Blockchain, block.Height); err != nil {
		return err
	}

	if err := b.publisher.ConfirmBlock(ctx, b.cfg.Blockchain, block); err != nil {
		return err
	}

	return nil
}
