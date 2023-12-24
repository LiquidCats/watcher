package postgres

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
	"watcher/internal/app/domain/entity"
	"watcher/test/mock"
)

func TestStorageRepository_GetHeight(t *testing.T) {
	conn, dbMock := mock.MockDB(t)
	defer conn.Close()

	rows := dbMock.NewRows([]string{"height_value"}).
		AddRow(10000)

	dbMock.ExpectQuery("SELECT height_value FROM heights WHERE blockchain = \\$1 LIMIT 1").
		WithArgs(entity.Ethereum).
		WillReturnRows(rows)

	repo := NewStorageRepository(conn)

	ctx := context.Background()

	height, err := repo.GetHeight(ctx, entity.Ethereum)

	require.NoError(t, err)

	require.Equal(t, entity.BlockHeight(10000), height)
}

func TestStorageRepository_StoreBlock(t *testing.T) {
	conn, dbMock := mock.MockDB(t)
	defer conn.Close()

	block := &entity.Block{
		Height:   10000,
		Hash:     "current_hash",
		Previous: "previous_hash",
	}

	dbMock.ExpectExec("INSERT INTO blocks_bitcoin").
		WithArgs(10000, "current_hash", "previous_hash").
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewStorageRepository(conn)

	ctx := context.Background()

	err := repo.StoreBlock(ctx, entity.Bitcoin, block)

	require.NoError(t, err)
}

func TestStorageRepository_UpdateBlock(t *testing.T) {
	conn, dbMock := mock.MockDB(t)
	defer conn.Close()

	block := &entity.Block{
		Height:   10000,
		Hash:     "current_hash",
		Previous: "previous_hash",
	}

	dbMock.ExpectExec("UPDATE blocks_bitcoin SET hash = \\$1, previous = \\$2 WHERE height = \\$3").
		WithArgs("current_hash", "previous_hash", 10000).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewStorageRepository(conn)

	ctx := context.Background()

	err := repo.UpdateBlock(ctx, entity.Bitcoin, block)

	require.NoError(t, err)
}

func TestStorageRepository_GetBlock(t *testing.T) {
	conn, dbMock := mock.MockDB(t)
	defer conn.Close()

	rows := dbMock.NewRows([]string{"height", "hash", "previous"}).
		AddRow(10000, "current_hash", "previous_hash")

	dbMock.ExpectQuery("SELECT height, hash, previous FROM blocks_ethereum WHERE height = \\$1 LIMIT 1").
		WithArgs(10000).
		WillReturnRows(rows)

	repo := NewStorageRepository(conn)

	ctx := context.Background()

	block, err := repo.GetBlock(ctx, entity.Ethereum, entity.BlockHeight(10000))

	require.NoError(t, err)

	require.Equal(t, entity.BlockHeight(10000), block.Height)
	require.Equal(t, entity.BlockHash("current_hash"), block.Hash)
	require.Equal(t, entity.BlockHash("previous_hash"), block.Previous)
}

func TestStorageRepository_ConfirmBlock(t *testing.T) {
	conn, dbMock := mock.MockDB(t)
	defer conn.Close()

	dbMock.ExpectExec("UPDATE blocks_bitcoin SET is_confirmed = \\$1 WHERE height = \\$2").
		WithArgs(true, 10000).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewStorageRepository(conn)

	ctx := context.Background()

	err := repo.ConfirmBlock(ctx, entity.Bitcoin, entity.BlockHeight(10000))

	require.NoError(t, err)
}

func TestStorageRepository_GetAllUnconfirmedBlocks(t *testing.T) {
	conn, dbMock := mock.MockDB(t)
	defer conn.Close()

	rows := dbMock.NewRows([]string{"height", "hash", "previous"}).
		AddRow(10002, "current_hash_2", "previous_hash_2").
		AddRow(10001, "current_hash_1", "previous_hash_1").
		AddRow(10000, "current_hash_0", "previous_hash_0")

	dbMock.ExpectQuery("SELECT height, hash, previous FROM blocks_ethereum WHERE \\(height <= \\$1 AND is_confirmed = \\$2\\) ORDER BY height DESC").
		WithArgs(10000, false).
		WillReturnRows(rows)

	repo := NewStorageRepository(conn)

	ctx := context.Background()

	blocks, err := repo.GetAllUnconfirmedBlocks(ctx, entity.Ethereum, entity.BlockHeight(10000))

	require.NoError(t, err)

	require.Len(t, blocks, 3)

	require.Equal(t, entity.BlockHeight(10002), blocks[0].Height)
	require.Equal(t, entity.BlockHash("current_hash_2"), blocks[0].Hash)
	require.Equal(t, entity.BlockHash("previous_hash_2"), blocks[0].Previous)

	require.Equal(t, entity.BlockHeight(10000), blocks[2].Height)
	require.Equal(t, entity.BlockHash("current_hash_0"), blocks[2].Hash)
	require.Equal(t, entity.BlockHash("previous_hash_0"), blocks[2].Previous)

	require.Equal(t, entity.BlockHeight(10001), blocks[1].Height)
	require.Equal(t, entity.BlockHash("current_hash_1"), blocks[1].Hash)
	require.Equal(t, entity.BlockHash("previous_hash_1"), blocks[1].Previous)
}
