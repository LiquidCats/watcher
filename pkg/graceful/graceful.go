package graceful

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type Runner func(ctx context.Context) error

func Wait(runners ...Runner) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	for _, runner := range runners {
		group.Go(func() error {
			return runner(ctx)
		})
	}

	return group.Wait()
}
