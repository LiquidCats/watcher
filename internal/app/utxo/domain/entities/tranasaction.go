package entities

import (
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
)

type Transaction struct {
	TxID      entities.TxID      `json:"tx_id"`
	BlockHash entities.BlockHash `json:"block_hash"`
	Inputs    []*Input           `json:"inputs,omitempty"`
	Outputs   []*Output          `json:"outputs,omitempty"`
}

type Input struct {
	TxID entities.TxID `json:"tx_id"`
	Vout uint32        `json:"vout"`
}

type Output struct {
	Address entities.Address
	Value   string
	N       uint64
}
