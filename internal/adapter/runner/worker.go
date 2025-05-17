package runner

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/port/runner"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type Worker[T any] struct {
	workerCh     runner.ChanRead[T]
	handler      runner.Handler[T]
	workersCount uint
}

func NewWorker[T any](workersCount uint, workerCh runner.ChanRead[T], handler runner.Handler[T]) *Worker[T] {
	return &Worker[T]{
		workersCount: workersCount,
		workerCh:     workerCh,
		handler:      handler,
	}
}

func (w *Worker[T]) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < int(w.workersCount); i++ {
		g.Go(func() error {
			return w.runner(ctx)
		})
	}

	return g.Wait()
}

func (w *Worker[T]) runner(ctx context.Context) error {
	logger := zerolog.Ctx(ctx).With().Str("name", "runner").Logger()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case v, ok := <-w.workerCh:
			if !ok {
				logger.Error().Msg("runner channel closed")
				return nil
			}

			if err := w.handler.Handle(ctx, v); err != nil {
				logger.Error().Err(err).Msg("runner handler error")
			}
		}
	}
}
