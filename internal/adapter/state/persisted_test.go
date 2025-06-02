package state_test

import (
	"testing"
	"time"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/adapter/state"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestState_Get(t *testing.T) {
	ctx := t.Context()
	stateDB := mocks.NewMockStateDB(t)
	stateDB.On("GetStateByKey", ctx, "test").Return(database.State{
		Key:       "test",
		Value:     []byte(`["test_value"]`),
		CreatedAt: pgtype.Timestamp{Time: time.Now().Add(-(time.Hour * 24)), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}, nil)

	st := state.NewPersister[string](configs.PersistConfig{Capacity: 6}, stateDB)

	val, err := st.Get(ctx, "test")
	require.NoError(t, err)

	require.Equal(t, "test_value", val[0])
}

func TestState_Set(t *testing.T) {
	ctx := t.Context()
	stateDB := mocks.NewMockStateDB(t)
	stateDB.On("SetState", ctx, database.SetStateParams{
		Key:   "test",
		Value: []byte(`["test_value"]`),
	}).Return(nil)

	st := state.NewPersister[string](configs.PersistConfig{Capacity: 6}, stateDB)

	err := st.Set(ctx, "test", "test_value")
	require.NoError(t, err)
}
