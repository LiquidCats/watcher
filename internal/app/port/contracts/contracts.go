package contracts

import (
	"context"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type Getter interface {
	GetContractInfoByAddress(ctx context.Context, contract entities.Address) (*entities.Contract, error)
}
