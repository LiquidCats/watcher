package bus

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/utxo/domain/entities"
)

type NilPublisher struct{}

func (p *NilPublisher) PublishTransaction(_ context.Context, _ *entities.UtxoTransaction) error {
	return nil
}

func (p *NilPublisher) PublishBlock(_ context.Context, _ *entities.UtxoBlock) error {
	return nil
}
