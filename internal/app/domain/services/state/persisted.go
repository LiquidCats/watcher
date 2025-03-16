package state

import (
	"context"
	"reflect"
	"time"

	database2 "github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/app/port/database"
	"github.com/bytedance/sonic"
	"github.com/go-faster/errors"
)

type PersistedStateService[T any] struct {
	persistedStorage database.StateDB
	value            T
	lastUpdated      time.Time
}

func NewPersisterState[T any](persistedStorage database.StateDB) *PersistedStateService[T] {
	return &PersistedStateService[T]{
		persistedStorage: persistedStorage,
	}
}

func (s *PersistedStateService[T]) Set(ctx context.Context, key string, value T, period time.Duration) error {
	s.value = value
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

	return nil
}

func (s *PersistedStateService[T]) shouldPersist(period time.Duration) bool {
	var zero T
	if reflect.DeepEqual(s.value, zero) {
		return true
	}

	return time.Since(s.lastUpdated) >= period
}

func (s *PersistedStateService[T]) Get(ctx context.Context, key string) (T, error) {
	var zero T
	if !reflect.DeepEqual(s.value, zero) {
		return s.value, nil
	}

	state, err := s.persistedStorage.GetStateByKey(ctx, key)
	if err != nil {
		return *new(T), errors.Wrap(err, "persisted state")
	}

	var value T

	if err := sonic.Unmarshal(state.Value, &value); err != nil {
		return *new(T), errors.Wrap(err, "failed to decode state")
	}

	return value, nil
}
