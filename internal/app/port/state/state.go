package state

import (
	"context"
)

type State[T any] interface {
	Set(ctx context.Context, key string, value T) error
	Get(ctx context.Context, key string) ([]T, error)
}
