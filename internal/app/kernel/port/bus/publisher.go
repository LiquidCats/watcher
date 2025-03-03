package bus

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/utxo/domain/entities"
)

type TransactionPublisher interface {
	PublishTransaction(ctx context.Context, transaction *entities.Transaction) error
}
