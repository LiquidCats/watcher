package entities

type TxID string

type Transaction struct {
	TxID      TxID      `json:"tx_id"`
	BlockHash BlockHash `json:"block_hash"`
	Inputs    []Input   `json:"inputs"`
	Outputs   []Output  `json:"outputs"`
}

type Input struct {
	TxID TxID   `json:"tx_id"`
	Vout uint64 `json:"vout"`
}
type Output struct {
	Address Address
	Value   string
	N       uint64
}
