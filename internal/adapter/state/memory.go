package state

import (
	"sync/atomic"

	"github.com/LiquidCats/watcher/v2/configs"
)

type PersistedState[T any] struct {
	cfg configs.PersistConfig

	value atomic.Value
}

func NewMemoryState[T any](cfg configs.PersistConfig) *PersistedState[T] {
	p := &PersistedState[T]{
		cfg: cfg,
	}
	emptySlice := make([]T, 0, cfg.Capacity)

	p.value.Store(&emptySlice)
	return p
}

func (s *PersistedState[T]) Set(value T) {
	for {
		// Read current value
		oldValPtr := s.value.Load().(*[]T) //nolint:errcheck
		oldVal := *oldValPtr

		// Create a new slice with the added value
		newVal := make([]T, len(oldVal), len(oldVal)+1)
		copy(newVal, oldVal)
		newVal = append(newVal, value)

		// Enforce capacity limit
		if len(newVal) >= s.cfg.Capacity {
			newVal = newVal[1:]
		}

		// Try to atomically swap the value
		if s.value.CompareAndSwap(oldValPtr, &newVal) {
			return
		}
		// If CAS failed, retry (another goroutine modified the value)
	}
}

func (s *PersistedState[T]) Get() []T {
	valPtr := s.value.Load().(*[]T) //nolint:errcheck
	val := *valPtr

	if len(val) != 0 {
		// Return a copy to prevent external modifications
		result := make([]T, len(val))
		copy(result, val)
		return result
	}

	return nil
}
