package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"watcher/internal/app/domain/entity"
)

func (s *StorageRepository) UpdateHeight(ctx context.Context, blockchain entity.Blockchain, height entity.BlockHeight) error {
	_, err := sq.Update(HeightsTable).
		Where(sq.Eq{"blockchain": blockchain}).
		Set("height_value", height).
		RunWith(s.conn).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)

	return err
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
