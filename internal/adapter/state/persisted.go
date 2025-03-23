package state

import (
	"context"
	"database/sql"
	"sync"
	"time"

	database2 "github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/app/port/database"
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

type PersistedStateService[T any] struct {
	persistedStorage database.StateDB
	value            []T
	lastUpdated      time.Time
	mu               sync.Mutex
}

func NewPersisterState[T any](persistedStorage database.StateDB) *PersistedStateService[T] {
	return &PersistedStateService[T]{
		persistedStorage: persistedStorage,
	}
}

func (s *PersistedStateService[T]) Set(ctx context.Context, key string, value []T, period time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.shouldPersist(period) {
		valueBytes, err := sonic.Marshal(value)
		if err != nil {
			return err
		}

		if err := s.persistedStorage.SetState(ctx, database2.SetStateParams{
			Key:   key,
			Value: valueBytes,
		}); err != nil {
			return errors.Wrap(err, "failed to persist state")
		}
		s.lastUpdated = time.Now()
	}

	s.value = value

	return nil
}

func (s *PersistedStateService[T]) shouldPersist(period time.Duration) bool {
	if s.value == nil || len(s.value) == 0 {
		return true
	}

	return time.Since(s.lastUpdated) >= period
}

func (s *PersistedStateService[T]) Get(ctx context.Context, key string) ([]T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.value) != 0 {
		return s.value, nil
	}

	state, err := s.persistedStorage.GetStateByKey(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "persisted state")
	}

	var value []T

	if err := sonic.Unmarshal(state.Value, &value); err != nil {
		return nil, errors.Wrap(err, "failed to decode state")
	}

	s.lastUpdated = state.UpdatedAt.Time

	return value, nil
}
