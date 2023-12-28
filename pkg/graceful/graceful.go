package graceful

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

type Runner func(ctx context.Context) error

func Wait(runners ...Runner) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)

	// Goroutine to handle OS signals
	group.Go(func() error {
		select {
		case <-sigs:
			fmt.Println("Signal received, initiating shutdown")
			cancel() // cancel the context
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	})

	// Start a goroutine for each runner
	for _, r := range runners {
		runner := r // Capture the current value of r in a local variable
		group.Go(func() error {
			return runner(ctx) // Use the captured variable here
		})
	}

	// Wait for all goroutines to complete
	return group.Wait()
}
