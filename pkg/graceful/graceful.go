package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

var (
	gracefulCtx, gracefulCancel = context.WithCancel(context.Background())
)

type Runner func(ctx context.Context) error

func Signals(ctx context.Context) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigs)

	for {
		select {
		case <-sigs:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// WithContext allows setting a custom context for the graceful shutdown.
func WithContext(ctx context.Context) {
	gracefulCtx, gracefulCancel = context.WithCancel(ctx)
}

func Wait(runners ...Runner) error {
	defer gracefulCancel()

	group, ctx := errgroup.WithContext(gracefulCtx)

	// Start a goroutine for each runner
	for _, r := range runners {
		runner := r
		group.Go(func() error {
			return runner(ctx)
		})
	}

	return group.Wait()
}
