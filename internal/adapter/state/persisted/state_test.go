package persisted_test

import (
	"context"
	"testing"
	"time"

	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/adapter/state/persisted"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestState_Get(t *testing.T) {
	ctx := context.Background()
	stateDB := mocks.NewStateDB(t)
	stateDB.On("GetStateByKey", ctx, "test").Return(database.State{
		Key:       "test",
		Value:     []byte(`"test_value"`),
		CreatedAt: pgtype.Timestamp{Time: time.Now().Add(-(time.Hour * 24)), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}, nil)

	state := persisted.NewPersisterState[string](stateDB)

	val, err := state.Get(ctx, "test")
	require.NoError(t, err)

	require.Equal(t, "test_value", val)
}

func TestState_Set(t *testing.T) {
	ctx := context.Background()
	stateDB := mocks.NewStateDB(t)
	stateDB.On("SetState", ctx, database.SetStateParams{
		Key:   "test",
		Value: []byte(`"test_value"`),
	}).Return(nil)

	state := persisted.NewPersisterState[string](stateDB)

	err := state.Set(ctx, "test", "test_value", time.Second)
	require.NoError(t, err)
}
