package postgres

import (
	"context"
	"database/sql"
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
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}
