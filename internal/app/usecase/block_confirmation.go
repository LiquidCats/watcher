package usecase

import (
	"context"
	"watcher/configs"
	"watcher/internal/app/domain/entity"
	"watcher/internal/port"
)

type BlockConfirmationUsecase struct {
	cfg       configs.Config
	rpc       port.RpcRepository
	storage   port.StorageRepository
	publisher port.EventPublisher
}

func NewBlockConfirmationUsecase(
	cfg configs.Config,
	rpc port.RpcRepository,
	storage port.StorageRepository,
	publisher port.EventPublisher,
) *BlockConfirmationUsecase {
	return &BlockConfirmationUsecase{
		cfg:       cfg,
		rpc:       rpc,
		storage:   storage,
		publisher: publisher,
	}
}

func (b *BlockConfirmationUsecase) Confirm(ctx context.Context, savedHeight entity.BlockHeight) error {
	confirmationHeight := savedHeight - b.cfg.Gap

	block, err := b.rpc.GetBlockByHeight(ctx, confirmationHeight)
	if err != nil {
		return err
	}

	return b.storage.Transaction(ctx, func(ctx context.Context) error {
		storedBlock, err := b.storage.GetBlock(ctx, b.cfg.Blockchain, confirmationHeight)
		if err != nil {
			return err
		}

		if block.Hash != storedBlock.Hash {
			if err := b.reject(ctx, storedBlock, block); err != nil {
				return err
			}
		}

		if err := b.confirm(ctx, block); err != nil {
			return err
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
