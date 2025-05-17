package runner

import (
	"context"
	"time"

	"github.com/LiquidCats/watcher/v2/internal/app/port/runner"
	"github.com/rs/zerolog"
)

type Processor struct {
	name     string
	interval time.Duration
	job      runner.Job
}

func NewProcessor(
	name string,
	interval time.Duration,
	job runner.Job,
) *Processor {
	return &Processor{
		name:     name,
		interval: interval,
		job:      job,
	}
}

func (bp *Processor) Run(ctx context.Context) error {
	ticker := time.NewTicker(bp.interval)
	defer ticker.Stop()

	logger := zerolog.Ctx(ctx)

	logger.Info().Msgf("background processor [name: %s] started", bp.name)
	defer logger.Info().Msgf("background processor [name: %s] stopped", bp.name)

	for {
		if err := bp.job.Handle(ctx); err != nil {
			logger.Error().Stack().Err(err).Stack().Msgf("background processor [name: %s] failed", bp.name)
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
