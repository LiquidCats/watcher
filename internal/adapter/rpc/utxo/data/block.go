package data

import (
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type Block[T any] struct {
	Hash              entities.BlockHash   `json:"hash"`
	Height            entities.BlockHeight `json:"height"`
	PreviousBlockHash entities.BlockHash   `json:"previousblockhash"`
	Tx                []T                  `json:"tx"`
}
