package state

import (
	"context"
	"database/sql"
	"sync/atomic"
	"time"

	db "github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/app/port/database"
	"github.com/bytedance/sonic"
	"github.com/go-faster/errors"
)

type stateData[T any] struct {
	value       []T
	lastUpdated time.Time
}

type Persister[T any] struct {
	persistedStorage database.StateDB
	data             atomic.Value
}

func NewPersister[T any](persistedStorage database.StateDB) *Persister[T] {
	p := &Persister[T]{
		persistedStorage: persistedStorage,
	}

	p.data.Store(stateData[T]{})

	return p
}

func (p *Persister[T]) Set(ctx context.Context, key string, value []T, period time.Duration) error {
	// Load the current snapshot
	oldH, ok := p.data.Load().(stateData[T])
	if !ok {
		return errors.New("persisted state not loaded")
	}

	h := stateData[T]{
		value:       value,
		lastUpdated: oldH.lastUpdated,
	}

	// Decide whether to push to database
	if oldH.value == nil || time.Since(oldH.lastUpdated) >= period {
		// marshal & persist
		data, err := sonic.Marshal(value)
		if err != nil {
			return err
		}
		if err = p.persistedStorage.SetState(ctx, db.SetStateParams{
			Key:   key,
			Value: data,
		}); err != nil {
			return errors.Wrap(err, "failed to persist state")
		}
		h.lastUpdated = time.Now()
	}

	// Atomically publish the updated in-memory state
	p.data.Store(h)
	return nil
}

func (p *Persister[T]) Get(ctx context.Context, key string) ([]T, error) {
	// Fast path: if we already have something in memory, just return it
	h, ok := p.data.Load().(stateData[T])
	if !ok {
		return nil, errors.New("persisted state not loaded")
	}

	if len(h.value) != 0 {
		return h.value, nil
	}

	state, err := p.persistedStorage.GetStateByKey(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "persisted state")
	}

	var value []T

	if err = sonic.Unmarshal(state.Value, &value); err != nil {
		return nil, errors.Wrap(err, "failed to decode state")
	}

	h = stateData[T]{
		value:       value,
		lastUpdated: state.UpdatedAt.Time,
	}
	p.data.Store(h)

	return value, nil
}
