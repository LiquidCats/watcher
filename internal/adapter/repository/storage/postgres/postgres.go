package postgres

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"watcher/internal/app/domain/entity"
)

type StorageRepository struct {
	conn *sql.DB
}

const HeightsTable = "heights"

func NewStorageRepository(conn *sql.DB) *StorageRepository {
	return &StorageRepository{
		conn: conn,
	}
}

func (s *StorageRepository) Transaction(ctx context.Context, cb func(ctx context.Context) error) error {
	tx, err := s.conn.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if err := cb(ctx); err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

func (s *StorageRepository) GetHeight(ctx context.Context, blockchain entity.Blockchain) (entity.BlockHeight, error) {
	row := sq.Select("height_value").
		From(HeightsTable).
		Where(sq.Eq{"blockchain": blockchain}).
		Limit(1).
		RunWith(s.conn).
		PlaceholderFormat(sq.Dollar).
		QueryRowContext(ctx)

	var height entity.BlockHeight

	if err := row.Scan(&height); err != nil {
		return 0, err
	}

	return height, nil
}