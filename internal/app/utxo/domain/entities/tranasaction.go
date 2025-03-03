package entities

import (
	entities2 "github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
)

type Transaction struct {
	TxID      entities2.TxID      `json:"tx_id"`
	BlockHash entities2.BlockHash `json:"block_hash"`
	Inputs    []Input             `json:"inputs,omitempty"`
	Outputs   []Output            `json:"outputs,omitempty"`
}

type Input struct {
	TxID entities2.TxID `json:"tx_id"`
	Vout uint32         `json:"vout"`
}

type Output struct {
	Address entities2.Address
	Value   string
	N       uint64
}
