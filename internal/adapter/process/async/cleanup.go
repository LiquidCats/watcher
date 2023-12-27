package async

import (
	"context"
	"time"
)

func (w *WatchProcess) Cleanup(ctx context.Context) error {
	if !w.cfg.Cleanup.Enabled {
		return nil
	}

	ticker := time.NewTicker(time.Second * time.Duration(w.cfg.Interval))
	defer ticker.Stop()

	for {
		if err := w.cleanupUsecase.Clean(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
