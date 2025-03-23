package state

import (
	"context"
	"time"
)

type State[T any] interface {
	Set(ctx context.Context, key string, value []T, period time.Duration) error
	Get(ctx context.Context, key string) ([]T, error)
}
