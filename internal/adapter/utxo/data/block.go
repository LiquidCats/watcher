package data

import (
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/entities"
)

type Block[TX any] struct {
	Hash              entities.BlockHash `json:"hash"`
	Height            uint32             `json:"height"`
	Previousblockhash entities.BlockHash `json:"previousblockhash"`
	Tx                []TX               `json:"tx"`
}
