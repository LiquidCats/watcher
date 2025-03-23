package backbround

import (
	"context"
	"time"

	"github.com/LiquidCats/watcher/v2/pkg/graceful"
	"github.com/rs/zerolog"
)

type Job interface {
	Execute(ctx context.Context) error
}

func BackgroundRunner(
	name string,
	interval time.Duration,
	job Job,
) graceful.Runner {
	timer := time.NewTicker(interval)

	return func(ctx context.Context) error {
		defer timer.Stop()

		logger := zerolog.Ctx(ctx)

		logger.Info().Msgf("backbround processor [name: %s] started", name)
		defer logger.Info().Msgf("backbround processor [name: %s] stopped", name)

		for {
			if err := job.Execute(ctx); err != nil {
				logger.Error().Stack().Err(err).Stack().Msgf("backbround processor [name: %s] failed", name)
			}

			select {
			case <-timer.C:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
