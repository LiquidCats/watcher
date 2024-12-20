package entities

import "github.com/LiquidCats/watcher/v2/internal/app/domain/shared/entities"

type TxHash string

type Transaction struct {
	TxId      TxHash             `json:"tx_id"`
	BlockHash entities.BlockHash `json:"block_hash"`
	Inputs    []Input            `json:"inputs"`
	Outputs   []Output           `json:"outputs"`
}

type Input struct {
	TxId TxHash `json:"tx_id"`
	Vout uint64 `json:"vout"`
}
type Output struct {
	Address entities.Address
	Value   string
	N       uint64
}
