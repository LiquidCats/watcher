package repository

import "context"

type P2PClient interface {
	Subscribe(ctx context.Context) error
	Unsubscribe(ctx context.Context) error
}
