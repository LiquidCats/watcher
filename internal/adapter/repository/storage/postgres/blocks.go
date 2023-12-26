package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"watcher/internal/app/domain/entity"
)

const BlocksTable = "blocks"

func (s *StorageRepository) StoreBlock(ctx context.Context, blockchain entity.Blockchain, block *entity.Block) error {
	_, err := sq.Insert(fmt.Sprint(BlocksTable, "_", blockchain)).
		Columns("height", "hash", "previous").
		Values(block.Height, block.Hash, block.Previous).
		RunWith(s.conn).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (s *StorageRepository) UpdateBlock(ctx context.Context, blockchain entity.Blockchain, newBlock *entity.Block) error {
	_, err := sq.Update(fmt.Sprint(BlocksTable, "_", blockchain)).
		Where(sq.Eq{"height": newBlock.Height}).
		Set("hash", newBlock.Hash).
		Set("previous", newBlock.Previous).
		RunWith(s.conn).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (s *StorageRepository) GetBlock(ctx context.Context, blockchain entity.Blockchain, height entity.BlockHeight) (*entity.Block, error) {
	row := sq.Select("height", "hash", "previous").
		From(fmt.Sprint(BlocksTable, "_", blockchain)).
		Where(sq.Eq{"height": height}).
		Limit(1).
		RunWith(s.conn).
		PlaceholderFormat(sq.Dollar).
		QueryRowContext(ctx)

	var block entity.Block

	if err := row.Scan(&block.Height, &block.Hash, &block.Previous); err != nil {
		return nil, err
	}

	return &block, nil
}

func (s *StorageRepository) GetAllUnconfirmedBlocks(ctx context.Context, blockchain entity.Blockchain, fromHeight entity.BlockHeight) ([]*entity.Block, error) {
	rows, err := sq.Select("height", "hash", "previous").
		From(fmt.Sprint(BlocksTable, "_", blockchain)).
		Where(sq.And{
			sq.LtOrEq{"height": fromHeight},
			sq.Eq{"is_confirmed": false},
		}).
		OrderBy("height").
		RunWith(s.conn).
		PlaceholderFormat(sq.Dollar).
		QueryContext(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var blocks []*entity.Block

	for rows.Next() {
		var block entity.Block

		if err := rows.Scan(&block.Height, &block.Hash, &block.Previous); err != nil {
			return nil, err
		}

		blocks = append(blocks, &block)
	}

	return blocks, nil
}

func (s *StorageRepository) ConfirmBlock(ctx context.Context, blockchain entity.Blockchain, height entity.BlockHeight) error {
	_, err := sq.Update(fmt.Sprint(BlocksTable, "_", blockchain)).
		Where(sq.Eq{"height": height}).
		Set("is_confirmed", true).
		RunWith(s.conn).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)

	if err != nil {
		return err
	}

	return nil
}
