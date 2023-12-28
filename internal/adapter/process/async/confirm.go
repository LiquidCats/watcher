package async

import (
	"context"
	"time"
)

func (w *WatchProcess) Confirm(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * time.Duration(w.cfg.Interval))
	defer ticker.Stop()

	for {
		if err := w.confirmationUsecase.Confirm(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
