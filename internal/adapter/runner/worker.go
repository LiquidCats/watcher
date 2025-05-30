package runner

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/port/runner"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type Worker[T any] struct {
	name         string
	workerCh     chan T
	handler      runner.Handler[T]
	workersCount uint
}

func NewWorker[T any](
	name string,
	workersCount uint,
	workerCh chan T,
	handler runner.Handler[T],
) *Worker[T] {
	return &Worker[T]{
		name:         name,
		workersCount: workersCount,
		workerCh:     workerCh,
		handler:      handler,
	}
}

func (w *Worker[T]) Run(ctx context.Context) error {
	defer close(w.workerCh)

	logger := zerolog.Ctx(ctx).
		With().
		Caller().
		Str("worker_name", w.name).
		Uint("workers_count", w.workersCount).
		Logger()

	g, ctx := errgroup.WithContext(ctx)

	for range w.workersCount {
		g.Go(func() error {
			return w.runner(ctx)
		})
	}

	logger.Info().Msg("background workers started")
	defer logger.Info().Msg("background workers stopped")

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (w *Worker[T]) runner(ctx context.Context) error {
	logger := zerolog.Ctx(ctx).
		With().
		Caller().
		Str("worker_name", w.name).
		Logger()

	logger.Debug().Msg("runner started")
	defer logger.Debug().Msg("runner stopped")

	for {
		select {
		case <-ctx.Done():
			logger.Debug().Msg("context closed")
			return ctx.Err()
		case v, ok := <-w.workerCh:
			if !ok {
				logger.Error().Msg("runner channel closed")
				return nil
			}

			logger.Debug().Any("msg", v).Msg("worker message")

			if err := w.handler.Handle(ctx, v); err != nil {
				logger.Error().Err(err).Msg("runner")
			}
		}
	}
}
