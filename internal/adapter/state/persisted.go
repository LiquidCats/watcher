package state

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/LiquidCats/watcher/v2/configs"
	db "github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/app/port/database"
	"github.com/bytedance/sonic"
	"github.com/rotisserie/eris"
)

type PersistedState[T any] struct {
	cfg              configs.PersistConfig
	persistedStorage database.StateDB

	value       []T
	lastUpdated time.Time

	mu sync.Mutex
}

func NewPersister[T any](cfg configs.PersistConfig, persistedStorage database.StateDB) *PersistedState[T] {
	return &PersistedState[T]{
		cfg:              cfg,
		persistedStorage: persistedStorage,
	}
}

func (s *PersistedState[T]) Set(ctx context.Context, key string, value T, period time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.value = append(s.value, value)
	if len(s.value) >= s.cfg.Capacity {
		s.value = s.value[1:]
	}

	if s.shouldPersist(period) {
		valueBytes, err := sonic.Marshal(s.value)
		if err != nil {
			return err
		}

		if err = s.persistedStorage.SetState(ctx, db.SetStateParams{
			Key:   key,
			Value: valueBytes,
		}); err != nil {
			return eris.Wrap(err, "failed to persist state")
		}
		s.lastUpdated = time.Now()
	}

	return nil
}

func (s *PersistedState[T]) shouldPersist(period time.Duration) bool {
	if s.value == nil {
		return true
	}

	return time.Since(s.lastUpdated).Seconds() >= period.Seconds()
}

func (s *PersistedState[T]) Get(ctx context.Context, key string) ([]T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.value) != 0 {
		return s.value, nil
	}

	state, err := s.persistedStorage.GetStateByKey(ctx, key)
	if err != nil {
		if eris.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, eris.Wrap(err, "persisted state")
	}

	var value []T

	if err = sonic.Unmarshal(state.Value, &value); err != nil {
		return nil, eris.Wrap(err, "failed to decode state")
	}

	s.value = value[:]
	s.lastUpdated = state.UpdatedAt.Time

	return s.value, nil
}
