package runner

import (
	"context"
	"time"

	"github.com/LiquidCats/watcher/v2/internal/app/port/runner"
	"github.com/rotisserie/eris"
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

	logger := zerolog.Ctx(ctx).With().Str("job_name", bp.name).Logger()

	logger.Info().Msg("background processor started")
	defer logger.Info().Msg("background processor stopped")

	for {
		if err := bp.job.Handle(ctx); err != nil {
			logger.Error().Any("err", eris.ToJSON(err, true)).Msg("background processor failed")
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
