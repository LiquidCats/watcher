package entities

import (
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
)

type Block[T any] struct {
	Height       entities.BlockHeight `json:"height"`
	Hash         entities.BlockHash   `json:"hash"`
	PrevHash     entities.BlockHash   `json:"prevHash"`
	Transactions []*Transaction       `json:"transactions"`
}
