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

func RunOnBackground(
	name string,
	interval time.Duration,
	job Job,
) graceful.Runner {
	return func(ctx context.Context) error {
		logger := zerolog.Ctx(ctx)

		logger.Info().Msgf("backbround processor [%s] started", name)
		defer logger.Info().Msgf("backbround processor [%s] stopped", name)

		timer := time.NewTimer(interval)
		defer timer.Stop()

		for {
			if err := job.Execute(ctx); err != nil {
				logger.Error().Err(err).Msgf("backbround processor [%s] failed", name)
			}

			select {
			case <-timer.C:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
