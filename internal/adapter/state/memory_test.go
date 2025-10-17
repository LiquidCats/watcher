package state_test

import (
	"testing"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/state"
	"github.com/stretchr/testify/assert"
)

func TestState_Get(t *testing.T) {
	st := state.NewMemoryState[string](configs.PersistConfig{Capacity: 6})

	st.Set("test_value1")
	st.Set("test_value2")
	st.Set("test_value3")

	val := st.Get()
	assert.Len(t, val, 3)
	assert.Equal(t, "test_value1", val[0])
	assert.Equal(t, "test_value2", val[1])
	assert.Equal(t, "test_value3", val[2])
}
