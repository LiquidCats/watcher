package data

import (
	"github.com/LiquidCats/watcher/v2/internal/adapter/rpc/evm/data/common"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type Block[T any] struct {
	Hash         entities.BlockHash `json:"hash"`
	Number       *common.Big        `json:"number"`
	ParentHash   entities.BlockHash `json:"parentHash"`
	Transactions []T                `json:"transactions"`
}
